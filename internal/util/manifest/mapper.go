// Copyright 2020-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manifest

import (
	"context"
	"fmt"
	modulev1alpha "github.com/ProtobufMan/bufman/internal/gen/module/v1alpha"
	"go.uber.org/multierr"
	"io"
)

var (
	protoDigestTypeToDigestType = map[modulev1alpha.DigestType]DigestType{
		modulev1alpha.DigestType_DIGEST_TYPE_SHAKE256: DigestTypeShake256,
	}
	digestTypeToProtoDigestType = map[DigestType]modulev1alpha.DigestType{
		DigestTypeShake256: modulev1alpha.DigestType_DIGEST_TYPE_SHAKE256,
	}
)

// NewDigestFromProtoDigest maps a *gen.Digest to a Digest.
func NewDigestFromProtoDigest(digest *modulev1alpha.Digest) (*Digest, error) {
	if digest == nil {
		return nil, fmt.Errorf("nil digest")
	}
	dType, ok := protoDigestTypeToDigestType[digest.DigestType]
	if !ok {
		return nil, fmt.Errorf("unsupported digest kind: %s", digest.DigestType.String())
	}
	return NewDigestFromBytes(dType, digest.Digest)
}

// AsProtoBlob returns the passed blob as a proto module blob.
func AsProtoBlob(ctx context.Context, b Blob) (_ *modulev1alpha.Blob, retErr error) {
	digestType, ok := digestTypeToProtoDigestType[b.Digest().Type()]
	if !ok {
		return nil, fmt.Errorf("digest type %q not supported by module proto", b.Digest().Type())
	}
	rc, err := b.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot open blob: %w", err)
	}
	defer func() {
		retErr = multierr.Append(retErr, rc.Close())
	}()
	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("cannot read blob contents: %w", err)
	}
	return &modulev1alpha.Blob{
		Digest: &modulev1alpha.Digest{
			DigestType: digestType,
			Digest:     b.Digest().Bytes(),
		},
		Content: content,
	}, nil
}

// NewManifestFromProto returns a Manifest from a proto module blob. It makes sure the
// digest and content matches.
func NewManifestFromProto(ctx context.Context, b *modulev1alpha.Blob) (_ *Manifest, retErr error) {
	blob, err := NewBlobFromProto(b)
	if err != nil {
		return nil, fmt.Errorf("invalid manifest: %w", err)
	}
	r, err := blob.Open(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		retErr = multierr.Append(retErr, r.Close())
	}()

	m, err := NewFromReader(r)
	m.digest, err = NewDigestFromProtoDigest(b.Digest)
	if err != nil {
		return nil, err
	}

	return m, err
}

// NewBlobSetFromProto returns a BlobSet from a slice of proto module blobs.
// It makes sure the digest and content matches for each blob.
func NewBlobSetFromProto(ctx context.Context, blobs []*modulev1alpha.Blob) (*BlobSet, error) {
	var memBlobs []Blob
	for i, modBlob := range blobs {
		memBlob, err := NewBlobFromProto(modBlob)
		if err != nil {
			return nil, fmt.Errorf("invalid blob at index %d: %w", i, err)
		}
		memBlobs = append(memBlobs, memBlob)
	}
	return NewBlobSet(ctx, memBlobs)
}

// NewBlobFromProto returns a Blob from a proto module blob. It makes sure the
// digest and content matches.
func NewBlobFromProto(b *modulev1alpha.Blob) (Blob, error) {
	if b == nil {
		return nil, fmt.Errorf("nil blob")
	}
	digest, err := NewDigestFromProtoDigest(b.Digest)
	if err != nil {
		return nil, fmt.Errorf("digest from proto digest: %w", err)
	}
	memBlob, err := NewMemoryBlob(
		*digest,
		b.Content,
		MemoryBlobWithDigestValidation(),
	)
	if err != nil {
		return nil, fmt.Errorf("new memory blob: %w", err)
	}
	return memBlob, nil
}

// ToProtoManifestAndBlobs converts a Manifest and BlobSet to the protobuf types.
func ToProtoManifestAndBlobs(ctx context.Context, manifest *Manifest, blobs *BlobSet) (*modulev1alpha.Blob, []*modulev1alpha.Blob, error) {
	manifestBlob, err := manifest.Blob()
	if err != nil {
		return nil, nil, err
	}
	manifestProtoBlob, err := AsProtoBlob(ctx, manifestBlob)
	if err != nil {
		return nil, nil, err
	}
	filesBlobs := blobs.Blobs()
	filesProtoBlobs := make([]*modulev1alpha.Blob, len(filesBlobs))
	for i, b := range filesBlobs {
		pb, err := AsProtoBlob(ctx, b)
		if err != nil {
			return nil, nil, err
		}
		filesProtoBlobs[i] = pb
	}
	return manifestProtoBlob, filesProtoBlobs, nil
}

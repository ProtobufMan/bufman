package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util/es"
	"io"
)

type ESStorageHelperImpl struct {
	EsClient es.Client
}

type FileMapping struct {
	Content []byte
}

func (helper *ESStorageHelperImpl) StoreBlob(ctx context.Context, blob *model.FileBlob) error {
	return helper.store(ctx, constant.ESFileBlobIndex, blob.Digest, blob)
}

func (helper *ESStorageHelperImpl) StoreManifest(ctx context.Context, manifest *model.FileManifest) error {
	return helper.store(ctx, constant.ESManifestIndex, manifest.Digest, manifest)
}

func (helper *ESStorageHelperImpl) StoreDocumentation(ctx context.Context, blob *model.FileBlob) error {
	return helper.store(ctx, constant.ESDocumentIndex, blob.Digest, blob)
}

func (helper *ESStorageHelperImpl) store(ctx context.Context, index, digest string, v interface{}) error {
	// 转为json
	jsonBody, err := json.Marshal(&v)
	if err != nil {
		return err
	}

	// 存储在es中
	err = helper.EsClient.Create(ctx, index, digest, jsonBody)
	if err != nil {
		return err
	}

	return nil
}

func (helper *ESStorageHelperImpl) ReadBlobToReader(ctx context.Context, digest string) (io.Reader, error) {
	content, err := helper.ReadBlob(ctx, digest)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

func (helper *ESStorageHelperImpl) ReadBlob(ctx context.Context, digest string) ([]byte, error) {
	b := &model.FileBlob{}
	err := helper.read(ctx, constant.ESManifestIndex, digest, b)
	if err != nil {
		return nil, err
	}

	return b.Content, nil
}

func (helper *ESStorageHelperImpl) ReadManifestToReader(ctx context.Context, digest string) (io.Reader, error) {
	content, err := helper.ReadManifest(ctx, digest)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

func (helper *ESStorageHelperImpl) ReadManifest(ctx context.Context, digest string) ([]byte, error) {
	m := &model.FileManifest{}
	err := helper.read(ctx, constant.ESManifestIndex, digest, m)
	if err != nil {
		return nil, err
	}

	return m.Content, nil
}

func (helper *ESStorageHelperImpl) read(ctx context.Context, index string, digest string, v interface{}) error {
	// 存储在es中
	data, err := helper.EsClient.Find(ctx, index, digest)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	return nil
}

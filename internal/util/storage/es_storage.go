package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/util/es"
	"io"
)

type ESStorageHelperImpl struct {
	EsClient es.Client
}

type FileMapping struct {
	Content []byte
}

func (helper *ESStorageHelperImpl) StoreBlobFromReader(ctx context.Context, digest string, readCloser io.ReadCloser) error {
	// 读取数据
	content, err := io.ReadAll(readCloser)
	defer readCloser.Close()
	if err != nil {
		return err
	}

	return helper.StoreBlob(ctx, digest, content)
}

func (helper *ESStorageHelperImpl) StoreBlob(ctx context.Context, digest string, content []byte) error {
	return helper.store(ctx, constant.ESFileBlobIndex, digest, content)
}

func (helper *ESStorageHelperImpl) StoreManifestFromReader(ctx context.Context, digest string, readCloser io.ReadCloser) error {
	// 读取数据
	content, err := io.ReadAll(readCloser)
	defer readCloser.Close()
	if err != nil {
		return err
	}

	return helper.StoreManifest(ctx, digest, content)
}

func (helper *ESStorageHelperImpl) StoreManifest(ctx context.Context, digest string, content []byte) error {
	return helper.store(ctx, constant.ESManifestIndex, digest, content)
}

func (helper *ESStorageHelperImpl) store(ctx context.Context, index, digest string, content []byte) error {
	// 转为json
	jsonBody, err := json.Marshal(FileMapping{
		Content: content,
	})
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
	return helper.read(ctx, constant.ESFileBlobIndex, digest)
}

func (helper *ESStorageHelperImpl) ReadManifestToReader(ctx context.Context, digest string) (io.Reader, error) {
	content, err := helper.ReadManifest(ctx, digest)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

func (helper *ESStorageHelperImpl) ReadManifest(ctx context.Context, fileName string) ([]byte, error) {
	return helper.read(ctx, constant.ESManifestIndex, fileName)
}

func (helper *ESStorageHelperImpl) read(ctx context.Context, index string, digest string) ([]byte, error) {
	// 存储在es中
	data, err := helper.EsClient.Find(ctx, index, digest)
	if err != nil {
		return nil, err
	}

	mapping := &FileMapping{}
	err = json.Unmarshal(data, mapping)
	if err != nil {
		return nil, err
	}

	return mapping.Content, nil
}

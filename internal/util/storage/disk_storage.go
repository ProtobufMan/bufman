package storage

import (
	"bytes"
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	"io"
	"os"
	"path"
	"sync"
)

type DiskStorageHelperImpl struct {
	mu     sync.Mutex
	muDict map[string]*sync.RWMutex

	pluginMu     sync.Mutex
	pluginMuDict map[string]*sync.RWMutex
}

func (helper *DiskStorageHelperImpl) StoreBlobFromReader(ctx context.Context, digest string, readCloser io.ReadCloser) error {
	content, err := io.ReadAll(readCloser)
	defer readCloser.Close()
	if err != nil {
		return err
	}

	return helper.StoreBlob(ctx, digest, content)
}

func (helper *DiskStorageHelperImpl) StoreBlob(ctx context.Context, digest string, content []byte) error {
	return helper.store(ctx, digest, content)
}

func (helper *DiskStorageHelperImpl) store(ctx context.Context, digest string, content []byte) error {
	helper.mu.Lock()
	defer helper.mu.Unlock()

	if _, ok := helper.muDict[digest]; !ok {
		helper.muDict[digest] = &sync.RWMutex{}
	}

	// 上写锁
	helper.muDict[digest].Lock()
	defer helper.muDict[digest].Unlock()

	// 打开文件
	filePath := helper.GetFilePath(digest)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if os.IsExist(err) {
		// 已经存在，直接返回
		return nil
	}
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func (helper *DiskStorageHelperImpl) StoreManifestFromReader(ctx context.Context, digest string, readCloser io.ReadCloser) error {
	content, err := io.ReadAll(readCloser)
	defer readCloser.Close()
	if err != nil {
		return err
	}

	return helper.StoreManifest(ctx, digest, content)
}

func (helper *DiskStorageHelperImpl) StoreManifest(ctx context.Context, digest string, content []byte) error {
	return helper.store(ctx, digest, content)
}

func (helper *DiskStorageHelperImpl) ReadBlobToReader(ctx context.Context, fileName string) (io.Reader, error) {
	content, err := helper.ReadBlob(ctx, fileName)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

func (helper *DiskStorageHelperImpl) ReadBlob(ctx context.Context, fileName string) ([]byte, error) {
	return helper.read(ctx, fileName)
}

func (helper *DiskStorageHelperImpl) ReadManifestToReader(ctx context.Context, fileName string) (io.Reader, error) {
	content, err := helper.ReadManifest(ctx, fileName)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

func (helper *DiskStorageHelperImpl) ReadManifest(ctx context.Context, fileName string) ([]byte, error) {
	return helper.read(ctx, fileName)
}

func (helper *DiskStorageHelperImpl) read(ctx context.Context, fileName string) ([]byte, error) {
	helper.mu.Lock()
	defer helper.mu.Unlock()

	if _, ok := helper.muDict[fileName]; !ok {
		helper.muDict[fileName] = &sync.RWMutex{}
	}

	// 上读锁
	helper.muDict[fileName].RLock()
	defer helper.muDict[fileName].RUnlock()

	// 读取文件
	filePath := helper.GetFilePath(fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (helper *DiskStorageHelperImpl) GetFilePath(fileName string) string {
	return path.Join(constant.FileSavaDir, fileName)
}

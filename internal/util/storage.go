package util

import (
	"bytes"
	"github.com/ProtobufMan/bufman/internal/constant"
	"io"
	"os"
	"path"
	"sync"
)

type StorageHelper interface {
	Store(fileName string, readCloser io.ReadCloser) error // 存储内容
	Read(fileName string) (io.Reader, error)               // 读取内容
	GetFilePath(fileName string) string                    // 获取文件实际存储地址
}

type StorageHelperImpl struct {
	mu     sync.Mutex
	muDict map[string]*sync.RWMutex
}

var storageHelperImpl = &StorageHelperImpl{
	muDict: map[string]*sync.RWMutex{},
}

func NewStorageHelper() StorageHelper {
	return storageHelperImpl
}

func (helper *StorageHelperImpl) Store(fileName string, readCloser io.ReadCloser) error {
	helper.mu.Lock()
	defer helper.mu.Unlock()

	if _, ok := helper.muDict[fileName]; !ok {
		helper.muDict[fileName] = &sync.RWMutex{}
	}

	// 上写锁
	helper.muDict[fileName].Lock()
	defer helper.muDict[fileName].Unlock()

	// 打开文件
	filePath := helper.GetFilePath(fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if os.IsExist(err) {
		// 已经存在，直接返回
		return nil
	}
	if err != nil {
		return err
	}

	// 写入文件
	_, err = io.Copy(file, readCloser)
	if err != nil {
		return err
	}

	return readCloser.Close()
}

func (helper *StorageHelperImpl) Read(fileName string) (io.Reader, error) {
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

	return bytes.NewReader(content), nil
}

func (helper *StorageHelperImpl) GetFilePath(fileName string) string {
	return path.Join(constant.FileSavaDir, fileName)
}

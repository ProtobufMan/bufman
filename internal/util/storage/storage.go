package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/ProtobufMan/bufman/internal/constant"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type StorageHelper interface {
	Store(digest string, readCloser io.ReadCloser) error // 存储内容
	Read(digest string) (io.Reader, error)               // 读取内容
	GetFilePath(digest string) string                    // 获取文件实际存储地址

	StorePlugin(pluginName string, version string, reversion uint32, binaryData []byte) (fileName string, err error) // 存储插件
	GetPluginFileName(pluginName string, version string, reversion uint32, binaryData []byte) string                 // 获取插件名称
}

type StorageHelperImpl struct {
	mu     sync.Mutex
	muDict map[string]*sync.RWMutex

	pluginMu     sync.Mutex
	pluginMuDict map[string]*sync.RWMutex
}

func (helper *StorageHelperImpl) StorePlugin(pluginName string, version string, reversion uint32, binaryData []byte) (fileName string, err error) {
	fileName = helper.GetPluginFileName(pluginName, version, reversion, binaryData)

	helper.pluginMu.Lock()
	defer helper.pluginMu.Unlock()

	if _, ok := helper.pluginMuDict[fileName]; !ok {
		helper.pluginMuDict[fileName] = &sync.RWMutex{}
	}

	// 上写锁
	helper.pluginMuDict[fileName].Lock()
	defer helper.pluginMuDict[fileName].Unlock()

	// 打开文件
	if !strings.HasSuffix(pluginName, ".wasm") {
		fileName = "proto-gen-" + fileName
		if runtime.GOOS == "windows" {
			fileName = fileName + ".exe"
		}
	}

	filePath := filepath.Join(constant.PluginSaveDir, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	defer file.Close()
	if os.IsExist(err) {
		// 已经存在，直接返回
		return fileName, nil
	}
	if err != nil {
		return "", err
	}

	// 写入文件

	_, err = file.Write(binaryData)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (helper *StorageHelperImpl) GetPluginFileName(pluginName string, version string, reversion uint32, binaryData []byte) string {
	sha := sha256.New()
	sha.Write([]byte(pluginName))
	sha.Write([]byte(version))
	sha.Write([]byte(strconv.Itoa(int(reversion))))
	sha.Write(binaryData)
	fileName := hex.EncodeToString(sha.Sum(nil))

	if strings.HasSuffix(pluginName, ".wasm") {
		fileName = fileName + ".wasm"
	}

	return fileName
}

var storageHelperImpl = &StorageHelperImpl{
	muDict:       map[string]*sync.RWMutex{},
	pluginMuDict: map[string]*sync.RWMutex{},
}

func NewStorageHelper() StorageHelper {
	return storageHelperImpl
}

func (helper *StorageHelperImpl) Store(digest string, readCloser io.ReadCloser) error {
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

package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufconfig"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/model"
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
	StoreFromReader(digest string, readCloser io.ReadCloser) error // 存储内容
	Store(digest string, content []byte) error
	ReadToReader(digest string) (io.Reader, error) // 读取内容
	Read(fileName string) ([]byte, error)
	GetDocumentAndLicenseFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, manifest.Blob, error)
	GetBufManConfigFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error)
	GetDocumentFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error)
	GetLicenseFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error)
	ReadToManifestAndBlobSet(ctx context.Context, modelFileManifest *model.FileManifest, fileBlobs model.FileBlobs) (*manifest.Manifest, *manifest.BlobSet, error) // 读取为manifest和blob set

	StorePlugin(pluginName string, version string, reversion uint32, binaryData []byte) (fileName string, err error) // 存储插件
}

type File struct {
	FileName string
	Digest   string
	Content  []byte
}

type StorageHelperImpl struct {
	mu     sync.Mutex
	muDict map[string]*sync.RWMutex

	pluginMu     sync.Mutex
	pluginMuDict map[string]*sync.RWMutex
}

var storageHelperImpl = &StorageHelperImpl{
	muDict:       map[string]*sync.RWMutex{},
	pluginMuDict: map[string]*sync.RWMutex{},
}

func NewStorageHelper() StorageHelper {
	return storageHelperImpl
}

func (helper *StorageHelperImpl) GetDocumentAndLicenseFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, manifest.Blob, error) {
	var documentDataExists, licenseExists bool
	var documentBlob, licenseBlob manifest.Blob

	externalPaths := []string{
		bufmodule.LicenseFilePath,
	}
	externalPaths = append(externalPaths, bufmodule.AllDocumentationPaths...)

	err := fileManifest.Range(func(path string, digest manifest.Digest) error {
		blob, ok := blobSet.BlobFor(digest.String())
		if !ok {
			// 文件清单中有的文件，在file blobs中没有
			return errors.New("check manifest and file blobs failed")
		}

		// 如果遇到配置文件，就记录下来
		for _, externalPath := range externalPaths {
			if documentDataExists && licenseExists {
				break
			}

			if path == externalPath {
				if path == bufmodule.LicenseFilePath {
					// license文件
					licenseBlob = blob
					licenseExists = true
				} else {
					if documentDataExists {
						break
					}
					// document文件
					documentBlob = blob
					documentDataExists = true
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return documentBlob, licenseBlob, nil
}

func (helper *StorageHelperImpl) GetBufManConfigFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error) {
	var configFileExist bool
	var configFileBlob manifest.Blob

	err := fileManifest.Range(func(path string, digest manifest.Digest) error {
		blob, ok := blobSet.BlobFor(digest.String())
		if !ok {
			// 文件清单中有的文件，在file blobs中没有
			return errors.New("check manifest and file blobs failed")
		}

		// 如果遇到配置文件，就记录下来
		for _, configFilePath := range bufconfig.AllConfigFilePaths {
			if configFileExist {
				break
			}

			if path == configFilePath {
				configFileBlob = blob
				configFileExist = true
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return configFileBlob, nil
}

func (helper *StorageHelperImpl) GetDocumentFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error) {
	var documentExist bool
	var documentBlob manifest.Blob

	err := fileManifest.Range(func(path string, digest manifest.Digest) error {
		blob, ok := blobSet.BlobFor(digest.String())
		if !ok {
			// 文件清单中有的文件，在file blobs中没有
			return errors.New("check manifest and file blobs failed")
		}

		// 如果遇到README文件，就记录下来
		for _, documentationPath := range bufmodule.AllDocumentationPaths {
			if documentExist {
				break
			}

			if path == documentationPath {
				documentBlob = blob
				documentExist = true
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return documentBlob, nil
}

func (helper *StorageHelperImpl) GetLicenseFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error) {
	var licenseExist bool
	var licenseBlob manifest.Blob

	err := fileManifest.Range(func(path string, digest manifest.Digest) error {
		if licenseExist {
			return nil
		}

		blob, ok := blobSet.BlobFor(digest.String())
		if !ok {
			// 文件清单中有的文件，在file blobs中没有
			return errors.New("check manifest and file blobs failed")
		}

		// 如果遇到license，就记录下来
		if path == bufmodule.LicenseFilePath {
			licenseBlob = blob
			licenseExist = true
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return licenseBlob, nil
}

func (helper *StorageHelperImpl) ReadToManifestAndBlobSet(ctx context.Context, modelFileManifest *model.FileManifest, fileBlobs model.FileBlobs) (*manifest.Manifest, *manifest.BlobSet, error) {
	// 读取文件清单
	reader, err := helper.ReadToReader(modelFileManifest.Digest)
	if err != nil {
		return nil, nil, err
	}
	fileManifest, err := manifest.NewFromReader(reader)
	if err != nil {
		return nil, nil, err
	}

	// 读取文件blobs
	blobs := make([]manifest.Blob, 0, len(fileBlobs))
	for i := 0; i < len(fileBlobs); i++ {
		// 读取文件
		reader, err := helper.ReadToReader(fileBlobs[i].Digest)
		if err != nil {
			return nil, nil, err
		}

		// 生成blob
		blob, err := manifest.NewMemoryBlobFromReader(reader)
		if err != nil {
			return nil, nil, err
		}
		blobs = append(blobs, blob)
	}

	blobSet, err := manifest.NewBlobSet(ctx, blobs)
	if err != nil {
		return nil, nil, err
	}

	return fileManifest, blobSet, nil
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

func (helper *StorageHelperImpl) StoreFromReader(digest string, readCloser io.ReadCloser) error {
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

func (helper *StorageHelperImpl) Store(digest string, content []byte) error {
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

func (helper *StorageHelperImpl) ReadToReader(fileName string) (io.Reader, error) {
	content, err := helper.Read(fileName)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(content), nil
}

func (helper *StorageHelperImpl) Read(fileName string) ([]byte, error) {
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

func (helper *StorageHelperImpl) GetFilePath(fileName string) string {
	return path.Join(constant.FileSavaDir, fileName)
}

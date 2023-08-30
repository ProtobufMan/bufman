package storage

import (
	"context"
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufconfig"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/config"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util/es"
	"io"
	"sync"
)

type BaseStorageHelper interface {
	StoreBlob(ctx context.Context, blob *model.FileBlob) error
	StoreManifest(ctx context.Context, manifest *model.FileManifest) error
	StoreDocumentation(ctx context.Context, blob *model.FileBlob) error
	ReadBlobToReader(ctx context.Context, digest string) (io.Reader, error) // 读取内容
	ReadBlob(ctx context.Context, fileName string) ([]byte, error)
	ReadManifestToReader(ctx context.Context, fileName string) (io.Reader, error)
	ReadManifest(ctx context.Context, fileName string) ([]byte, error)
}

type StorageHelper interface {
	BaseStorageHelper
	ReadToManifestAndBlobSet(ctx context.Context, modelFileManifest *model.FileManifest, fileBlobs model.FileBlobs) (*manifest.Manifest, *manifest.BlobSet, error) // 读取为manifest和blob set
	GetDocumentAndLicenseFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, manifest.Blob, error)
	GetBufManConfigFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error)
	GetDocumentFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error)
	GetLicenseFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (manifest.Blob, error)
}

type StorageHelperImpl struct {
	BaseStorageHelper
}

// 单例模式
var storageHelperImpl *StorageHelperImpl
var once sync.Once

func NewStorageHelper() StorageHelper {
	if storageHelperImpl == nil {
		// 对象初始化
		once.Do(func() {
			if len(config.Properties.ElasticSearch.Urls) == 0 {
				storageHelperImpl = &StorageHelperImpl{
					BaseStorageHelper: &DiskStorageHelperImpl{
						muDict:       map[string]*sync.RWMutex{},
						pluginMuDict: map[string]*sync.RWMutex{},
					},
				}
			} else {
				esClient, err := es.NewEsClient(config.Properties.ElasticSearch.Username, config.Properties.ElasticSearch.Password, config.Properties.ElasticSearch.Urls...)
				if err != nil {
					panic(err)
				}

				storageHelperImpl = &StorageHelperImpl{
					BaseStorageHelper: &ESStorageHelperImpl{
						EsClient: esClient,
					},
				}
			}
		})
	}

	return storageHelperImpl
}

func (helper *StorageHelperImpl) ReadToManifestAndBlobSet(ctx context.Context, modelFileManifest *model.FileManifest, fileBlobs model.FileBlobs) (*manifest.Manifest, *manifest.BlobSet, error) {
	// 读取文件清单
	reader, err := helper.ReadManifestToReader(ctx, modelFileManifest.Digest)
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
		reader, err := helper.ReadBlobToReader(ctx, fileBlobs[i].Digest)
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

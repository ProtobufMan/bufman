package services

import (
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/util/storage"
	"gorm.io/gorm"
)

type DownloadService interface {
	DownloadManifestAndBlobs(registerID string, reference string) (*manifest.Manifest, *manifest.BlobSet, e.ResponseError)
}

type DownloadServiceImpl struct {
	commitMapper  mapper.CommitMapper
	fileMapper    mapper.FileMapper
	storageHelper storage.StorageHelper
}

func NewDownloadService() DownloadService {
	return &DownloadServiceImpl{
		commitMapper:  &mapper.CommitMapperImpl{},
		fileMapper:    &mapper.FileMapperImpl{},
		storageHelper: storage.NewStorageHelper(),
	}
}

func (downloadService *DownloadServiceImpl) DownloadManifestAndBlobs(registerID string, reference string) (*manifest.Manifest, *manifest.BlobSet, e.ResponseError) {
	// 查询reference对应的commit
	commit, err := downloadService.commitMapper.FindByRepositoryIDAndReference(registerID, reference)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, e.NewNotFoundError(fmt.Sprintf("reference %s", reference))
		}

		return nil, nil, e.NewInternalError(err.Error())
	}

	// 查询文件清单
	modelFileManifest, err := downloadService.fileMapper.FindManifestByCommitID(commit.CommitID)
	if err != nil {
		if err != nil {
			return nil, nil, e.NewInternalError(err.Error())
		}
	}

	// 接着查询blobs
	fileBlobs, err := downloadService.fileMapper.FindAllBlobsByCommitID(commit.CommitID)
	if err != nil {
		return nil, nil, e.NewInternalError(err.Error())
	}

	// 读取
	fileManifest, blobSet, err := downloadService.storageHelper.ReadToManifestAndBlobSet(modelFileManifest, fileBlobs)
	if err != nil {
		return nil, nil, e.NewInternalError(err.Error())
	}

	return fileManifest, blobSet, nil
}

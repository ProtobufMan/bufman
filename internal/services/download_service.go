package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
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

		return nil, nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	// 查询文件清单
	modelFileManifest, err := downloadService.fileMapper.FindManifestByCommitID(commit.CommitID)
	if err != nil {
		if err != nil {
			return nil, nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
		}
	}

	// 接着查询blobs
	fileBlobs, err := downloadService.fileMapper.FindAllBlobsByCommitID(commit.CommitID)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	// 读取文件清单
	reader, err := downloadService.storageHelper.Read(modelFileManifest.Digest)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}
	fileManifest, err := manifest.NewFromReader(reader)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	// 读取文件blobs
	blobs := make([]manifest.Blob, 0, len(fileBlobs))
	for i := 0; i < len(fileBlobs); i++ {
		// 读取文件
		reader, err := downloadService.storageHelper.Read(fileBlobs[i].Digest)
		if err != nil {
			return nil, nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
		}

		// 生成blob
		blob, err := manifest.NewMemoryBlobFromReader(reader)
		if err != nil {
			return nil, nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
		}
		blobs = append(blobs, blob)
	}

	blobSet, err := manifest.NewBlobSet(context.Background(), blobs)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	return fileManifest, blobSet, nil
}

package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman/internal/e"
	modulev1alpha "github.com/ProtobufMan/bufman/internal/gen/module/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/util"
	"github.com/ProtobufMan/bufman/internal/util/manifest"
	"gorm.io/gorm"
)

type DownloadService interface {
	DownloadManifestAndBlobs(registerID string, reference string) (*modulev1alpha.Blob, []*modulev1alpha.Blob, e.ResponseError)
}

type DownloadServiceImpl struct {
	commitMapper  mapper.CommitMapper
	fileMapper    mapper.FileMapper
	storageHelper util.StorageHelper
}

func NewDownloadService() DownloadService {
	return &DownloadServiceImpl{
		commitMapper:  &mapper.CommitMapperImpl{},
		fileMapper:    &mapper.FileMapperImpl{},
		storageHelper: util.NewStorageHelper(),
	}
}

func (downloadService *DownloadServiceImpl) DownloadManifestAndBlobs(registerID string, reference string) (*modulev1alpha.Blob, []*modulev1alpha.Blob, e.ResponseError) {
	// 查询reference对应的commit
	commit, err := downloadService.commitMapper.FindByRepositoryIDAndReference(registerID, reference)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, e.NewNotFoundError(fmt.Sprintf("reference %s", reference))
		}

		return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	// 查询文件清单
	modelFileManifest, err := downloadService.fileMapper.FindManifestByCommitID(commit.CommitID)
	if err != nil {
		if err != nil {
			return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
		}
	}

	// 接着查询blobs
	fileBlobs, err := downloadService.fileMapper.FindAllBlobsByCommitID(commit.CommitID)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	// 读取文件清单
	reader, err := downloadService.storageHelper.Read(modelFileManifest.Digest)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}
	fileManifest, err := manifest.NewFromReader(reader)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	// 读取文件blobs
	blobs := make([]manifest.Blob, 0, len(fileBlobs))
	for i := 0; i < len(fileBlobs); i++ {
		// 读取文件
		reader, err := downloadService.storageHelper.Read(fileBlobs[i].Digest)
		if err != nil {
			return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
		}

		// 生成blob
		blob, err := manifest.NewMemoryBlobFromReader(reader)
		if err != nil {
			return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
		}
		blobs = append(blobs, blob)
	}

	blobSet, err := manifest.NewBlobSet(context.Background(), blobs)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	protoManifest, protoBlobs, err := manifest.ToProtoManifestAndBlobs(context.Background(), fileManifest, blobSet)
	if err != nil {
		return nil, nil, e.NewInternalError(registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure)
	}

	return protoManifest, protoBlobs, nil
}

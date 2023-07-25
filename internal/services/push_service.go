package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util/security"
	"github.com/ProtobufMan/bufman/internal/util/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PushService interface {
	PushManifestAndBlobs(ctx context.Context, userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet) (*model.Commit, e.ResponseError)
	PushManifestAndBlobsWithTags(ctx context.Context, userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet, tagNames []string) (*model.Commit, e.ResponseError)
	PushManifestAndBlobsWithDraft(ctx context.Context, userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet, draftName string) (*model.Commit, e.ResponseError)
	GetManifestAndBlobSet(ctx context.Context, repositoryID string, reference string) (*manifest.Manifest, *manifest.BlobSet, e.ResponseError)
}

type PushServiceImpl struct {
	userMapper       mapper.UserMapper
	repositoryMapper mapper.RepositoryMapper
	fileMapper       mapper.FileMapper
	commitMapper     mapper.CommitMapper
	storageHelper    storage.StorageHelper
}

func NewPushService() PushService {
	return &PushServiceImpl{
		userMapper:       &mapper.UserMapperImpl{},
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
		fileMapper:       &mapper.FileMapperImpl{},
		storageHelper:    storage.NewStorageHelper(),
	}
}

func (pushService *PushServiceImpl) GetManifestAndBlobSet(ctx context.Context, repositoryID string, reference string) (*manifest.Manifest, *manifest.BlobSet, e.ResponseError) {
	// 查询reference对应的commit
	commit, err := pushService.commitMapper.FindByRepositoryIDAndReference(repositoryID, reference)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, e.NewNotFoundError(fmt.Sprintf("repository %s", repositoryID))
		}

		return nil, nil, e.NewInternalError(err.Error())
	}

	// 查询文件清单
	modelFileManifest, err := pushService.fileMapper.FindManifestByCommitID(commit.CommitID)
	if err != nil {
		if err != nil {
			return nil, nil, e.NewInternalError(err.Error())
		}
	}

	// 接着查询blobs
	fileBlobs, err := pushService.fileMapper.FindAllBlobsByCommitID(commit.CommitID)
	if err != nil {
		return nil, nil, e.NewInternalError(err.Error())
	}

	// 读取
	fileManifest, blobSet, err := pushService.storageHelper.ReadToManifestAndBlobSet(ctx, modelFileManifest, fileBlobs)
	if err != nil {
		return nil, nil, e.NewInternalError(err.Error())
	}

	return fileManifest, blobSet, nil
}

func (pushService *PushServiceImpl) PushManifestAndBlobs(ctx context.Context, userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet) (*model.Commit, e.ResponseError) {
	commit, err := pushService.toCommit(userID, ownerName, repositoryName, fileManifest)
	if err != nil {
		return nil, err
	}

	// 写入文件
	err = pushService.saveFileManifestAndBlobs(ctx, fileManifest, fileBlobs)
	if err != nil {
		return nil, err
	}

	// 写入数据库
	createErr := pushService.commitMapper.Create(commit)
	if createErr != nil {
		if errors.Is(createErr, gorm.ErrDuplicatedKey) {
			return nil, e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
		}
		if errors.Is(createErr, mapper.ErrLastCommitDuplicated) {
			return nil, e.NewAlreadyExistsError("last commit")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
	}

	return commit, nil
}

func (pushService *PushServiceImpl) PushManifestAndBlobsWithTags(ctx context.Context, userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet, tagNames []string) (*model.Commit, e.ResponseError) {
	commit, err := pushService.toCommit(userID, ownerName, repositoryName, fileManifest)
	if err != nil {
		return nil, err
	}

	// 生成tags
	var tags []*model.Tag
	for i := 0; i < len(tagNames); i++ {
		tags = append(tags, &model.Tag{
			UserID:       commit.UserID,
			RepositoryID: commit.RepositoryID,
			CommitID:     commit.CommitID,
			TagID:        uuid.NewString(),
			TagName:      tagNames[i],
		})
	}
	commit.Tags = tags

	// 写入文件
	err = pushService.saveFileManifestAndBlobs(ctx, fileManifest, fileBlobs)
	if err != nil {
		return nil, err
	}

	createErr := pushService.commitMapper.Create(commit)
	if createErr != nil {
		if errors.Is(createErr, mapper.ErrTagAndDraftDuplicated) || errors.Is(createErr, gorm.ErrDuplicatedKey) {
			return nil, e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
		}
		if errors.Is(createErr, mapper.ErrLastCommitDuplicated) {
			return nil, e.NewAlreadyExistsError("last commit")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
	}

	return commit, nil
}

func (pushService *PushServiceImpl) PushManifestAndBlobsWithDraft(ctx context.Context, userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet, draftName string) (*model.Commit, e.ResponseError) {
	commit, err := pushService.toCommit(userID, ownerName, repositoryName, fileManifest)
	if err != nil {
		return nil, err
	}
	commit.DraftName = draftName

	// 写入文件
	err = pushService.saveFileManifestAndBlobs(ctx, fileManifest, fileBlobs)
	if err != nil {
		return nil, err
	}

	createErr := pushService.commitMapper.Create(commit)
	if createErr != nil {
		if errors.Is(createErr, mapper.ErrTagAndDraftDuplicated) {
			return nil, e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
		}
		if errors.Is(createErr, mapper.ErrLastCommitDuplicated) {
			return nil, e.NewAlreadyExistsError("last commit")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
	}

	return commit, nil
}

func (pushService *PushServiceImpl) toCommit(userID, ownerName, repositoryName string, fileManifest *manifest.Manifest) (*model.Commit, e.ResponseError) {
	// 获取user
	user, err := pushService.userMapper.FindByUserID(userID)
	if err != nil || user.UserName != ownerName {
		return nil, e.NewPermissionDeniedError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
	}

	// 获取repo
	repository, err := pushService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		return nil, e.NewNotFoundError("repository")
	}

	commitID := uuid.NewString()

	// 生成file blobs
	modelBlobs := make([]*model.FileBlob, 0, len(fileManifest.Paths()))
	_ = fileManifest.Range(func(path string, digest manifest.Digest) error {
		modelBlobs = append(modelBlobs, &model.FileBlob{
			Digest:   digest.Hex(),
			CommitID: commitID,
			FileName: path,
		})

		return nil
	})
	fileManifestBlob, err := fileManifest.Blob()
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}
	modelFileManifest := &model.FileManifest{
		ID:       0,
		Digest:   fileManifestBlob.Digest().Hex(),
		CommitID: commitID,
	}

	commit := &model.Commit{
		UserID:         user.UserID,
		UserName:       user.UserName,
		RepositoryID:   repository.RepositoryID,
		RepositoryName: repositoryName,
		CommitID:       commitID,
		CommitName:     security.GenerateCommitName(user.UserName, repositoryName),
		ManifestDigest: fileManifestBlob.Digest().Hex(),
		FileManifest:   modelFileManifest,
		FileBlobs:      modelBlobs,
	}

	return commit, nil
}

func (pushService *PushServiceImpl) saveFileManifestAndBlobs(ctx context.Context, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet) e.ResponseError {
	// 保存file blobs
	err := fileManifest.Range(func(path string, digest manifest.Digest) error {
		blob, ok := fileBlobs.BlobFor(digest.String())
		if !ok {
			return e.NewInvalidArgumentError("file blobs")
		}

		readCloser, err := blob.Open(ctx)
		if err != nil {
			return e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
		}

		// 写入文件
		err = pushService.storageHelper.Store(digest.Hex(), readCloser)
		if err != nil {
			return e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
		}

		return nil
	})
	if err != nil {
		return e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
	}

	// 保存file manifest
	blob, err := fileManifest.Blob()
	if err != nil {
		return e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
	}
	readCloser, err := blob.Open(ctx)
	if err != nil {
		return e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
	}
	err = pushService.storageHelper.Store(blob.Digest().Hex(), readCloser)
	if err != nil {
		return e.NewInternalError(registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure)
	}

	return nil
}

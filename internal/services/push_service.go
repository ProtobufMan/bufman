package services

import (
	"context"
	"errors"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util"
	"github.com/ProtobufMan/bufman/internal/util/manifest"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PushService interface {
	PushManifestAndBlobs(userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet) (*model.Commit, e.ResponseError)
	PushManifestAndBlobsWithTags(userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet, tagNames []string) (*model.Commit, e.ResponseError)
	PushManifestAndBlobsWithDraft(userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet, draftName string) (*model.Commit, e.ResponseError)
}

type PushServiceImpl struct {
	userMapper       mapper.UserMapper
	repositoryMapper mapper.RepositoryMapper
	commitMapper     mapper.CommitMapper
	storageHelper    util.StorageHelper
}

func NewPushService() PushService {
	return &PushServiceImpl{
		userMapper:       &mapper.UserMapperImpl{},
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
		storageHelper:    util.NewStorageHelper(),
	}
}

func (pushService *PushServiceImpl) PushManifestAndBlobs(userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet) (*model.Commit, e.ResponseError) {
	commit, err := pushService.toCommit(userID, ownerName, repositoryName, fileManifest)
	if err != nil {
		return nil, err
	}

	// 写入文件
	err = pushService.saveFileBlobs(fileManifest, fileBlobs)
	if err != nil {
		return nil, err
	}

	// 写入数据库
	createErr := pushService.commitMapper.Create(commit)
	if createErr != nil {
		if errors.Is(createErr, gorm.ErrDuplicatedKey) {
			return nil, e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
		}
		if errors.Is(createErr, mapper.ErrLastCommitDuplicated) {
			return nil, e.NewAlreadyExistsError("last commit")
		}

		return nil, e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
	}

	return commit, nil
}

func (pushService *PushServiceImpl) PushManifestAndBlobsWithTags(userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet, tagNames []string) (*model.Commit, e.ResponseError) {
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
	err = pushService.saveFileBlobs(fileManifest, fileBlobs)
	if err != nil {
		return nil, err
	}

	createErr := pushService.commitMapper.Create(commit)
	if createErr != nil {
		if errors.Is(createErr, mapper.ErrTagAndDraftDuplicated) || errors.Is(createErr, gorm.ErrDuplicatedKey) {
			return nil, e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
		}
		if errors.Is(createErr, mapper.ErrLastCommitDuplicated) {
			return nil, e.NewAlreadyExistsError("last commit")
		}

		return nil, e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
	}

	return commit, nil
}

func (pushService *PushServiceImpl) PushManifestAndBlobsWithDraft(userID, ownerName, repositoryName string, fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet, draftName string) (*model.Commit, e.ResponseError) {
	commit, err := pushService.toCommit(userID, ownerName, repositoryName, fileManifest)
	if err != nil {
		return nil, err
	}
	commit.DraftName = draftName

	// 写入文件
	err = pushService.saveFileBlobs(fileManifest, fileBlobs)
	if err != nil {
		return nil, err
	}

	createErr := pushService.commitMapper.Create(commit)
	if createErr != nil {
		if errors.Is(createErr, mapper.ErrTagAndDraftDuplicated) {
			return nil, e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
		}
		if errors.Is(createErr, mapper.ErrLastCommitDuplicated) {
			return nil, e.NewAlreadyExistsError("last commit")
		}

		return nil, e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
	}

	return commit, nil
}

func (pushService *PushServiceImpl) toCommit(userID, ownerName, repositoryName string, fileManifest *manifest.Manifest) (*model.Commit, e.ResponseError) {
	// 获取user
	user, err := pushService.userMapper.FindByUserID(userID)
	if err != nil || user.UserName != ownerName {
		return nil, e.NewPermissionDeniedError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
	}

	// 获取repo
	repository, err := pushService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		return nil, e.NewNotFoundError("repository")
	}

	commitID := uuid.NewString()

	// 生成文件清单file manifests
	modelFileManifests := make([]*model.FileManifest, 0, len(fileManifest.Paths()))
	_ = fileManifest.Range(func(path string, digest manifest.Digest) error {
		modelFileManifests = append(modelFileManifests, &model.FileManifest{
			Digest:   digest.Hex(),
			CommitID: commitID,
			FileName: path,
		})

		return nil
	})

	commit := &model.Commit{
		UserID:         user.UserID,
		UserName:       user.UserName,
		RepositoryID:   repository.RepositoryID,
		RepositoryName: repositoryName,
		CommitID:       commitID,
		CommitName:     util.GenerateCommitName(user.UserName, repositoryName),
		ManifestDigest: fileManifest.Digest().Hex(),
		FileManifests:  modelFileManifests,
	}

	return commit, nil
}

func (pushService *PushServiceImpl) saveFileBlobs(fileManifest *manifest.Manifest, fileBlobs *manifest.BlobSet) e.ResponseError {
	err := fileManifest.Range(func(path string, digest manifest.Digest) error {
		blob, ok := fileBlobs.BlobFor(digest.Hex())
		if !ok {
			return e.NewInvalidArgumentError("file blobs")
		}

		readCloser, err := blob.Open(context.Background())
		if err != nil {
			return e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
		}

		// 写入文件
		err = pushService.storageHelper.Store(digest.Hex(), readCloser)
		if err != nil {
			return e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
		}

		return nil
	})

	respErr, ok := err.(e.ResponseError)
	if !ok {
		return e.NewInternalError(registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure)
	}
	return respErr
}

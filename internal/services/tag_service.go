package services

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/e"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/validity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagService interface {
	CreateRepositoryTag(userID, repositoryID, TagName, commitName string) (*model.Tag, e.ResponseError)
	ListRepositoryTags(repositoryID string, offset, limit int, reverse bool) (model.Tags, e.ResponseError)
	ListRepositoryTagsWithUserID(userID, repositoryID string, offset, limit int, reverse bool) (model.Tags, e.ResponseError)
}

func NewTagService() TagService {
	return &TagServiceImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
		tagMapper:        &mapper.TagMapperImpl{},
		validator:        validity.NewValidator(),
	}
}

type TagServiceImpl struct {
	repositoryMapper mapper.RepositoryMapper
	commitMapper     mapper.CommitMapper
	tagMapper        mapper.TagMapper
	validator        validity.Validator
}

func (tagService *TagServiceImpl) CreateRepositoryTag(userID, repositoryID, TagName, commitName string) (*model.Tag, e.ResponseError) {
	// 查询用户是否是repo的owner
	repository, err := tagService.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if err != nil {
			return nil, e.NewNotFoundError("repository")
		}
	}

	if repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryTagServiceCreateRepositoryTagProcedure)
	}

	// 查询commitName
	commit, err := tagService.commitMapper.FindByRepositoryIDAndCommitName(repositoryID, commitName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("commit")
		}
	}

	tag := &model.Tag{
		UserID:       userID,
		UserName:     repository.UserName,
		RepositoryID: repositoryID,
		CommitID:     commit.CommitID,
		CommitName:   commitName,
		TagID:        uuid.NewString(),
		TagName:      TagName,
	}
	err = tagService.tagMapper.Create(tag)
	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryTagServiceCreateRepositoryTagProcedure)
	}

	return tag, nil
}

func (tagService *TagServiceImpl) ListRepositoryTags(repositoryID string, offset, limit int, reverse bool) (model.Tags, e.ResponseError) {
	// 查询用户是否是repo的owner
	repository, err := tagService.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if err != nil {
			return nil, e.NewNotFoundError("repository")
		}
	}
	if repository.Visibility != uint8(registryv1alpha.Visibility_VISIBILITY_PUBLIC) {
		return nil, e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryTagServiceListRepositoryTagsProcedure)
	}

	tags, err := tagService.tagMapper.FindPageByRepositoryID(repositoryID, limit, offset, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryTagServiceListRepositoryTagsProcedure)
	}

	return tags, nil
}

func (tagService *TagServiceImpl) ListRepositoryTagsWithUserID(userID, repositoryID string, offset, limit int, reverse bool) (model.Tags, e.ResponseError) {
	// 查询用户是否是repo的owner
	repository, err := tagService.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if err != nil {
			return nil, e.NewNotFoundError("repository")
		}
	}
	if repository.Visibility != uint8(registryv1alpha.Visibility_VISIBILITY_PUBLIC) && repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryTagServiceListRepositoryTagsProcedure)
	}

	tags, err := tagService.tagMapper.FindPageByRepositoryID(repositoryID, limit, offset, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryTagServiceListRepositoryTagsProcedure)
	}

	return tags, nil
}

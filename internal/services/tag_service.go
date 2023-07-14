package services

import (
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/validity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagService interface {
	CreateRepositoryTag(repositoryID, TagName, commitName string) (*model.Tag, e.ResponseError)
	ListRepositoryTags(repositoryID string, offset, limit int, reverse bool) (model.Tags, e.ResponseError)
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

func (tagService *TagServiceImpl) CreateRepositoryTag(repositoryID, TagName, commitName string) (*model.Tag, e.ResponseError) {
	// 查询commitName
	commit, err := tagService.commitMapper.FindByRepositoryIDAndCommitName(repositoryID, commitName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("commit")
		}
	}

	tag := &model.Tag{
		UserID:       commit.UserID,
		UserName:     commit.UserName,
		RepositoryID: repositoryID,
		CommitID:     commit.CommitID,
		CommitName:   commitName,
		TagID:        uuid.NewString(),
		TagName:      TagName,
	}
	err = tagService.tagMapper.Create(tag)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryTagServiceCreateRepositoryTagProcedure)
	}

	return tag, nil
}

func (tagService *TagServiceImpl) ListRepositoryTags(repositoryID string, offset, limit int, reverse bool) (model.Tags, e.ResponseError) {
	tags, err := tagService.tagMapper.FindPageByRepositoryID(repositoryID, limit, offset, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryTagServiceListRepositoryTagsProcedure)
	}

	return tags, nil
}

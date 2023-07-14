package services

import (
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type RepositoryService interface {
	GetRepository(repositoryID string) (*model.Repository, e.ResponseError)
	GetRepositoryByUserNameAndRepositoryName(userName, repositoryName string) (*model.Repository, e.ResponseError)
	GetRepositoryCounts(repositoryID string) (*model.RepositoryCounts, e.ResponseError)
	ListRepositories(offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	ListUserRepositories(userID string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	ListRepositoriesUserCanAccess(userID string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	CreateRepositoryByUserNameAndRepositoryName(userID, userName, repositoryName string, visibility registryv1alpha1.Visibility) (*model.Repository, e.ResponseError)
	DeleteRepository(repositoryID string) e.ResponseError
	DeleteRepositoryByUserNameAndRepositoryName(userName, repositoryName string) e.ResponseError
	DeprecateRepositoryByName(ownerName, repositoryName, deprecateMsg string) (*model.Repository, e.ResponseError)
	UndeprecateRepositoryByName(ownerName, repositoryName string) (*model.Repository, e.ResponseError)
	UpdateRepositorySettingsByName(ownerName, repositoryName string, visibility registryv1alpha1.Visibility, description string) e.ResponseError
}

type RepositoryServiceImpl struct {
	repositoryMapper mapper.RepositoryMapper
	userMapper       mapper.UserMapper
	commitMapper     mapper.CommitMapper
	tagMapper        mapper.TagMapper
}

func NewRepositoryService() RepositoryService {
	return &RepositoryServiceImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		userMapper:       &mapper.UserMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
		tagMapper:        &mapper.TagMapperImpl{},
	}
}

func (repositoryService *RepositoryServiceImpl) GetRepository(repositoryID string) (*model.Repository, e.ResponseError) {
	// 查询
	repository, err := repositoryService.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("repository")
		}
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceGetRepositoryProcedure)
	}

	return repository, nil
}

func (repositoryService *RepositoryServiceImpl) GetRepositoryByUserNameAndRepositoryName(userName, repositoryName string) (*model.Repository, e.ResponseError) {
	// 查询
	repository, err := repositoryService.repositoryMapper.FindByUserNameAndRepositoryName(userName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("repository")
		}
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceGetRepositoryByFullNameProcedure)
	}

	return repository, nil
}

func (repositoryService *RepositoryServiceImpl) GetRepositoryCounts(repositoryID string) (*model.RepositoryCounts, e.ResponseError) {
	// 忽略错误
	draftCounts, _ := repositoryService.commitMapper.GetDraftCountsByRepositoryID(repositoryID)
	tagCounts, _ := repositoryService.tagMapper.GetCountsByRepositoryID(repositoryID)
	repositoryCounts := &model.RepositoryCounts{
		TagsCount:   tagCounts,
		DraftsCount: draftCounts,
	}

	return repositoryCounts, nil
}

func (repositoryService *RepositoryServiceImpl) ListRepositories(offset, limit int, reverse bool) (model.Repositories, e.ResponseError) {
	repositories, err := repositoryService.repositoryMapper.FindPage(offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceListRepositoriesProcedure)
	}

	return repositories, nil
}

func (repositoryService *RepositoryServiceImpl) ListUserRepositories(userID string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError) {
	repositories, err := repositoryService.repositoryMapper.FindPageByUserID(userID, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceListUserRepositoriesProcedure)
	}

	return repositories, nil
}

func (repositoryService *RepositoryServiceImpl) ListRepositoriesUserCanAccess(userID string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError) {
	repositories, err := repositoryService.repositoryMapper.FindAccessiblePageByUserID(userID, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceListRepositoriesUserCanAccessProcedure)
	}

	return repositories, nil
}

func (repositoryService *RepositoryServiceImpl) CreateRepositoryByUserNameAndRepositoryName(userID, userName, repositoryName string, visibility registryv1alpha1.Visibility) (*model.Repository, e.ResponseError) {
	// 查询用户
	user, err := repositoryService.userMapper.FindByUserName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("user")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceCreateRepositoryByFullNameProcedure)
	}
	if user.UserID != userID {
		return nil, e.NewPermissionDeniedError(registryv1alpha1connect.RepositoryServiceCreateRepositoryByFullNameProcedure)
	}

	// 创建repo
	repository := &model.Repository{
		UserID:         user.UserID,
		UserName:       user.UserName,
		RepositoryID:   uuid.NewString(),
		RepositoryName: repositoryName,
		Visibility:     uint8(visibility),
	}

	err = repositoryService.repositoryMapper.Create(repository)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, e.NewAlreadyExistsError("repository")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceCreateRepositoryByFullNameProcedure)
	}

	return repository, nil
}

func (repositoryService *RepositoryServiceImpl) DeleteRepository(repositoryID string) e.ResponseError {
	// 删除
	err := repositoryService.repositoryMapper.DeleteByRepositoryID(repositoryID)
	if err != nil {
		return e.NewInternalError(registryv1alpha1connect.RepositoryServiceDeleteRepositoryProcedure)
	}

	return nil
}

func (repositoryService *RepositoryServiceImpl) DeleteRepositoryByUserNameAndRepositoryName(userName, repositoryName string) e.ResponseError {
	err := repositoryService.repositoryMapper.DeleteByUserNameAndRepositoryName(userName, repositoryName)
	if err != nil {
		return e.NewInternalError(registryv1alpha1connect.RepositoryServiceDeleteRepositoryByFullNameProcedure)
	}

	return nil
}

func (repositoryService *RepositoryServiceImpl) DeprecateRepositoryByName(ownerName, repositoryName, deprecateMsg string) (*model.Repository, e.ResponseError) {
	// 修改数据库
	updatedRepository := &model.Repository{
		Deprecated:     true,
		DeprecationMsg: deprecateMsg,
	}
	err := repositoryService.repositoryMapper.UpdateByUserNameAndRepositoryName(ownerName, repositoryName, updatedRepository)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceDeprecateRepositoryByNameProcedure)
	}

	return updatedRepository, nil
}

func (repositoryService *RepositoryServiceImpl) UndeprecateRepositoryByName(ownerName, repositoryName string) (*model.Repository, e.ResponseError) {
	// 修改数据库
	updatedRepository := &model.Repository{
		Deprecated: false,
	}
	err := repositoryService.repositoryMapper.UpdateByUserNameAndRepositoryName(ownerName, repositoryName, updatedRepository)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryServiceUndeprecateRepositoryByNameProcedure)
	}

	return updatedRepository, nil
}

func (repositoryService *RepositoryServiceImpl) UpdateRepositorySettingsByName(ownerName, repositoryName string, visibility registryv1alpha1.Visibility, description string) e.ResponseError {
	// 修改数据库
	updatedRepository := &model.Repository{
		Visibility:  uint8(visibility),
		Description: description,
	}
	err := repositoryService.repositoryMapper.UpdateByUserNameAndRepositoryName(ownerName, repositoryName, updatedRepository)
	if err != nil {
		return e.NewInternalError(registryv1alpha1connect.RepositoryServiceUpdateRepositorySettingsByNameProcedure)
	}

	return nil
}

func SplitFullName(fullName string) (userName, repositoryName string, ok bool) {
	split := strings.SplitN(fullName, "/", 2)
	if len(split) != 2 {
		return
	}
	userName, repositoryName = split[0], split[1]
	ok = userName != "" && repositoryName != ""
	return
}

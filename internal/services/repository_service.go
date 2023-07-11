package services

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/e"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type RepositoryService interface {
	GetRepository(repositoryID string) (*model.Repository, e.ResponseError)
	GetRepositoryByFullName(fullName string) (*model.Repository, e.ResponseError)
	GetRepositoryCounts(repositoryID string) (*model.RepositoryCounts, e.ResponseError)
	ListRepositories(offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	ListUserRepositories(userID string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	ListRepositoriesUserCanAccess(userID string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	CreateRepositoryByFullName(userID string, fullName string, visibility registryv1alpha.Visibility) (*model.Repository, e.ResponseError)
	DeleteRepository(userID, repositoryID string) e.ResponseError
	DeleteRepositoryByFullName(userID, fullName string) e.ResponseError
	DeprecateRepositoryByName(userID, ownerName, repositoryName, deprecateMsg string) (*model.Repository, e.ResponseError)
	UndeprecateRepositoryByName(userID, ownerName, repositoryName string) (*model.Repository, e.ResponseError)
	UpdateRepositorySettingsByName(userID, ownerName, repositoryName string, visibility registryv1alpha.Visibility, description string) e.ResponseError
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
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceGetRepositoryProcedure)
	}

	return repository, nil
}

func (repositoryService *RepositoryServiceImpl) GetRepositoryByFullName(fullName string) (*model.Repository, e.ResponseError) {
	// 查询
	userName, repositoryName, ok := SplitFullName(fullName)
	if !ok {
		return nil, e.NewNotFoundError("repository")
	}
	repository, err := repositoryService.repositoryMapper.FindByUserNameAndRepositoryName(userName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("repository")
		}
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceGetRepositoryByFullNameProcedure)
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
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceListRepositoriesProcedure)
	}

	return repositories, nil
}

func (repositoryService *RepositoryServiceImpl) ListUserRepositories(userID string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError) {
	repositories, err := repositoryService.repositoryMapper.FindPageByUserID(userID, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceListUserRepositoriesProcedure)
	}

	return repositories, nil
}

func (repositoryService *RepositoryServiceImpl) ListRepositoriesUserCanAccess(userID string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError) {
	repositories, err := repositoryService.repositoryMapper.FindAccessiblePageByUserID(userID, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceListRepositoriesUserCanAccessProcedure)
	}

	return repositories, nil
}

func (repositoryService *RepositoryServiceImpl) CreateRepositoryByFullName(userID string, fullName string, visibility registryv1alpha.Visibility) (*model.Repository, e.ResponseError) {
	userName, repositoryName, ok := SplitFullName(fullName)
	if !ok {
		return nil, e.NewInvalidArgumentError("full name")
	}

	// 查询用户
	user, err := repositoryService.userMapper.FindByUserName(userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("user")
		}

		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceCreateRepositoryByFullNameProcedure)
	}
	if user.UserID != userID {
		return nil, e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryServiceCreateRepositoryByFullNameProcedure)
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

		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceCreateRepositoryByFullNameProcedure)
	}

	return repository, nil
}

func (repositoryService *RepositoryServiceImpl) DeleteRepository(userID, repositoryID string) e.ResponseError {
	// 查询repository，检查是否可以删除
	repository, err := repositoryService.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("repository")
		}

		return e.NewInternalError(registryv1alphaconnect.RepositoryServiceDeleteRepositoryProcedure)
	}

	if repository.UserID != userID {
		return e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryServiceDeleteRepositoryProcedure)
	}

	// 删除
	err = repositoryService.repositoryMapper.DeleteByRepositoryID(repositoryID)
	if err != nil {
		return e.NewInternalError(registryv1alphaconnect.RepositoryServiceDeleteRepositoryProcedure)
	}

	return nil
}

func (repositoryService *RepositoryServiceImpl) DeleteRepositoryByFullName(userID, fullName string) e.ResponseError {
	userName, repositoryName, ok := SplitFullName(fullName)
	if !ok {
		return e.NewInvalidArgumentError("full name")
	}

	// 查询repository，检查是否可以删除
	repository, err := repositoryService.repositoryMapper.FindByUserNameAndRepositoryName(userName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("repository")
		}

		return e.NewInternalError(registryv1alphaconnect.RepositoryServiceDeleteRepositoryByFullNameProcedure)
	}

	// 不能删除别人的repo
	if repository.UserID != userID {
		return e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryServiceDeleteRepositoryByFullNameProcedure)
	}

	err = repositoryService.repositoryMapper.DeleteByUserNameAndRepositoryName(userName, repositoryName)
	if err != nil {
		return e.NewInternalError(registryv1alphaconnect.RepositoryServiceDeleteRepositoryByFullNameProcedure)
	}

	return nil
}

func (repositoryService *RepositoryServiceImpl) DeprecateRepositoryByName(userID, ownerName, repositoryName, deprecateMsg string) (*model.Repository, e.ResponseError) {
	// 查询repository，检查是否可以修改
	repository, err := repositoryService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("repository")
		}

		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceDeprecateRepositoryByNameProcedure)
	}

	// 不能丢弃别人的repo
	if repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryServiceDeprecateRepositoryByNameProcedure)
	}

	// 修改数据库
	updatedRepository := &model.Repository{
		Deprecated:     true,
		DeprecationMsg: deprecateMsg,
	}
	err = repositoryService.repositoryMapper.UpdateByUserNameAndRepositoryName(ownerName, repositoryName, updatedRepository)
	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceDeprecateRepositoryByNameProcedure)
	}

	return updatedRepository, nil
}

func (repositoryService *RepositoryServiceImpl) UndeprecateRepositoryByName(userID, ownerName, repositoryName string) (*model.Repository, e.ResponseError) {
	// 查询repository，检查是否可以修改
	repository, err := repositoryService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("repository")
		}

		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceUndeprecateRepositoryByNameProcedure)
	}

	if repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryServiceUndeprecateRepositoryByNameProcedure)
	}

	// 修改数据库
	updatedRepository := &model.Repository{
		Deprecated: false,
	}
	err = repositoryService.repositoryMapper.UpdateByUserNameAndRepositoryName(ownerName, repositoryName, updatedRepository)
	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryServiceUndeprecateRepositoryByNameProcedure)
	}

	return updatedRepository, nil
}

func (repositoryService *RepositoryServiceImpl) UpdateRepositorySettingsByName(userID, ownerName, repositoryName string, visibility registryv1alpha.Visibility, description string) e.ResponseError {
	// 查询repository，检查是否可以修改
	repository, err := repositoryService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("repository")
		}

		return e.NewInternalError(registryv1alphaconnect.RepositoryServiceUpdateRepositorySettingsByNameProcedure)
	}

	if repository.UserID != userID {
		return e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryServiceUpdateRepositorySettingsByNameProcedure)
	}

	// 修改数据库
	updatedRepository := &model.Repository{
		Visibility:  uint8(visibility),
		Description: description,
	}
	err = repositoryService.repositoryMapper.UpdateByUserNameAndRepositoryName(ownerName, repositoryName, updatedRepository)
	if err != nil {
		return e.NewInternalError(registryv1alphaconnect.RepositoryServiceUpdateRepositorySettingsByNameProcedure)
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

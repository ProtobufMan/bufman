package services

import (
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/gorm"
)

// AuthorizationService 用户权限验证
type AuthorizationService interface {
	CheckRepositoryCanAccess(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) // 检查user id用户是否可以访问repo
	CheckRepositoryCanAccessByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError)
	CheckRepositoryCanEdit(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) // 检查user是否可以修改repo
	CheckRepositoryCanEditByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError)
	CheckRepositoryCanDelete(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) // 检查用户是否可以删除repo
	CheckRepositoryCanDeleteByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError)
}

func NewAuthorizationService() AuthorizationService {
	return &AuthorizationServiceImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
	}
}

type AuthorizationServiceImpl struct {
	repositoryMapper mapper.RepositoryMapper
}

func (authorizationService *AuthorizationServiceImpl) CheckRepositoryCanAccess(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := authorizationService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError(fmt.Sprintf("repository [name=%s/%s]", ownerName, repositoryName))
		}

		return nil, e.NewInternalError(procedure)
	}

	if registryv1alpha1.Visibility(repository.Visibility) != registryv1alpha1.Visibility_VISIBILITY_PUBLIC && repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(fmt.Sprintf("repository [name=%s/%s]", ownerName, repositoryName))
	}

	return repository, nil
}

func (authorizationService *AuthorizationServiceImpl) CheckRepositoryCanAccessByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := authorizationService.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError(fmt.Sprintf("repository [id=%s]", repositoryID))
		}

		return nil, e.NewInternalError(procedure)
	}

	if registryv1alpha1.Visibility(repository.Visibility) != registryv1alpha1.Visibility_VISIBILITY_PUBLIC && repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(fmt.Sprintf("repository [id=%s]", repositoryID))
	}

	return repository, nil
}

func (authorizationService *AuthorizationServiceImpl) CheckRepositoryCanEdit(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := authorizationService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError(fmt.Sprintf("repository [name=%s/%s]", ownerName, repositoryName))
		}

		return nil, e.NewInternalError(procedure)
	}

	// 只有所属用户才能修改
	if repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(fmt.Sprintf("repository [name=%s/%s]", ownerName, repositoryName))
	}

	return repository, nil
}

func (authorizationService *AuthorizationServiceImpl) CheckRepositoryCanEditByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := authorizationService.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError(fmt.Sprintf("repository [id=%s]", repositoryID))
		}

		return nil, e.NewInternalError(procedure)
	}

	// 只有所属用户才能修改
	if repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(fmt.Sprintf("repository [id=%s]", repositoryID))
	}

	return repository, nil
}

func (authorizationService *AuthorizationServiceImpl) CheckRepositoryCanDelete(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := authorizationService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError(fmt.Sprintf("repository [name=%s/%s]", ownerName, repositoryName))
		}

		return nil, e.NewInternalError(procedure)
	}

	// 只有所属用户才能修改
	if repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(fmt.Sprintf("repository [name=%s/%s]", ownerName, repositoryName))
	}

	return repository, nil
}

func (authorizationService *AuthorizationServiceImpl) CheckRepositoryCanDeleteByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := authorizationService.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError(fmt.Sprintf("repository [id=%s]", repositoryID))
		}

		return nil, e.NewInternalError(procedure)
	}

	// 只有所属用户才能修改
	if repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(fmt.Sprintf("repository [id=%s]", repositoryID))
	}

	return repository, nil
}

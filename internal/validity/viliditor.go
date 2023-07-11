package validity

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/e"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/gorm"
)

type Validator interface {
	CheckRepositoryCanAccess(ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError)
	CheckRepositoryCanAccessWithUserID(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError)
}

func NewValidator() Validator {
	return &ValidatorImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
	}
}

type ValidatorImpl struct {
	repositoryMapper mapper.RepositoryMapper
}

func (validator *ValidatorImpl) CheckRepositoryCanAccess(ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := validator.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("repository")
		}

		return nil, e.NewInternalError(procedure)
	}

	if registryv1alpha.Visibility(repository.Visibility) != registryv1alpha.Visibility_VISIBILITY_PUBLIC {
		return nil, e.NewPermissionDeniedError(procedure)
	}

	return repository, nil
}

func (validator *ValidatorImpl) CheckRepositoryCanAccessWithUserID(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := validator.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("repository")
		}

		return nil, e.NewInternalError(procedure)
	}

	if registryv1alpha.Visibility(repository.Visibility) != registryv1alpha.Visibility_VISIBILITY_PUBLIC && repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(procedure)
	}

	return repository, nil
}

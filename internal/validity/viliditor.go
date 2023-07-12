package validity

import (
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

type Validator interface {
	CheckUserName(username string) e.ResponseError                                                                     // 检查用户名合法性
	CheckPassword(password string) e.ResponseError                                                                     // 检查密码合法性
	CheckRepositoryName(repositoryName string) e.ResponseError                                                         // 检查repo name合法性
	CheckTagName(tagName string) e.ResponseError                                                                       // 检查tag name合法性
	CheckDraftName(draftName string) e.ResponseError                                                                   // 检查draft name合法性
	CheckPageSize(pageSize uint32) e.ResponseError                                                                     // 检查page size合法性
	SplitFullName(fullName string) (userName, repositoryName string, respErr e.ResponseError)                          // 分割full name
	CheckRepositoryCanAccess(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) // 检查user id用户是否可以访问repo
}

func NewValidator() Validator {
	return &ValidatorImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
	}
}

type ValidatorImpl struct {
	repositoryMapper mapper.RepositoryMapper
}

func (validator *ValidatorImpl) CheckUserName(username string) e.ResponseError {
	err := validator.doCheckByLengthAndPattern(username, constant.MinUserNameLength, constant.MaxUserNameLength, constant.UserNamePattern)
	if err != nil {
		return e.NewInvalidArgumentError("username:" + err.Error())
	}

	return nil
}

func (validator *ValidatorImpl) CheckPassword(password string) e.ResponseError {
	err := validator.doCheckByLengthAndPattern(password, constant.MinPasswordLength, constant.MaxPasswordLength, constant.PasswordPattern)
	if err != nil {
		return e.NewInvalidArgumentError("password:" + err.Error())
	}

	return nil
}

func (validator *ValidatorImpl) CheckRepositoryName(repositoryName string) e.ResponseError {
	err := validator.doCheckByLengthAndPattern(repositoryName, constant.MinRepositoryNameLength, constant.MaxRepositoryNameLength, constant.RepositoryNamePattern)
	if err != nil {
		return e.NewInvalidArgumentError("repo name:" + err.Error())
	}

	return nil
}

func (validator *ValidatorImpl) CheckTagName(tagName string) e.ResponseError {
	err := validator.doCheckByLengthAndPattern(tagName, constant.MinTagLength, constant.MaxTagLength, constant.TagPattern)
	if err != nil {
		return e.NewInvalidArgumentError("tag name:" + err.Error())
	}

	return nil
}

func (validator *ValidatorImpl) CheckDraftName(draftName string) e.ResponseError {
	if draftName == constant.DefaultBranch { // draft name不能为 main
		return e.NewInvalidArgumentError(fmt.Sprintf("draft (can not be '%v')", constant.DefaultBranch))
	}

	err := validator.doCheckByLengthAndPattern(draftName, constant.MinDraftLength, constant.MaxDraftLength, constant.DraftPattern)
	if err != nil {
		return e.NewInvalidArgumentError("draft name:" + err.Error())
	}

	return nil
}

func (validator *ValidatorImpl) CheckPageSize(pageSize uint32) e.ResponseError {
	if pageSize < constant.MinPageSize || pageSize > constant.MaxPageSize {
		return e.NewInvalidArgumentError(fmt.Sprintf("page size: length is limited between %v and %v", constant.MinPageSize, constant.MaxPageSize))
	}

	return nil
}

func (validator *ValidatorImpl) SplitFullName(fullName string) (userName, repositoryName string, respErr e.ResponseError) {
	split := strings.SplitN(fullName, "/", 2)
	if len(split) != 2 {
		respErr = e.NewInvalidArgumentError("full name")
		return
	}

	ok := split[0] != "" && split[1] != ""
	if !ok {
		respErr = e.NewInvalidArgumentError("full name")
		return
	}

	userName, repositoryName = split[0], split[1]
	return userName, repositoryName, nil
}

func (validator *ValidatorImpl) doCheckByLengthAndPattern(str string, minLength, maxLength int, pattern string) error {
	// 长度检查
	if len(str) < minLength || len(str) > maxLength {
		return fmt.Errorf("length is limited between %v and %v", minLength, maxLength)
	}

	// 正则匹配
	match, _ := regexp.MatchString(pattern, str)
	if !match {
		return fmt.Errorf("pattern dont math %s", pattern)
	}

	return nil
}

func (validator *ValidatorImpl) CheckRepositoryCanAccess(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) {
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

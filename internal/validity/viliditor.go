package validity

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufconfig"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/buflock"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	modulev1alpha "github.com/ProtobufMan/bufman/internal/gen/module/v1alpha"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/gorm"
	"io"
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
	CheckRepositoryCanAccessByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError)
	CheckRepositoryCanEdit(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) // 检查user是否可以修改repo
	CheckRepositoryCanEditByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError)
	CheckRepositoryCanDelete(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) // 检查用户是否可以删除repo
	CheckRepositoryCanDeleteByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError)
	CheckManifestAndBlobs(ctx context.Context, protoManifest *modulev1alpha.Blob, protoBlobs []*modulev1alpha.Blob) (*manifest.Manifest, *manifest.BlobSet, *bufconfig.Config, e.ResponseError)
}

func NewValidator() Validator {
	return &ValidatorImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
	}
}

type ValidatorImpl struct {
	repositoryMapper mapper.RepositoryMapper
	commitMapper     mapper.CommitMapper
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
			return nil, e.NewNotFoundError(fmt.Sprintf("repository [name=%s/%s]", ownerName, repositoryName))
		}

		return nil, e.NewInternalError(procedure)
	}

	if registryv1alpha.Visibility(repository.Visibility) != registryv1alpha.Visibility_VISIBILITY_PUBLIC && repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(fmt.Sprintf("repository [name=%s/%s]", ownerName, repositoryName))
	}

	return repository, nil
}

func (validator *ValidatorImpl) CheckRepositoryCanAccessByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := validator.repositoryMapper.FindByRepositoryID(repositoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError(fmt.Sprintf("repository [id=%s]", repositoryID))
		}

		return nil, e.NewInternalError(procedure)
	}

	if registryv1alpha.Visibility(repository.Visibility) != registryv1alpha.Visibility_VISIBILITY_PUBLIC && repository.UserID != userID {
		return nil, e.NewPermissionDeniedError(fmt.Sprintf("repository [id=%s]", repositoryID))
	}

	return repository, nil
}

func (validator *ValidatorImpl) CheckRepositoryCanEdit(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := validator.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
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

func (validator *ValidatorImpl) CheckRepositoryCanEditByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := validator.repositoryMapper.FindByRepositoryID(repositoryID)
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

func (validator *ValidatorImpl) CheckRepositoryCanDelete(userID, ownerName, repositoryName, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := validator.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
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

func (validator *ValidatorImpl) CheckRepositoryCanDeleteByID(userID, repositoryID, procedure string) (*model.Repository, e.ResponseError) {
	repository, err := validator.repositoryMapper.FindByRepositoryID(repositoryID)
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

func (validator *ValidatorImpl) CheckManifestAndBlobs(ctx context.Context, protoManifest *modulev1alpha.Blob, protoBlobs []*modulev1alpha.Blob) (*manifest.Manifest, *manifest.BlobSet, *bufconfig.Config, e.ResponseError) {
	// 读取文件清单
	reader := bytes.NewReader(protoManifest.Content)
	fileManifest, err := manifest.NewFromReader(reader)
	if err != nil {
		return nil, nil, nil, e.NewInvalidArgumentError(err.Error())
	}
	if fileManifest.Empty() {
		// 不允许上次空的commit
		return nil, nil, nil, e.NewInvalidArgumentError("no files")
	}

	// 读取文件列表
	blobs := make([]manifest.Blob, 0, len(protoBlobs))
	for i := 0; i < len(protoBlobs); i++ {
		reader = bytes.NewReader(protoBlobs[i].Content)
		blob, err := manifest.NewMemoryBlobFromReader(reader)
		if err != nil {
			return nil, nil, nil, e.NewInvalidArgumentError(err.Error())
		}
		blobs = append(blobs, blob)
	}

	blobSet, err := manifest.NewBlobSet(ctx, blobs)
	if err != nil {
		return nil, nil, nil, e.NewInvalidArgumentError(err.Error())
	}

	// 检查文件清单和blobs
	var configFileExist bool
	var configFileData []byte
	externalPaths := []string{
		buflock.ExternalConfigFilePath,
		bufmodule.LicenseFilePath,
	}
	externalPaths = append(externalPaths, bufconfig.AllConfigFilePaths...)
	err = fileManifest.Range(func(path string, digest manifest.Digest) error {
		blob, ok := blobSet.BlobFor(digest.String())
		if !ok {
			// 文件清单中有的文件，在file blobs中没有
			return errors.New("check manifest and file blobs failed")
		}

		// 如果遇到配置文件，就记录下来
		for _, configFilePath := range bufconfig.AllConfigFilePaths {
			if path == configFilePath {
				reader, err := blob.Open(ctx)
				if err != nil {
					return err
				}
				configFileData, err = io.ReadAll(reader)
				if err != nil {
					return err
				}
				configFileExist = true
			}
		}

		// 仅仅允许上传.proto、readme、license、配置文件
		unexpected := true
		if !strings.HasSuffix(path, ".proto") {
			for _, externalPath := range externalPaths {
				if path == externalPath {
					unexpected = false
					break
				}
			}
		}
		if unexpected {
			return errors.New("only allow update .proto、readme、license、bufman config file")
		}

		return nil
	})
	if err != nil {
		return nil, nil, nil, e.NewInvalidArgumentError(err.Error())
	}
	if !configFileExist {
		// 不存在配置文件
		return nil, nil, nil, e.NewInvalidArgumentError("no config file")
	}

	// 生成Config，并验证其中的依赖关系
	bufConfig, err := bufconfig.GetConfigForData(ctx, configFileData)
	if err != nil {
		// 无法解析配置文件
		return nil, nil, nil, e.NewInvalidArgumentError(err.Error())
	}

	return fileManifest, blobSet, bufConfig, nil
}

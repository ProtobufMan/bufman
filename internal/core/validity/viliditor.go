package validity

import (
	"context"
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufconfig"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/buflock"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmanifest"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule"
	modulev1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/module/v1alpha1"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/docker"
	"github.com/ProtobufMan/bufman/internal/e"
	"golang.org/x/mod/semver"
	"regexp"
	"strings"
)

type Validator interface {
	CheckUserName(username string) e.ResponseError                                            // 检查用户名合法性
	CheckPassword(password string) e.ResponseError                                            // 检查密码合法性
	CheckRepositoryName(repositoryName string) e.ResponseError                                // 检查repo name合法性
	CheckTagName(tagName string) e.ResponseError                                              // 检查tag name合法性
	CheckPluginName(pluginName string) e.ResponseError                                        // 检查插件名称合法性
	CheckDockerRepoName(dockerRepoName string) e.ResponseError                                // 检查docker repo name合法性
	CheckQuery(query string) e.ResponseError                                                  // 检查query字符串的合法性
	CheckVersion(version string) e.ResponseError                                              // 检查版本号是否合法
	CheckDraftName(draftName string) e.ResponseError                                          // 检查draft name合法性
	CheckPageSize(pageSize uint32) e.ResponseError                                            // 检查page size合法性
	SplitFullName(fullName string) (userName, repositoryName string, respErr e.ResponseError) // 分割full name

	// CheckRegistryAuth 检查是否可以登录registry
	CheckRegistryAuth(ctx context.Context, address, username, password string) e.ResponseError

	// CheckManifestAndBlobs 检查上传的文件是否合法
	CheckManifestAndBlobs(ctx context.Context, protoManifest *modulev1alpha1.Blob, protoBlobs []*modulev1alpha1.Blob) (*manifest.Manifest, *manifest.BlobSet, e.ResponseError)
}

func NewValidator() Validator {
	return &ValidatorImpl{}
}

type ValidatorImpl struct {
}

func (validator *ValidatorImpl) CheckRegistryAuth(ctx context.Context, address, username, password string) e.ResponseError {
	// 验证username和password是否合法
	d, err := docker.NewDockerClient(address, username, password)
	if err != nil {
		return e.NewInternalError(err.Error())
	}
	_, err = d.RegistryLogin(ctx)
	if err != nil {
		return e.NewInternalError(err.Error())
	}

	return nil
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

func (validator *ValidatorImpl) CheckPluginName(pluginName string) e.ResponseError {
	err := validator.doCheckByLengthAndPattern(pluginName, constant.MinPluginLength, constant.MaxPluginLength, constant.PluginNamePattern)
	if err != nil {
		return e.NewInvalidArgumentError("plugin name:" + err.Error())
	}

	return nil
}

func (validator *ValidatorImpl) CheckDockerRepoName(dockerRepoName string) e.ResponseError {
	err := validator.doCheckByLengthAndPattern(dockerRepoName, constant.MinDockerRepoNameLength, constant.MaxDockerRepoNameLength, constant.DockerRepoNamePattern)
	if err != nil {
		return e.NewInvalidArgumentError("docker repo name:" + err.Error())
	}

	return nil
}

func (validator *ValidatorImpl) CheckQuery(query string) e.ResponseError {
	err := validator.doCheckByLengthAndPattern(query, constant.MinQueryLength, constant.MaxQueryLength, constant.QueryPattern)
	if err != nil {
		return e.NewInvalidArgumentError("query string:" + err.Error())
	}

	return nil
}

func (validator *ValidatorImpl) CheckVersion(version string) e.ResponseError {
	if !semver.IsValid(version) {
		// 版本号不合法
		return e.NewInvalidArgumentError("version")
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

func (validator *ValidatorImpl) CheckManifestAndBlobs(ctx context.Context, protoManifest *modulev1alpha1.Blob, protoBlobs []*modulev1alpha1.Blob) (*manifest.Manifest, *manifest.BlobSet, e.ResponseError) {
	// 读取文件清单
	fileManifest, err := bufmanifest.NewManifestFromProto(ctx, protoManifest)
	if err != nil {
		return nil, nil, e.NewInvalidArgumentError(err.Error())
	}
	if fileManifest.Empty() {
		// 不允许上次空的commit
		return nil, nil, e.NewInvalidArgumentError("no files")
	}

	// 读取文件列表
	blobSet, err := bufmanifest.NewBlobSetFromProto(ctx, protoBlobs)
	if err != nil {
		return nil, nil, e.NewInvalidArgumentError(err.Error())
	}

	// 检查文件清单和blobs
	externalPaths := []string{
		buflock.ExternalConfigFilePath,
		bufmodule.LicenseFilePath,
	}
	externalPaths = append(externalPaths, bufmodule.AllDocumentationPaths...)
	externalPaths = append(externalPaths, bufconfig.AllConfigFilePaths...)
	err = fileManifest.Range(func(path string, digest manifest.Digest) error {
		_, ok := blobSet.BlobFor(digest.String())
		if !ok {
			// 文件清单中有的文件，在file blobs中没有
			return errors.New("check manifest and file blobs failed")
		}

		// 仅仅允许上传.proto、readme、license、配置文件
		if !strings.HasSuffix(path, ".proto") {
			unexpected := true
			for _, externalPath := range externalPaths {
				if path == externalPath {
					unexpected = false
					break
				}
			}

			if unexpected {
				return errors.New("only allow update .proto、readme、license、bufman config file")
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, e.NewInvalidArgumentError(err.Error())
	}

	return fileManifest, blobSet, nil
}

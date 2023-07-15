package resolve

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
	"github.com/ProtobufMan/bufman/internal/config"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util/storage"
	"gorm.io/gorm"
	"io"
)

type Resolver interface {
	GetBufConfigFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (*bufconfig.Config, e.ResponseError)
	GetBufConfigFromProto(ctx context.Context, protoManifest *modulev1alpha1.Blob, protoBlobs []*modulev1alpha1.Blob) (*bufconfig.Config, e.ResponseError)
	GetAllDependenciesFromBufConfig(ctx context.Context, bufConfig *bufconfig.Config) (model.Commits, e.ResponseError)    // 获取全部依赖
	GetDirectDependenciesFromBufConfig(ctx context.Context, bufConfig *bufconfig.Config) (model.Commits, e.ResponseError) // 获取直接依赖
}

type ResolverImpl struct {
	repositoryMapper mapper.RepositoryMapper
	commitMapper     mapper.CommitMapper
	fileMapper       mapper.FileMapper
	storageHelper    storage.StorageHelper
}

func NewResolver() Resolver {
	return &ResolverImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
		fileMapper:       &mapper.FileMapperImpl{},
		storageHelper:    storage.NewStorageHelper(),
	}
}

func (resolver *ResolverImpl) GetAllDependenciesFromBufConfig(ctx context.Context, bufConfig *bufconfig.Config) (model.Commits, e.ResponseError) {
	var dependentCommitSet map[string]*model.Commit
	err := resolver.doGetDependenciesFromBufConfig(ctx, dependentCommitSet, bufConfig, true)
	if err != nil {
		return nil, err
	}

	commits := make([]*model.Commit, 0, len(dependentCommitSet))
	for _, commit := range dependentCommitSet {
		commits = append(commits, commit)
	}

	return commits, nil
}

func (resolver *ResolverImpl) GetDirectDependenciesFromBufConfig(ctx context.Context, bufConfig *bufconfig.Config) (model.Commits, e.ResponseError) {
	var dependentCommitSet map[string]*model.Commit
	err := resolver.doGetDependenciesFromBufConfig(ctx, dependentCommitSet, bufConfig, false)
	if err != nil {
		return nil, err
	}

	commits := make([]*model.Commit, 0, len(dependentCommitSet))
	for _, commit := range dependentCommitSet {
		commits = append(commits, commit)
	}

	return commits, nil
}

func (resolver *ResolverImpl) GetBufConfigFromCommitID(ctx context.Context, commitID string) (*bufconfig.Config, e.ResponseError) {
	// 查询manifest名称
	manifestModel, err := resolver.fileMapper.FindManifestByCommitID(commitID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError(fmt.Sprintf("manifest(commit id = %s)", commitID))
		}
		return nil, e.NewInternalError("GetDependenciesByCommitID")
	}

	// 读取
	reader, err := resolver.storageHelper.Read(manifestModel.Digest)
	if err != nil {
		return nil, nil
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, e.NewInternalError("GetDependenciesByCommitID")
	}
	bufConfig, err := bufconfig.GetConfigForData(ctx, data)
	if err != nil {
		return nil, e.NewInternalError("GetDependenciesByCommitID")
	}

	return bufConfig, nil
}

func (resolver *ResolverImpl) GetBufConfigFromBlob(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) (*bufconfig.Config, e.ResponseError) {
	var configFileExist bool
	var configFileData []byte
	externalPaths := []string{
		buflock.ExternalConfigFilePath,
		bufmodule.LicenseFilePath,
	}
	externalPaths = append(externalPaths, bufconfig.AllConfigFilePaths...)
	err := fileManifest.Range(func(path string, digest manifest.Digest) error {
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
				break
			}
		}

		return nil
	})
	if err != nil {
		return nil, e.NewInvalidArgumentError(err.Error())
	}
	if !configFileExist {
		// 不存在配置文件
		return nil, e.NewInvalidArgumentError("no config file")
	}

	// 生成Config，并验证其中的依赖关系
	bufConfig, err := bufconfig.GetConfigForData(ctx, configFileData)
	if err != nil {
		// 无法解析配置文件
		return nil, e.NewInvalidArgumentError(err.Error())
	}

	return bufConfig, nil
}

func (resolver *ResolverImpl) GetBufConfigFromProto(ctx context.Context, protoManifest *modulev1alpha1.Blob, protoBlobs []*modulev1alpha1.Blob) (*bufconfig.Config, e.ResponseError) {
	fileManifest, err := bufmanifest.NewManifestFromProto(ctx, protoManifest)
	if err != nil {
		return nil, e.NewInvalidArgumentError(err.Error())
	}

	blobSet, err := bufmanifest.NewBlobSetFromProto(ctx, protoBlobs)
	if err != nil {
		return nil, e.NewInvalidArgumentError(err.Error())
	}

	return resolver.GetBufConfigFromBlob(ctx, fileManifest, blobSet)
}

func (resolver *ResolverImpl) doGetDependenciesFromBufConfig(ctx context.Context, dependentCommitSet map[string]*model.Commit, bufConfig *bufconfig.Config, getAll bool) e.ResponseError {
	if dependentCommitSet == nil {
		dependentCommitSet = map[string]*model.Commit{}
	}
	dependencyReferences := bufConfig.Build.DependencyModuleReferences

	for i := 0; i < len(dependencyReferences); i++ {
		dependencyReference := dependencyReferences[i]
		if dependencyReference.Remote() == config.Properties.BufMan.ServerHost {
			// 查询repo
			repo, err := resolver.repositoryMapper.FindByUserNameAndRepositoryName(dependencyReference.Owner(), dependencyReference.Repository())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return e.NewNotFoundError(dependencyReference.IdentityString())
				}
				return e.NewInternalError(fmt.Sprintf("find repository(%s)", err.Error()))
			}

			dependentCommit, ok := dependentCommitSet[dependencyReference.IdentityString()]
			if ok && dependentCommit.CommitName == dependencyReference.Reference() {
				continue
			}

			// 查询当前reference，版本号是否一样
			commit, err := resolver.commitMapper.FindByRepositoryIDAndReference(repo.RepositoryID, dependencyReference.Reference())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return e.NewNotFoundError(fmt.Sprintf("%s:%s", dependencyReference.IdentityString(), dependencyReference.Reference()))
				}
				return e.NewInternalError(fmt.Sprintf("find reference(%s)", err.Error()))
			}

			if commit.CommitName != dependentCommit.CommitName {
				// 同一个仓库下的依赖版本号不同，返回错误
				return e.NewInternalError(fmt.Sprintf("two different version %s and %s for %s", dependentCommit.CommitName, dependencyReference.Reference(), dependencyReference.IdentityString()))
			}
			if !ok {
				// 如果之前没有记录过，记录依赖commit
				dependentCommitSet[dependencyReference.IdentityString()] = commit

				if getAll {
					// 需要获取全部依赖，记录这个依赖下的依赖关系
					dependentBufConfig, configErr := resolver.GetBufConfigFromCommitID(ctx, commit.CommitID)
					if err != nil {
						return configErr
					}
					dependentErr := resolver.doGetDependenciesFromBufConfig(ctx, dependentCommitSet, dependentBufConfig, true)
					if dependentErr != nil {
						return dependentErr
					}
				}

			} // else 当前依赖已经记录，跳过
		}
	}

	// 通过
	return nil
}

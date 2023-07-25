package services

import (
	"context"
	"errors"
	"fmt"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util/parse"
	"github.com/ProtobufMan/bufman/internal/util/resolve"
	"github.com/ProtobufMan/bufman/internal/util/storage"
	"gorm.io/gorm"
	"io"
)

type DocsService interface {
	GetSourceDirectoryInfo(repositoryID, reference string) (model.FileBlobs, e.ResponseError)
	GetSourceFile(repositoryID, reference, path string) ([]byte, e.ResponseError)
	GetModulePackages(repositoryID, reference string) ([]*registryv1alpha1.ModulePackage, e.ResponseError)
	GetModuleDocumentation(repositoryID, reference string) (*registryv1alpha1.ModuleDocumentation, e.ResponseError)
	GetPackageDocumentation(repositoryID, reference, packageName string) (*registryv1alpha1.PackageDocumentation, e.ResponseError)
}

type DocsServiceImpl struct {
	commitMapper  mapper.CommitMapper
	fileMapper    mapper.FileMapper
	storageHelper storage.StorageHelper
	protoParser   parse.ProtoParser
	resolver      resolve.Resolver
}

func (docsService *DocsServiceImpl) GetSourceDirectoryInfo(repositoryID, reference string) (model.FileBlobs, e.ResponseError) {
	// 根据reference查询commit
	commit, err := docsService.commitMapper.FindByRepositoryIDAndReference(repositoryID, reference)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("commit")
		}

		return nil, e.NewInternalError(err.Error())
	}

	// 查询所有文件
	fileBlobs, err := docsService.fileMapper.FindAllBlobsByCommitID(commit.CommitID)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return fileBlobs, nil
}

func (docsService *DocsServiceImpl) GetSourceFile(repositoryID, reference, path string) ([]byte, e.ResponseError) {
	// 根据reference查询commit
	commit, err := docsService.commitMapper.FindByRepositoryIDAndReference(repositoryID, reference)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("commit")
		}

		return nil, e.NewInternalError(err.Error())
	}

	// 查询file
	fileBlob, err := docsService.fileMapper.FindBlobByCommitIDAndPath(commit.CommitID, path)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("file path")
		}

		return nil, e.NewInternalError(err.Error())
	}

	// 读取文件
	reader, err := docsService.storageHelper.Read(fileBlob.Digest)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return content, nil
}

func (docsService *DocsServiceImpl) GetModulePackages(repositoryID, reference string) ([]*registryv1alpha1.ModulePackage, e.ResponseError) {
	// 读取commit文件
	fileManifest, blobSet, err := docsService.getManifestAndBlobSet(repositoryID, reference)
	if err != nil {
		return nil, err
	}

	// 读取依赖
	dependentManifests, dependentBlobSets, err := docsService.getDependentManifestsAndBlobSets(fileManifest, blobSet)
	if err != nil {
		return nil, err
	}

	// 获取所有的packages
	packages, err := docsService.protoParser.GetPackages(context.Background(), fileManifest, blobSet, dependentManifests, dependentBlobSets)
	if err != nil {
		return nil, err
	}

	return packages, nil
}

func (docsService *DocsServiceImpl) GetModuleDocumentation(repositoryID, reference string) (*registryv1alpha1.ModuleDocumentation, e.ResponseError) {
	// 读取commit文件
	fileManifest, blobSet, err := docsService.getManifestAndBlobSet(repositoryID, reference)
	if err != nil {
		return nil, err
	}

	documentBlob, licenseBlob, readErr := docsService.storageHelper.GetDocumentAndLicenseFromBlob(context.Background(), fileManifest, blobSet)
	if readErr != nil {
		return nil, e.NewInternalError(readErr.Error())
	}

	// 读取document
	documentReader, readErr := documentBlob.Open(context.Background())
	if readErr != nil {
		return nil, e.NewInternalError(readErr.Error())
	}
	documentData, readErr := io.ReadAll(documentReader)
	if readErr != nil {
		return nil, e.NewInternalError(readErr.Error())
	}

	// 读取license
	licenseReader, readErr := licenseBlob.Open(context.Background())
	if readErr != nil {
		return nil, e.NewInternalError(readErr.Error())
	}
	licenseData, readErr := io.ReadAll(licenseReader)
	if readErr != nil {
		return nil, e.NewInternalError(readErr.Error())
	}

	// 获取documentation path
	paths, _ := fileManifest.PathsFor(documentBlob.Digest().String())
	documentPath := paths[0]

	return &registryv1alpha1.ModuleDocumentation{
		Documentation:     string(documentData),
		License:           string(licenseData),
		DocumentationPath: documentPath,
	}, nil
}

func (docsService *DocsServiceImpl) GetPackageDocumentation(repositoryID, reference, packageName string) (*registryv1alpha1.PackageDocumentation, e.ResponseError) {
	//TODO implement me
	panic("implement me")
}

// getDependentManifestsAndBlobSets 获取依赖的manifests和blob sets
func (docsService *DocsServiceImpl) getDependentManifestsAndBlobSets(fileManifest *manifest.Manifest, blobSet *manifest.BlobSet) ([]*manifest.Manifest, []*manifest.BlobSet, e.ResponseError) {
	// 获取bufConfig
	bufConfig, configErr := docsService.storageHelper.GetBufConfigFromBlob(context.Background(), fileManifest, blobSet)
	if configErr != nil {
		return nil, nil, e.NewInternalError(configErr.Error())
	}

	// 获取全部依赖commits
	dependentCommits, dependenceErr := docsService.resolver.GetAllDependenciesFromBufConfig(context.Background(), bufConfig)
	if dependenceErr != nil {
		return nil, nil, e.NewInternalError(dependenceErr.Error())
	}

	// 读取依赖文件
	dependentManifests := make([]*manifest.Manifest, 0, len(dependentCommits))
	dependentBlobSets := make([]*manifest.BlobSet, 0, len(dependentCommits))
	for i := 0; i < len(dependentCommits); i++ {
		dependentCommit := dependentCommits[i]
		dependentManifest, dependentBlobSet, getErr := docsService.getManifestAndBlobSet(dependentCommit.RepositoryID, dependentCommit.CommitName)
		if getErr != nil {
			return nil, nil, getErr
		}

		dependentManifests = append(dependentManifests, dependentManifest)
		dependentBlobSets = append(dependentBlobSets, dependentBlobSet)
	}

	return dependentManifests, dependentBlobSets, nil
}

func (docsService *DocsServiceImpl) getManifestAndBlobSet(repositoryID string, reference string) (*manifest.Manifest, *manifest.BlobSet, e.ResponseError) {
	// 查询reference对应的commit
	commit, err := docsService.commitMapper.FindByRepositoryIDAndReference(repositoryID, reference)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, e.NewNotFoundError(fmt.Sprintf("repository %s", repositoryID))
		}

		return nil, nil, e.NewInternalError(err.Error())
	}

	// 查询文件清单
	modelFileManifest, err := docsService.fileMapper.FindManifestByCommitID(commit.CommitID)
	if err != nil {
		if err != nil {
			return nil, nil, e.NewInternalError(err.Error())
		}
	}

	// 接着查询blobs
	fileBlobs, err := docsService.fileMapper.FindAllBlobsByCommitID(commit.CommitID)
	if err != nil {
		return nil, nil, e.NewInternalError(err.Error())
	}

	// 读取
	fileManifest, blobSet, err := docsService.storageHelper.ReadToManifestAndBlobSet(modelFileManifest, fileBlobs)
	if err != nil {
		return nil, nil, e.NewInternalError(err.Error())
	}

	return fileManifest, blobSet, nil
}

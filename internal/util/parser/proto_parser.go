package parser

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule/bufmoduleprotocompile"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule/bufmoduleref"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman-cli/private/pkg/thread"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/bufbuild/protocompile"
	"github.com/bufbuild/protocompile/linker"
	"strings"
)

type ProtoParser interface {
	// TryCompile 尝试编译，查看是否能够编译成功
	TryCompile(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) e.ResponseError
	// GetPackageDocumentation 获取package document
	GetPackageDocumentation(ctx context.Context, packageName string, moduleIdentity bufmoduleref.ModuleIdentity, commitName string, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentIdentities []bufmoduleref.ModuleIdentity, dependentCommits []string, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) (*registryv1alpha1.PackageDocumentation, e.ResponseError)
	// GetPackages 获取所有的package
	GetPackages(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) ([]*registryv1alpha1.ModulePackage, e.ResponseError)
}

func NewProtoParser() ProtoParser {
	return &ProtoParserImpl{}
}

type ProtoParserImpl struct {
}

func (protoParser *ProtoParserImpl) GetPackages(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) ([]*registryv1alpha1.ModulePackage, e.ResponseError) {
	module, dependentModules, err := protoParser.getModules(ctx, fileManifest, blobSet, dependentManifests, dependentBlobSets)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	// 编译proto文件
	linkers, _, err := protoParser.compile(ctx, fileManifest, module, dependentModules)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	packagesSet := map[string]struct{}{}
	for _, link := range linkers {
		packagesSet[string(link.Package())] = struct{}{}
	}

	modulePackages := make([]*registryv1alpha1.ModulePackage, 0, len(packagesSet))
	for packageName := range packagesSet {
		modulePackages = append(modulePackages, &registryv1alpha1.ModulePackage{
			Name: packageName,
		})
	}

	return modulePackages, nil
}

func (protoParser *ProtoParserImpl) GetPackageDocumentation(ctx context.Context, packageName string, moduleIdentity bufmoduleref.ModuleIdentity, commitName string, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentIdentities []bufmoduleref.ModuleIdentity, dependentCommits []string, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) (*registryv1alpha1.PackageDocumentation, e.ResponseError) {
	module, dependentModules, err := protoParser.getModulesWithModuleIdentityAndCommit(ctx, moduleIdentity, commitName, fileManifest, blobSet, dependentIdentities, dependentCommits, dependentManifests, dependentBlobSets)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	// 编译proto文件
	linkers, parserAccessorHandler, err := protoParser.compile(ctx, fileManifest, module, dependentModules)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	// 生成package文档
	documentGenerator := NewDocumentGenerator(commitName, linkers, parserAccessorHandler)
	return documentGenerator.GenerateDocument(packageName), nil
}

func (protoParser *ProtoParserImpl) TryCompile(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) e.ResponseError {
	module, dependentModules, err := protoParser.getModules(ctx, fileManifest, blobSet, dependentManifests, dependentBlobSets)
	if err != nil {
		return e.NewInternalError(err.Error())
	}

	// 尝试编译，查看是否成功
	_, _, err = protoParser.compile(ctx, fileManifest, module, dependentModules)
	if err != nil {
		return e.NewInternalError(err.Error())
	}

	return nil
}

func (protoParser *ProtoParserImpl) getModules(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) (bufmodule.Module, []bufmodule.Module, error) {
	module, err := bufmodule.NewModuleForManifestAndBlobSet(ctx, fileManifest, blobSet)
	if err != nil {
		return nil, nil, err
	}
	dependentModules := make([]bufmodule.Module, 0, len(dependentManifests))
	for i := 0; i < len(dependentManifests); i++ {
		dependentModule, err := bufmodule.NewModuleForManifestAndBlobSet(ctx, dependentManifests[i], dependentBlobSets[i])
		if err != nil {
			return nil, nil, err
		}
		dependentModules = append(dependentModules, dependentModule)
	}

	return module, dependentModules, nil
}

func (protoParser *ProtoParserImpl) getModulesWithModuleIdentityAndCommit(ctx context.Context, moduleIdentity bufmoduleref.ModuleIdentity, commit string, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentIdentities []bufmoduleref.ModuleIdentity, dependentCommits []string, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) (bufmodule.Module, []bufmodule.Module, error) {
	module, err := bufmodule.NewModuleForManifestAndBlobSet(ctx, fileManifest, blobSet, bufmodule.ModuleWithModuleIdentityAndCommit(moduleIdentity, commit))
	if err != nil {
		return nil, nil, err
	}
	dependentModules := make([]bufmodule.Module, 0, len(dependentManifests))
	for i := 0; i < len(dependentManifests); i++ {
		dependentModule, err := bufmodule.NewModuleForManifestAndBlobSet(ctx, dependentManifests[i], dependentBlobSets[i], bufmodule.ModuleWithModuleIdentityAndCommit(dependentIdentities[i], dependentCommits[i]))
		if err != nil {
			return nil, nil, err
		}
		dependentModules = append(dependentModules, dependentModule)
	}

	return module, dependentModules, nil
}

func (protoParser *ProtoParserImpl) compile(ctx context.Context, fileManifest *manifest.Manifest, module bufmodule.Module, dependentModules []bufmodule.Module) (linker.Files, bufmoduleprotocompile.ParserAccessorHandler, error) {
	moduleFileSet := bufmodule.NewModuleFileSet(module, dependentModules)
	parserAccessorHandler := bufmoduleprotocompile.NewParserAccessorHandler(ctx, moduleFileSet)
	compiler := protocompile.Compiler{
		MaxParallelism: thread.Parallelism(),
		SourceInfoMode: protocompile.SourceInfoStandard,
		Resolver:       &protocompile.SourceResolver{Accessor: parserAccessorHandler.Open},
	}

	// fileDescriptors are in the same order as paths per the documentation
	protoPaths := protoParser.getProtoPaths(fileManifest)
	linkers, err := compiler.Compile(ctx, protoPaths...)
	if err != nil {
		return nil, nil, err
	}

	return linkers, parserAccessorHandler, nil

}

func (protoParser *ProtoParserImpl) getProtoPaths(fileManifest *manifest.Manifest) []string {

	var protoPaths []string
	_ = fileManifest.Range(func(path string, digest manifest.Digest) error {
		if strings.HasSuffix(path, ".proto") {
			protoPaths = append(protoPaths, path)
		}

		return nil
	})

	return protoPaths
}

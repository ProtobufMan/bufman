package parse

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule"
	"github.com/ProtobufMan/bufman-cli/private/bufpkg/bufmodule/bufmoduleprotocompile"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman-cli/private/pkg/manifest"
	"github.com/ProtobufMan/bufman-cli/private/pkg/thread"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/bufbuild/protocompile"
	"strings"
)

type ProtoParser interface {
	// TryCompile 尝试编译，查看是否能够编译成功
	TryCompile(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) e.ResponseError
	// GetPackageDocumentation 获取package document
	GetPackageDocumentation(ctx context.Context, packageName string, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) (*registryv1alpha1.PackageDocumentation, e.ResponseError)
	// GetPackages 获取所有的package
	GetPackages(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) ([]*registryv1alpha1.ModulePackage, e.ResponseError)
}

func NewProtoParser() ProtoParser {
	return &ProtoParserImpl{}
}

type ProtoParserImpl struct {
}

func (protoParser *ProtoParserImpl) GetPackages(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) ([]*registryv1alpha1.ModulePackage, e.ResponseError) {
	module, err := bufmodule.NewModuleForManifestAndBlobSet(ctx, fileManifest, blobSet)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}
	dependentModules := make([]bufmodule.Module, 0, len(dependentManifests))
	for i := 0; i < len(dependentManifests); i++ {
		dependentModule, err := bufmodule.NewModuleForManifestAndBlobSet(ctx, dependentManifests[i], dependentBlobSets[i])
		if err != nil {
			return nil, e.NewInternalError(err.Error())
		}
		dependentModules = append(dependentModules, dependentModule)
	}
	moduleFileSet := bufmodule.NewModuleFileSet(module, dependentModules)
	parserAccessorHandler := bufmoduleprotocompile.NewParserAccessorHandler(ctx, moduleFileSet)
	compiler := protocompile.Compiler{
		MaxParallelism: thread.Parallelism(),
		SourceInfoMode: protocompile.SourceInfoStandard,
		Resolver:       &protocompile.SourceResolver{Accessor: parserAccessorHandler.Open},
	}

	// fileDescriptors are in the same order as paths per the documentation
	protoPaths := getProtoPaths(fileManifest)
	links, err := compiler.Compile(ctx, protoPaths...)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	packagesSet := map[string]struct{}{}
	for _, link := range links {
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

func (protoParser *ProtoParserImpl) GetPackageDocumentation(ctx context.Context, packageName string, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) (*registryv1alpha1.PackageDocumentation, e.ResponseError) {
	panic("unimplemented")
}

func (protoParser *ProtoParserImpl) TryCompile(ctx context.Context, fileManifest *manifest.Manifest, blobSet *manifest.BlobSet, dependentManifests []*manifest.Manifest, dependentBlobSets []*manifest.BlobSet) e.ResponseError {
	// 检查编译
	module, err := bufmodule.NewModuleForManifestAndBlobSet(ctx, fileManifest, blobSet)
	if err != nil {
		return e.NewInternalError(err.Error())
	}
	dependentModules := make([]bufmodule.Module, 0, len(dependentManifests))
	for i := 0; i < len(dependentManifests); i++ {
		dependentModule, err := bufmodule.NewModuleForManifestAndBlobSet(ctx, dependentManifests[i], dependentBlobSets[i])
		if err != nil {
			return e.NewInternalError(err.Error())
		}
		dependentModules = append(dependentModules, dependentModule)
	}
	moduleFileSet := bufmodule.NewModuleFileSet(module, dependentModules)
	parserAccessorHandler := bufmoduleprotocompile.NewParserAccessorHandler(ctx, moduleFileSet)
	compiler := protocompile.Compiler{
		MaxParallelism: thread.Parallelism(),
		SourceInfoMode: protocompile.SourceInfoStandard,
		Resolver:       &protocompile.SourceResolver{Accessor: parserAccessorHandler.Open},
	}

	// fileDescriptors are in the same order as paths per the documentation
	protoPaths := getProtoPaths(fileManifest)
	_, err = compiler.Compile(ctx, protoPaths...)
	if err != nil {
		return e.NewInternalError(err.Error())
	}

	return nil
}

func getProtoPaths(fileManifest *manifest.Manifest) []string {

	var protoPaths []string
	_ = fileManifest.Range(func(path string, digest manifest.Digest) error {
		if strings.HasSuffix(path, ".proto") {
			protoPaths = append(protoPaths, path)
		}

		return nil
	})

	return protoPaths
}

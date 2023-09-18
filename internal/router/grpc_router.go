package router

import (
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/ProtobufMan/bufman/internal/handlers/grpc_handlers"
	"github.com/ProtobufMan/bufman/internal/interceptors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitGRPCRouter(router *gin.Engine) {
	// UserService
	userServicePath, userServiceHandler := registryv1alpha1connect.NewUserServiceHandler(grpc_handlers.NewUserServiceHandler())
	registerHandler(router, userServicePath, userServiceHandler)

	// TokenService
	tokenServicePath, tokenServiceHandler := registryv1alpha1connect.NewTokenServiceHandler(grpc_handlers.NewTokenServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alpha1connect.TokenServiceGetTokenProcedure,
			registryv1alpha1connect.TokenServiceListTokensProcedure,
			registryv1alpha1connect.TokenServiceDeleteTokenProcedure,
		),
	)
	registerHandler(router, tokenServicePath, tokenServiceHandler)

	// AuthnService
	authnServicePath, authnServiceHandler := registryv1alpha1connect.NewAuthnServiceHandler(grpc_handlers.NewAuthnServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alpha1connect.AuthnServiceGetCurrentUserProcedure,
		),
	)
	registerHandler(router, authnServicePath, authnServiceHandler)

	// RepositoryService
	repositoryServicePath, repositoryServiceHandler := registryv1alpha1connect.NewRepositoryServiceHandler(grpc_handlers.NewRepositoryServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alpha1connect.RepositoryServiceListRepositoriesUserCanAccessProcedure,
			registryv1alpha1connect.RepositoryServiceCreateRepositoryByFullNameProcedure,
			registryv1alpha1connect.RepositoryServiceDeleteRepositoryProcedure,
			registryv1alpha1connect.RepositoryServiceDeleteRepositoryByFullNameProcedure,
			registryv1alpha1connect.RepositoryServiceDeprecateRepositoryByNameProcedure,
			registryv1alpha1connect.RepositoryServiceUndeprecateRepositoryByNameProcedure,
			registryv1alpha1connect.RepositoryServiceUpdateRepositorySettingsByNameProcedure,
		),
	)
	registerHandler(router, repositoryServicePath, repositoryServiceHandler)

	// PushService
	pushServicePath, pushServiceHandler := registryv1alpha1connect.NewPushServiceHandler(grpc_handlers.NewPushServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alpha1connect.PushServicePushManifestAndBlobsProcedure,
		),
	)
	registerHandler(router, pushServicePath, pushServiceHandler)

	// CommitService
	commitServicePath, commitServiceHandler := registryv1alpha1connect.NewRepositoryCommitServiceHandler(grpc_handlers.NewCommitServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alpha1connect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure,
			registryv1alpha1connect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure,
			registryv1alpha1connect.RepositoryCommitServiceListRepositoryDraftCommitsProcedure,
		),
		interceptors.WithAuthInterceptor(
			registryv1alpha1connect.RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure,
		),
	)
	registerHandler(router, commitServicePath, commitServiceHandler)

	// TagService
	tagServicePath, tagServiceHandler := registryv1alpha1connect.NewRepositoryTagServiceHandler(grpc_handlers.NewTagServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alpha1connect.RepositoryTagServiceListRepositoryTagsProcedure,
		),
		interceptors.WithAuthInterceptor(
			registryv1alpha1connect.RepositoryTagServiceCreateRepositoryTagProcedure,
		),
	)
	registerHandler(router, tagServicePath, tagServiceHandler)

	// ResolveService
	resolveServicePath, resolveServiceHandler := registryv1alpha1connect.NewResolveServiceHandler(grpc_handlers.NewResolveServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alpha1connect.ResolveServiceGetModulePinsProcedure,
		),
	)
	registerHandler(router, resolveServicePath, resolveServiceHandler)

	// DownloadService
	downloadServicePath, downloadServiceHandler := registryv1alpha1connect.NewDownloadServiceHandler(grpc_handlers.NewDownloadServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alpha1connect.DownloadServiceDownloadManifestAndBlobsProcedure,
		),
	)
	registerHandler(router, downloadServicePath, downloadServiceHandler)

	// PluginService
	pluginServicePath, pluginServiceHandler := registryv1alpha1connect.NewPluginCurationServiceHandler(grpc_handlers.NewPluginServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alpha1connect.PluginCurationServiceCreateCuratedPluginProcedure,
		),
	)
	registerHandler(router, pluginServicePath, pluginServiceHandler)

	// CodeGenerateService
	codeGenerateServicePath, codeGenerateServiceHandler := registryv1alpha1connect.NewCodeGenerationServiceHandler(grpc_handlers.NewCodeGenerateServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alpha1connect.CodeGenerationServiceGenerateCodeProcedure,
		),
	)
	registerHandler(router, codeGenerateServicePath, codeGenerateServiceHandler)

	// DocService
	docServicePath, docsServiceHandler := registryv1alpha1connect.NewDocServiceHandler(grpc_handlers.NewDocServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alpha1connect.DocServiceGetSourceDirectoryInfoProcedure,
			registryv1alpha1connect.DocServiceGetSourceFileProcedure,
			registryv1alpha1connect.DocServiceGetModulePackagesProcedure,
			registryv1alpha1connect.DocServiceGetPackageDocumentationProcedure,
			registryv1alpha1connect.DocServiceGetModuleDocumentationProcedure,
		),
	)
	registerHandler(router, docServicePath, docsServiceHandler)

	// DockerRepoService
	dockerRepoServicePath, dockerRepoServiceHandler := registryv1alpha1connect.NewDockerRepoServiceHandler(grpc_handlers.NewDockerRepoServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alpha1connect.DockerRepoServiceCreateDockerRepoProcedure,
			registryv1alpha1connect.DockerRepoServiceGetDockerRepoProcedure,
			registryv1alpha1connect.DockerRepoServiceGetDockerRepoByNameProcedure,
			registryv1alpha1connect.DockerRepoServiceListDockerReposProcedure,
			registryv1alpha1connect.DockerRepoServiceUpdateDockerRepoByIDProcedure,
			registryv1alpha1connect.DockerRepoServiceUpdateDockerRepoByNameProcedure,
		),
	)
	registerHandler(router, dockerRepoServicePath, dockerRepoServiceHandler)

	// SearchService
	searchServicePath, searchServiceHandler := registryv1alpha1connect.NewSearchServiceHandler(grpc_handlers.NewSearchServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alpha1connect.SearchServiceSearchTagProcedure,
			registryv1alpha1connect.SearchServiceSearchDraftProcedure,
		),
	)
	registerHandler(router, searchServicePath, searchServiceHandler)
}

func registerHandler(router *gin.Engine, path string, handler http.Handler) {
	router.Handle(http.MethodPost, path+"/*action", gin.WrapH(handler))
}

package router

import (
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/handlers"
	"github.com/ProtobufMan/bufman/internal/interceptors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// UserService
	userServicePath, userServiceHandler := registryv1alphaconnect.NewUserServiceHandler(handlers.NewUserServiceHandler())
	registerHandler(router, userServicePath, userServiceHandler)

	// TokenService
	tokenServicePath, tokenServiceHandler := registryv1alphaconnect.NewTokenServiceHandler(handlers.NewTokenServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alphaconnect.TokenServiceGetTokenProcedure,
			registryv1alphaconnect.TokenServiceListTokensProcedure,
			registryv1alphaconnect.TokenServiceDeleteTokenProcedure,
		),
	)
	registerHandler(router, tokenServicePath, tokenServiceHandler)

	// AuthnService
	authnServicePath, authnServiceHandler := registryv1alphaconnect.NewAuthnServiceHandler(handlers.NewAuthnServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alphaconnect.AuthnServiceGetCurrentUserProcedure,
		),
	)
	registerHandler(router, authnServicePath, authnServiceHandler)

	// RepositoryService
	repositoryServicePath, repositoryServiceHandler := registryv1alphaconnect.NewRepositoryServiceHandler(handlers.NewRepositoryServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alphaconnect.RepositoryServiceListRepositoriesUserCanAccessProcedure,
			registryv1alphaconnect.RepositoryServiceCreateRepositoryByFullNameProcedure,
			registryv1alphaconnect.RepositoryServiceDeleteRepositoryProcedure,
			registryv1alphaconnect.RepositoryServiceDeleteRepositoryByFullNameProcedure,
			registryv1alphaconnect.RepositoryServiceDeprecateRepositoryByNameProcedure,
			registryv1alphaconnect.RepositoryServiceUndeprecateRepositoryByNameProcedure,
			registryv1alphaconnect.RepositoryServiceUpdateRepositorySettingsByNameProcedure,
		),
	)
	registerHandler(router, repositoryServicePath, repositoryServiceHandler)

	// PushService
	pushServicePath, pushServiceHandler := registryv1alphaconnect.NewPushServiceHandler(handlers.NewPushServiceHandler(),
		interceptors.WithAuthInterceptor(
			registryv1alphaconnect.PushServicePushManifestAndBlobsProcedure,
		),
	)
	registerHandler(router, pushServicePath, pushServiceHandler)

	// CommitService
	commitServicePath, commitServiceHandler := registryv1alphaconnect.NewRepositoryCommitServiceHandler(handlers.NewCommitServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alphaconnect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure,
			registryv1alphaconnect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure,
			registryv1alphaconnect.RepositoryCommitServiceListRepositoryDraftCommitsProcedure,
		),
		interceptors.WithAuthInterceptor(
			registryv1alphaconnect.RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure,
		),
	)
	registerHandler(router, commitServicePath, commitServiceHandler)

	// TagService
	tagServicePath, tagServiceHandler := registryv1alphaconnect.NewRepositoryTagServiceHandler(handlers.NewTagServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alphaconnect.RepositoryTagServiceListRepositoryTagsProcedure,
		),
		interceptors.WithAuthInterceptor(
			registryv1alphaconnect.RepositoryTagServiceCreateRepositoryTagProcedure,
		),
	)
	registerHandler(router, tagServicePath, tagServiceHandler)

	// ResolveService
	resolveServicePath, resolveServiceHandler := registryv1alphaconnect.NewResolveServiceHandler(handlers.NewResolveServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alphaconnect.ResolveServiceGetModulePinsProcedure,
		),
	)
	registerHandler(router, resolveServicePath, resolveServiceHandler)

	// DownloadService
	downloadServicePath, downloadServiceHandler := registryv1alphaconnect.NewDownloadServiceHandler(handlers.NewDownloadServiceHandler(),
		interceptors.WithOptionalAuthInterceptor(
			registryv1alphaconnect.DownloadServiceDownloadManifestAndBlobsProcedure,
		),
	)
	registerHandler(router, downloadServicePath, downloadServiceHandler)

	return router
}

func registerHandler(router *gin.Engine, path string, handler http.Handler) {
	router.Handle(http.MethodPost, path+"/*action", gin.WrapH(handler))
}

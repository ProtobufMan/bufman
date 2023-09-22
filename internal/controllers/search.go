package controllers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
)

type SearchController struct {
	validator            validity.Validator
	searchService        services.SearchService
	authorizationService services.AuthorizationService
}

func NewSearchController() *SearchController {
	return &SearchController{
		validator:            validity.NewValidator(),
		searchService:        services.NewSearchService(),
		authorizationService: services.NewAuthorizationService(),
	}
}

func (controller *SearchController) SearchUser(ctx context.Context, req *registryv1alpha1.SearchUserRequest) (*registryv1alpha1.SearchUserResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}
	argErr = controller.validator.CheckQuery(req.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	// 查询结果
	users, respErr := controller.searchService.SearchUser(ctx, req.GetQuery(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search user: %v\n", respErr.Error())

		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(users))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.SearchUserResponse{
		Users:         users.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	}

	return resp, nil
}

func (controller *SearchController) SearchRepository(ctx context.Context, req *registryv1alpha1.SearchRepositoryRequest) (*registryv1alpha1.SearchRepositoryResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}
	argErr = controller.validator.CheckQuery(req.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	// 查询结果
	repositories, respErr := controller.searchService.SearchRepository(ctx, req.GetQuery(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search repo: %v\n", respErr.Error())

		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(repositories))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.SearchRepositoryResponse{
		Repositories:  repositories.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	}

	return resp, nil
}

func (controller *SearchController) SearchLastCommitByContent(ctx context.Context, req *registryv1alpha1.SearchLastCommitByContentRequest) (*registryv1alpha1.SearchLastCommitByContentResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}
	argErr = controller.validator.CheckQuery(req.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	commits, respErr := controller.searchService.SearchLastCommitByContent(ctx, req.GetQuery(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search commit by content: %v\n", err.Error())

		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(commits))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.SearchLastCommitByContentResponse{
		Commits:       commits.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	}

	return resp, nil
}

func (controller *SearchController) SearchCurationPlugin(ctx context.Context, req *registryv1alpha1.SearchCuratedPluginRequest) (*registryv1alpha1.SearchCuratedPluginResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}
	argErr = controller.validator.CheckQuery(req.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	plugins, respErr := controller.searchService.SearchCurationPlugin(ctx, req.GetQuery(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search curation plugin: %v\n", err.Error())

		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(plugins))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.SearchCuratedPluginResponse{
		Plugins:       plugins.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	}

	return resp, nil
}

func (controller *SearchController) SearchTag(ctx context.Context, req *registryv1alpha1.SearchTagRequest) (*registryv1alpha1.SearchTagResponse, e.ResponseError) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}
	argErr = controller.validator.CheckQuery(req.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 查询权限
	repository, checkErr := controller.authorizationService.CheckRepositoryCanAccess(userID, req.GetRepositoryOwner(), req.GetRepositoryName(), registryv1alpha1connect.SearchServiceSearchTagProcedure)
	if checkErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, checkErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	tags, respErr := controller.searchService.SearchTag(ctx, repository.RepositoryID, req.GetQuery(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search tag: %v\n", err.Error())

		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(tags))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.SearchTagResponse{
		RepositoryTags: tags.ToProtoRepositoryTags(),
		NextPageToken:  nextPageToken,
	}

	return resp, nil
}

func (controller *SearchController) SearchDraft(ctx context.Context, req *registryv1alpha1.SearchDraftRequest) (*registryv1alpha1.SearchDraftResponse, e.ResponseError) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}
	argErr = controller.validator.CheckQuery(req.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 查询权限
	repository, checkErr := controller.authorizationService.CheckRepositoryCanAccess(userID, req.GetRepositoryOwner(), req.GetRepositoryName(), registryv1alpha1connect.SearchServiceSearchTagProcedure)
	if checkErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, checkErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	commits, respErr := controller.searchService.SearchDraft(ctx, repository.RepositoryID, req.GetQuery(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search draft: %v\n", err.Error())

		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(commits))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.SearchDraftResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
		NextPageToken:     nextPageToken,
	}

	return resp, nil
}

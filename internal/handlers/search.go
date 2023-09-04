package handlers

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
	"github.com/bufbuild/connect-go"
)

type SearchServiceHandler struct {
	validator     validity.Validator
	searchService services.SearchService
}

func NewSearchServiceHandler() *SearchServiceHandler {
	return &SearchServiceHandler{
		validator:     validity.NewValidator(),
		searchService: services.NewSearchService(),
	}
}

func (handler *SearchServiceHandler) SearchUser(ctx context.Context, req *connect.Request[registryv1alpha1.SearchUserRequest]) (*connect.Response[registryv1alpha1.SearchUserResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckQuery(req.Msg.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	users, respErr := handler.searchService.SearchUser(ctx, req.Msg.GetQuery(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search user: %v\n", respErr.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(users))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.SearchUserResponse{
		Users:         users.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	})

	return resp, nil
}

func (handler *SearchServiceHandler) SearchRepository(ctx context.Context, req *connect.Request[registryv1alpha1.SearchRepositoryRequest]) (*connect.Response[registryv1alpha1.SearchRepositoryResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckQuery(req.Msg.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	repositories, respErr := handler.searchService.SearchRepository(ctx, req.Msg.GetQuery(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search repo: %v\n", respErr.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(repositories))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.SearchRepositoryResponse{
		Repositories:  repositories.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	})

	return resp, nil
}

func (handler *SearchServiceHandler) SearchLastCommitByContent(ctx context.Context, req *connect.Request[registryv1alpha1.SearchLastCommitByContentRequest]) (*connect.Response[registryv1alpha1.SearchLastCommitByContentResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckQuery(req.Msg.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	commits, respErr := handler.searchService.SearchLastCommitByContent(ctx, req.Msg.GetQuery(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search commit by content: %v\n", err.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(commits))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.SearchLastCommitByContentResponse{
		Commits:       commits.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	})

	return resp, nil
}

func (handler *SearchServiceHandler) SearchCurationPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.SearchCuratedPluginRequest]) (*connect.Response[registryv1alpha1.SearchCuratedPluginResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckQuery(req.Msg.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	plugins, respErr := handler.searchService.SearchCurationPlugin(ctx, req.Msg.GetQuery(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search curation plugin: %v\n", err.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(plugins))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.SearchCuratedPluginResponse{
		Plugins:       plugins.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	})

	return resp, nil
}

func (handler *SearchServiceHandler) SearchTag(ctx context.Context, req *connect.Request[registryv1alpha1.SearchTagRequest]) (*connect.Response[registryv1alpha1.SearchTagResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckQuery(req.Msg.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 查询权限
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alpha1connect.SearchServiceSearchTagProcedure)
	if checkErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	tags, respErr := handler.searchService.SearchTag(ctx, repository.RepositoryID, req.Msg.GetQuery(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search tag: %v\n", err.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(tags))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.SearchTagResponse{
		RepositoryTags: tags.ToProtoRepositoryTags(),
		NextPageToken:  nextPageToken,
	})

	return resp, nil
}

func (handler *SearchServiceHandler) SearchDraft(ctx context.Context, req *connect.Request[registryv1alpha1.SearchDraftRequest]) (*connect.Response[registryv1alpha1.SearchDraftResponse], error) {
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckQuery(req.Msg.GetQuery())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 查询权限
	repository, checkErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alpha1connect.SearchServiceSearchTagProcedure)
	if checkErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	commits, respErr := handler.searchService.SearchDraft(ctx, repository.RepositoryID, req.Msg.GetQuery(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error search draft: %v\n", err.Error())

		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(commits))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.SearchDraftResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
		NextPageToken:     nextPageToken,
	})

	return resp, nil
}

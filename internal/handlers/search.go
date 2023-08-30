package handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/util/security"
	"github.com/ProtobufMan/bufman/internal/util/validity"
	"github.com/bufbuild/connect-go"
)

type SearchServiceHandler struct {
	validator     validity.Validator
	searchService services.SearchService
}

func (handler *SearchServiceHandler) SearchUser(ctx context.Context, req *connect.Request[registryv1alpha1.SearchUserRequest]) (*connect.Response[registryv1alpha1.SearchUserResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckQuery(req.Msg.GetQuery())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询结果
	users, respErr := handler.searchService.SearchUser(ctx, req.Msg.GetQuery(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(users))
	if err != nil {
		return nil, e.NewInternalError("generate next page token")
	}

	resp := connect.NewResponse(&registryv1alpha1.SearchUserResponse{
		Users:         users.ToProtoSearchResults(),
		NextPageToken: nextPageToken,
	})

	return resp, nil
}

func (handler *SearchServiceHandler) SearchRepository(ctx context.Context, req *connect.Request[registryv1alpha1.SearchRepositoryRequest]) (*connect.Response[registryv1alpha1.SearchRepositoryResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *SearchServiceHandler) SearchRepositoryByContent(ctx context.Context, req *connect.Request[registryv1alpha1.SearchRepositoryByContentRequest]) (*connect.Response[registryv1alpha1.SearchRepositoryByContentResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *SearchServiceHandler) SearchCurationPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.SearchCuratedPluginRequest]) (*connect.Response[registryv1alpha1.SearchCuratedPluginResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *SearchServiceHandler) SearchTag(ctx context.Context, req *connect.Request[registryv1alpha1.SearchTagRequest]) (*connect.Response[registryv1alpha1.SearchTagResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *SearchServiceHandler) SearchDraft(ctx context.Context, req *connect.Request[registryv1alpha1.SearchDraftRequest]) (*connect.Response[registryv1alpha1.SearchDraftResponse], error) {
	//TODO implement me
	panic("implement me")
}

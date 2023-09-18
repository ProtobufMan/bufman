package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type SearchServiceHandler struct {
	searchController *controllers.SearchController
}

func NewSearchServiceHandler() *SearchServiceHandler {
	return &SearchServiceHandler{
		searchController: controllers.NewSearchController(),
	}
}

func (handler *SearchServiceHandler) SearchUser(ctx context.Context, req *connect.Request[registryv1alpha1.SearchUserRequest]) (*connect.Response[registryv1alpha1.SearchUserResponse], error) {
	resp, err := handler.searchController.SearchUser(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *SearchServiceHandler) SearchRepository(ctx context.Context, req *connect.Request[registryv1alpha1.SearchRepositoryRequest]) (*connect.Response[registryv1alpha1.SearchRepositoryResponse], error) {
	resp, err := handler.searchController.SearchRepository(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *SearchServiceHandler) SearchLastCommitByContent(ctx context.Context, req *connect.Request[registryv1alpha1.SearchLastCommitByContentRequest]) (*connect.Response[registryv1alpha1.SearchLastCommitByContentResponse], error) {
	resp, err := handler.searchController.SearchLastCommitByContent(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *SearchServiceHandler) SearchCurationPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.SearchCuratedPluginRequest]) (*connect.Response[registryv1alpha1.SearchCuratedPluginResponse], error) {
	resp, err := handler.searchController.SearchCurationPlugin(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *SearchServiceHandler) SearchTag(ctx context.Context, req *connect.Request[registryv1alpha1.SearchTagRequest]) (*connect.Response[registryv1alpha1.SearchTagResponse], error) {
	resp, err := handler.searchController.SearchTag(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *SearchServiceHandler) SearchDraft(ctx context.Context, req *connect.Request[registryv1alpha1.SearchDraftRequest]) (*connect.Response[registryv1alpha1.SearchDraftResponse], error) {
	resp, err := handler.searchController.SearchDraft(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

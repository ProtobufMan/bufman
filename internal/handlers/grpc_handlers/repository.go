package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type RepositoryServiceHandler struct {
	repositoryController *controllers.RepositoryController
}

func NewRepositoryServiceHandler() *RepositoryServiceHandler {
	return &RepositoryServiceHandler{
		repositoryController: controllers.NewRepositoryController(),
	}
}

func (handler *RepositoryServiceHandler) GetRepository(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryRequest]) (*connect.Response[registryv1alpha1.GetRepositoryResponse], error) {
	resp, err := handler.repositoryController.GetRepository(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) GetRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha1.GetRepositoryByFullNameResponse], error) {
	resp, err := handler.repositoryController.GetRepositoryByFullName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) ListRepositories(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoriesRequest]) (*connect.Response[registryv1alpha1.ListRepositoriesResponse], error) {
	resp, err := handler.repositoryController.ListRepositories(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) ListUserRepositories(ctx context.Context, req *connect.Request[registryv1alpha1.ListUserRepositoriesRequest]) (*connect.Response[registryv1alpha1.ListUserRepositoriesResponse], error) {
	resp, err := handler.repositoryController.ListUserRepositories(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) ListRepositoriesUserCanAccess(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoriesUserCanAccessRequest]) (*connect.Response[registryv1alpha1.ListRepositoriesUserCanAccessResponse], error) {
	resp, err := handler.repositoryController.ListRepositoriesUserCanAccess(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) CreateRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha1.CreateRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha1.CreateRepositoryByFullNameResponse], error) {
	resp, err := handler.repositoryController.CreateRepositoryByFullName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) DeleteRepository(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteRepositoryRequest]) (*connect.Response[registryv1alpha1.DeleteRepositoryResponse], error) {
	resp, err := handler.repositoryController.DeleteRepository(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) DeleteRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha1.DeleteRepositoryByFullNameResponse], error) {
	resp, err := handler.repositoryController.DeleteRepositoryByFullName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) DeprecateRepositoryByName(ctx context.Context, req *connect.Request[registryv1alpha1.DeprecateRepositoryByNameRequest]) (*connect.Response[registryv1alpha1.DeprecateRepositoryByNameResponse], error) {
	resp, err := handler.repositoryController.DeprecateRepositoryByName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) UndeprecateRepositoryByName(ctx context.Context, req *connect.Request[registryv1alpha1.UndeprecateRepositoryByNameRequest]) (*connect.Response[registryv1alpha1.UndeprecateRepositoryByNameResponse], error) {
	resp, err := handler.repositoryController.UndeprecateRepositoryByName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) UpdateRepositorySettingsByName(ctx context.Context, req *connect.Request[registryv1alpha1.UpdateRepositorySettingsByNameRequest]) (*connect.Response[registryv1alpha1.UpdateRepositorySettingsByNameResponse], error) {
	resp, err := handler.repositoryController.UpdateRepositorySettingsByName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) GetRepositoriesByFullName(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoriesByFullNameRequest]) (*connect.Response[registryv1alpha1.GetRepositoriesByFullNameResponse], error) {
	resp, err := handler.repositoryController.GetRepositoriesByFullName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *RepositoryServiceHandler) GetRepositorySettings(ctx context.Context, c *connect.Request[registryv1alpha1.GetRepositorySettingsRequest]) (*connect.Response[registryv1alpha1.GetRepositorySettingsResponse], error) {
	//TODO implement me
	panic("implement me")
}

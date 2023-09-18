package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type DockerRepoServiceHandler struct {
	dockerRepoController *controllers.DockerRepoController
}

func NewDockerRepoServiceHandler() *DockerRepoServiceHandler {
	return &DockerRepoServiceHandler{
		dockerRepoController: controllers.NewDockerRepoController(),
	}
}

func (handler *DockerRepoServiceHandler) CreateDockerRepo(ctx context.Context, req *connect.Request[registryv1alpha1.CreateDockerRepoRequest]) (*connect.Response[registryv1alpha1.CreateDockerRepoResponse], error) {
	resp, err := handler.dockerRepoController.CreateDockerRepo(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DockerRepoServiceHandler) GetDockerRepo(ctx context.Context, req *connect.Request[registryv1alpha1.GetDockerRepoRequest]) (*connect.Response[registryv1alpha1.GetDockerRepoResponse], error) {
	resp, err := handler.dockerRepoController.GetDockerRepo(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DockerRepoServiceHandler) GetDockerRepoByName(ctx context.Context, req *connect.Request[registryv1alpha1.GetDockerRepoByNameRequest]) (*connect.Response[registryv1alpha1.GetDockerRepoByNameResponse], error) {
	resp, err := handler.dockerRepoController.GetDockerRepoByName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DockerRepoServiceHandler) ListDockerRepos(ctx context.Context, req *connect.Request[registryv1alpha1.ListDockerReposRequest]) (*connect.Response[registryv1alpha1.ListDockerReposResponse], error) {
	resp, err := handler.dockerRepoController.ListDockerRepos(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DockerRepoServiceHandler) UpdateDockerRepoByName(ctx context.Context, req *connect.Request[registryv1alpha1.UpdateDockerRepoByNameRequest]) (*connect.Response[registryv1alpha1.UpdateDockerRepoByNameResponse], error) {
	resp, err := handler.dockerRepoController.UpdateDockerRepoByName(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DockerRepoServiceHandler) UpdateDockerRepoByID(ctx context.Context, req *connect.Request[registryv1alpha1.UpdateDockerRepoByIDRequest]) (*connect.Response[registryv1alpha1.UpdateDockerRepoByIDResponse], error) {
	resp, err := handler.dockerRepoController.UpdateDockerRepoByID(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

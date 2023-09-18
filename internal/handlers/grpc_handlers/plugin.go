package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type PluginServiceHandler struct {
	pluginController *controllers.PluginController
}

func NewPluginServiceHandler() *PluginServiceHandler {
	return &PluginServiceHandler{
		pluginController: controllers.NewPluginController(),
	}
}

func (handler *PluginServiceHandler) ListCuratedPlugins(ctx context.Context, req *connect.Request[registryv1alpha1.ListCuratedPluginsRequest]) (*connect.Response[registryv1alpha1.ListCuratedPluginsResponse], error) {
	resp, err := handler.pluginController.ListCuratedPlugins(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *PluginServiceHandler) CreateCuratedPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.CreateCuratedPluginRequest]) (*connect.Response[registryv1alpha1.CreateCuratedPluginResponse], error) {
	resp, err := handler.pluginController.CreateCuratedPlugin(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *PluginServiceHandler) GetLatestCuratedPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.GetLatestCuratedPluginRequest]) (*connect.Response[registryv1alpha1.GetLatestCuratedPluginResponse], error) {
	resp, err := handler.pluginController.GetLatestCuratedPlugin(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *PluginServiceHandler) DeleteCuratedPlugin(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteCuratedPluginRequest]) (*connect.Response[registryv1alpha1.DeleteCuratedPluginResponse], error) {
	//TODO implement me
	panic("implement me")
}

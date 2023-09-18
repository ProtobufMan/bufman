package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type DocServiceHandler struct {
	docController *controllers.DocController
}

func NewDocServiceHandler() *DocServiceHandler {
	return &DocServiceHandler{
		docController: controllers.NewDocController(),
	}
}

func (handler *DocServiceHandler) GetSourceDirectoryInfo(ctx context.Context, req *connect.Request[registryv1alpha1.GetSourceDirectoryInfoRequest]) (*connect.Response[registryv1alpha1.GetSourceDirectoryInfoResponse], error) {
	resp, err := handler.docController.GetSourceDirectoryInfo(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DocServiceHandler) GetSourceFile(ctx context.Context, req *connect.Request[registryv1alpha1.GetSourceFileRequest]) (*connect.Response[registryv1alpha1.GetSourceFileResponse], error) {
	resp, err := handler.docController.GetSourceFile(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DocServiceHandler) GetModulePackages(ctx context.Context, req *connect.Request[registryv1alpha1.GetModulePackagesRequest]) (*connect.Response[registryv1alpha1.GetModulePackagesResponse], error) {
	resp, err := handler.docController.GetModulePackages(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DocServiceHandler) GetModuleDocumentation(ctx context.Context, req *connect.Request[registryv1alpha1.GetModuleDocumentationRequest]) (*connect.Response[registryv1alpha1.GetModuleDocumentationResponse], error) {
	resp, err := handler.docController.GetModuleDocumentation(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *DocServiceHandler) GetPackageDocumentation(ctx context.Context, req *connect.Request[registryv1alpha1.GetPackageDocumentationRequest]) (*connect.Response[registryv1alpha1.GetPackageDocumentationResponse], error) {
	resp, err := handler.docController.GetPackageDocumentation(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

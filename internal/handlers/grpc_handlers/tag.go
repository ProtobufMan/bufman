package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type TagServiceHandler struct {
	tagController *controllers.TagController
}

func NewTagServiceHandler() *TagServiceHandler {
	return &TagServiceHandler{
		tagController: controllers.NewTagController(),
	}
}

func (handler *TagServiceHandler) CreateRepositoryTag(ctx context.Context, req *connect.Request[registryv1alpha1.CreateRepositoryTagRequest]) (*connect.Response[registryv1alpha1.CreateRepositoryTagResponse], error) {
	resp, err := handler.tagController.CreateRepositoryTag(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *TagServiceHandler) ListRepositoryTags(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryTagsRequest]) (*connect.Response[registryv1alpha1.ListRepositoryTagsResponse], error) {
	resp, err := handler.tagController.ListRepositoryTags(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *TagServiceHandler) ListRepositoryTagsForReference(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryTagsForReferenceRequest]) (*connect.Response[registryv1alpha1.ListRepositoryTagsForReferenceResponse], error) {
	panic("implement me")
}

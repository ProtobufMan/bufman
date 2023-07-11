package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type TagServiceHandler struct {
	tagService services.TagService
}

func NewTagServiceHandler() *TagServiceHandler {
	return &TagServiceHandler{
		tagService: services.NewTagService(),
	}
}

func (handler *TagServiceHandler) CreateRepositoryTag(ctx context.Context, req *connect.Request[registryv1alpha.CreateRepositoryTagRequest]) (*connect.Response[registryv1alpha.CreateRepositoryTagResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	tag, err := handler.tagService.CreateRepositoryTag(userID, req.Msg.GetRepositoryId(), req.Msg.GetName(), req.Msg.GetCommitName())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.CreateRepositoryTagResponse{
		RepositoryTag: tag.ToProtoRepositoryTag(),
	})
	return resp, nil
}

func (handler *TagServiceHandler) ListRepositoryTags(ctx context.Context, req *connect.Request[registryv1alpha.ListRepositoryTagsRequest]) (*connect.Response[registryv1alpha.ListRepositoryTagsResponse], error) {
	var tags model.Tags
	var respErr e.ResponseError

	// 尝试获取user ID
	userID, ok := ctx.Value(constant.UserIDKey).(string)
	if ok {
		tags, respErr = handler.tagService.ListRepositoryTagsWithUserID(userID, req.Msg.GetRepositoryId(), int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	} else {
		tags, respErr = handler.tagService.ListRepositoryTags(req.Msg.GetRepositoryId(), int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	}
	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListRepositoryTagsResponse{
		RepositoryTags: tags.ToProtoRepositoryTags(),
	})
	return resp, nil
}

package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/validity"
	"github.com/bufbuild/connect-go"
)

type TagServiceHandler struct {
	tagService services.TagService
	validator  validity.Validator
}

func NewTagServiceHandler() *TagServiceHandler {
	return &TagServiceHandler{
		tagService: services.NewTagService(),
		validator:  validity.NewValidator(),
	}
}

func (handler *TagServiceHandler) CreateRepositoryTag(ctx context.Context, req *connect.Request[registryv1alpha.CreateRepositoryTagRequest]) (*connect.Response[registryv1alpha.CreateRepositoryTagResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckTagName(req.Msg.GetName())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 获取用户ID
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEditByID(userID, req.Msg.GetRepositoryId(), registryv1alphaconnect.RepositoryTagServiceCreateRepositoryTagProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	tag, err := handler.tagService.CreateRepositoryTag(req.Msg.GetRepositoryId(), req.Msg.GetName(), req.Msg.GetCommitName())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.CreateRepositoryTagResponse{
		RepositoryTag: tag.ToProtoRepositoryTag(),
	})
	return resp, nil
}

func (handler *TagServiceHandler) ListRepositoryTags(ctx context.Context, req *connect.Request[registryv1alpha.ListRepositoryTagsRequest]) (*connect.Response[registryv1alpha.ListRepositoryTagsResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanAccessByID(userID, req.Msg.GetRepositoryId(), registryv1alphaconnect.RepositoryTagServiceCreateRepositoryTagProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	tags, respErr := handler.tagService.ListRepositoryTags(req.Msg.GetRepositoryId(), int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())

	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListRepositoryTagsResponse{
		RepositoryTags: tags.ToProtoRepositoryTags(),
	})
	return resp, nil
}

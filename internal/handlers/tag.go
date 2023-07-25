package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/util/security"
	"github.com/ProtobufMan/bufman/internal/util/validity"
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

func (handler *TagServiceHandler) CreateRepositoryTag(ctx context.Context, req *connect.Request[registryv1alpha1.CreateRepositoryTagRequest]) (*connect.Response[registryv1alpha1.CreateRepositoryTagResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckTagName(req.Msg.GetName())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 获取用户ID
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEditByID(userID, req.Msg.GetRepositoryId(), registryv1alpha1connect.RepositoryTagServiceCreateRepositoryTagProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	tag, err := handler.tagService.CreateRepositoryTag(ctx, req.Msg.GetRepositoryId(), req.Msg.GetName(), req.Msg.GetCommitName())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.CreateRepositoryTagResponse{
		RepositoryTag: tag.ToProtoRepositoryTag(),
	})
	return resp, nil
}

func (handler *TagServiceHandler) ListRepositoryTags(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryTagsRequest]) (*connect.Response[registryv1alpha1.ListRepositoryTagsResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		return nil, e.NewInvalidArgumentError("page token")
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanAccessByID(userID, req.Msg.GetRepositoryId(), registryv1alpha1connect.RepositoryTagServiceListRepositoryTagsProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	tags, respErr := handler.tagService.ListRepositoryTags(ctx, req.Msg.GetRepositoryId(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(tags))
	if err != nil {
		return nil, e.NewInternalError("generate next page token")
	}

	resp := connect.NewResponse(&registryv1alpha1.ListRepositoryTagsResponse{
		RepositoryTags: tags.ToProtoRepositoryTags(),
		NextPageToken:  nextPageToken,
	})
	return resp, nil
}

func (handler *TagServiceHandler) ListRepositoryTagsForReference(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryTagsForReferenceRequest]) (*connect.Response[registryv1alpha1.ListRepositoryTagsForReferenceResponse], error) {
	panic("implement me")
}

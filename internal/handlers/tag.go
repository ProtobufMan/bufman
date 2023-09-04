package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
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
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 获取用户ID
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEditByID(userID, req.Msg.GetRepositoryId(), registryv1alpha1connect.RepositoryTagServiceCreateRepositoryTagProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	tag, err := handler.tagService.CreateRepositoryTag(ctx, req.Msg.GetRepositoryId(), req.Msg.GetName(), req.Msg.GetCommitName())
	if err != nil {
		logger.Errorf("Error create tag: %v", err.Error())

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
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanAccessByID(userID, req.Msg.GetRepositoryId(), registryv1alpha1connect.RepositoryTagServiceListRepositoryTagsProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	tags, respErr := handler.tagService.ListRepositoryTags(ctx, req.Msg.GetRepositoryId(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error list repo tags: %v", respErr.Error())

		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(tags))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
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

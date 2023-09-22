package controllers

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
)

type TagController struct {
	tagService           services.TagService
	authorizationService services.AuthorizationService
	validator            validity.Validator
}

func NewTagController() *TagController {
	return &TagController{
		tagService:           services.NewTagService(),
		authorizationService: services.NewAuthorizationService(),
		validator:            validity.NewValidator(),
	}
}

func (controller *TagController) CreateRepositoryTag(ctx context.Context, req *registryv1alpha1.CreateRepositoryTagRequest) (*registryv1alpha1.CreateRepositoryTagResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckTagName(req.GetName())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 获取用户ID
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := controller.authorizationService.CheckRepositoryCanEditByID(userID, req.GetRepositoryId(), registryv1alpha1connect.RepositoryTagServiceCreateRepositoryTagProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, permissionErr
	}

	tag, err := controller.tagService.CreateRepositoryTag(ctx, req.GetRepositoryId(), req.GetName(), req.GetCommitName())
	if err != nil {
		logger.Errorf("Error create tag: %v", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.CreateRepositoryTagResponse{
		RepositoryTag: tag.ToProtoRepositoryTag(),
	}
	return resp, nil
}

func (controller *TagController) ListRepositoryTags(ctx context.Context, req *registryv1alpha1.ListRepositoryTagsRequest) (*registryv1alpha1.ListRepositoryTagsResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := controller.authorizationService.CheckRepositoryCanAccessByID(userID, req.GetRepositoryId(), registryv1alpha1connect.RepositoryTagServiceListRepositoryTagsProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, permissionErr
	}

	tags, respErr := controller.tagService.ListRepositoryTags(ctx, req.GetRepositoryId(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error list repo tags: %v", respErr.Error())

		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(tags))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListRepositoryTagsResponse{
		RepositoryTags: tags.ToProtoRepositoryTags(),
		NextPageToken:  nextPageToken,
	}
	return resp, nil
}

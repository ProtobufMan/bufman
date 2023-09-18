package controllers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
)

type TokenController struct {
	tokenService services.TokenService
	validator    validity.Validator
}

func NewTokenController() *TokenController {
	return &TokenController{
		tokenService: services.NewTokenService(),
		validator:    validity.NewValidator(),
	}
}

func (controller *TokenController) CreateToken(ctx context.Context, req *registryv1alpha1.CreateTokenRequest) (*registryv1alpha1.CreateTokenResponse, e.ResponseError) {
	token, err := controller.tokenService.CreateToken(ctx, req.GetUsername(), req.GetPassword(), req.GetExpireTime().AsTime(), req.GetNote())
	if err != nil {
		logger.Errorf("Error create token: %v\n", err.Error())

		return nil, err
	}

	// success
	resp := &registryv1alpha1.CreateTokenResponse{
		Token: token.TokenName,
	}
	return resp, nil
}

func (controller *TokenController) GetToken(ctx context.Context, req *registryv1alpha1.GetTokenRequest) (*registryv1alpha1.GetTokenResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询token
	token, err := controller.tokenService.GetToken(ctx, userID, req.GetTokenId())
	if err != nil {
		logger.Errorf("Error get token: %v\n", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.GetTokenResponse{
		Token: token.ToProtoToken(),
	}
	return resp, nil
}

func (controller *TokenController) ListTokens(ctx context.Context, req *registryv1alpha1.ListTokensRequest) (*registryv1alpha1.ListTokensResponse, e.ResponseError) {
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

	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询token
	tokens, listErr := controller.tokenService.ListTokens(ctx, userID, pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if err != nil {
		logger.Errorf("Error list tokens: %v\n", listErr.Error())

		return nil, listErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(tokens))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListTokensResponse{
		Tokens:        tokens.ToProtoTokens(),
		NextPageToken: nextPageToken,
	}
	return resp, nil
}

func (controller *TokenController) DeleteToken(ctx context.Context, req *registryv1alpha1.DeleteTokenRequest) (*registryv1alpha1.DeleteTokenResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 删除token
	err := controller.tokenService.DeleteToken(ctx, userID, req.GetTokenId())
	if err != nil {
		logger.Errorf("Error delete token: %v\n", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.DeleteTokenResponse{}
	return resp, nil
}

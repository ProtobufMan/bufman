package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type TokenServiceHandler struct {
	tokenService services.TokenService
	validator    validity.Validator
}

func NewTokenServiceHandler() *TokenServiceHandler {
	return &TokenServiceHandler{
		tokenService: services.NewTokenService(),
		validator:    validity.NewValidator(),
	}
}

func (handler *TokenServiceHandler) CreateToken(ctx context.Context, req *connect.Request[registryv1alpha1.CreateTokenRequest]) (*connect.Response[registryv1alpha1.CreateTokenResponse], error) {
	token, err := handler.tokenService.CreateToken(ctx, req.Msg.GetUsername(), req.Msg.GetPassword(), req.Msg.GetExpireTime().AsTime(), req.Msg.GetNote())
	if err != nil {
		logger.Errorf("Error create token: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err.Err())
	}

	// success
	resp := connect.NewResponse(&registryv1alpha1.CreateTokenResponse{
		Token: token.TokenName,
	})
	return resp, nil
}

func (handler *TokenServiceHandler) GetToken(ctx context.Context, req *connect.Request[registryv1alpha1.GetTokenRequest]) (*connect.Response[registryv1alpha1.GetTokenResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询token
	token, err := handler.tokenService.GetToken(ctx, userID, req.Msg.GetTokenId())
	if err != nil {
		logger.Errorf("Error get token: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.GetTokenResponse{
		Token: token.ToProtoToken(),
	})
	return resp, nil
}

func (handler *TokenServiceHandler) ListTokens(ctx context.Context, req *connect.Request[registryv1alpha1.ListTokensRequest]) (*connect.Response[registryv1alpha1.ListTokensResponse], error) {
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

	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询token
	tokens, listErr := handler.tokenService.ListTokens(ctx, userID, pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if err != nil {
		logger.Errorf("Error list tokens: %v\n", listErr.Error())

		return nil, connect.NewError(listErr.Code(), listErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(tokens))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.ListTokensResponse{
		Tokens:        tokens.ToProtoTokens(),
		NextPageToken: nextPageToken,
	})
	return resp, nil
}

func (handler *TokenServiceHandler) DeleteToken(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteTokenRequest]) (*connect.Response[registryv1alpha1.DeleteTokenResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 删除token
	err := handler.tokenService.DeleteToken(ctx, userID, req.Msg.GetTokenId())
	if err != nil {
		logger.Errorf("Error delete token: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.DeleteTokenResponse{})
	return resp, nil
}

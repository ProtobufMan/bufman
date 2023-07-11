package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type TokenServiceHandler struct {
	tokenService services.TokenService
}

func NewTokenServiceHandler() *TokenServiceHandler {
	return &TokenServiceHandler{
		tokenService: services.NewTokenService(),
	}
}

func (handler *TokenServiceHandler) CreateToken(ctx context.Context, req *connect.Request[registryv1alpha.CreateTokenRequest]) (*connect.Response[registryv1alpha.CreateTokenResponse], error) {
	token, err := handler.tokenService.CreateToken(req.Msg.GetUsername(), req.Msg.GetPassword(), req.Msg.GetExpireTime().AsTime(), req.Msg.GetNote())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	// success
	resp := connect.NewResponse(&registryv1alpha.CreateTokenResponse{
		Token: token.TokenName,
	})
	return resp, nil
}

func (handler *TokenServiceHandler) GetToken(ctx context.Context, req *connect.Request[registryv1alpha.GetTokenRequest]) (*connect.Response[registryv1alpha.GetTokenResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询token
	token, err := handler.tokenService.GetToken(userID, req.Msg.GetTokenId())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.GetTokenResponse{
		Token: token.ToProtoToken(),
	})
	return resp, nil
}

func (handler *TokenServiceHandler) ListTokens(ctx context.Context, req *connect.Request[registryv1alpha.ListTokensRequest]) (*connect.Response[registryv1alpha.ListTokensResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询token
	tokens, err := handler.tokenService.ListTokens(userID, int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListTokensResponse{
		Tokens: tokens.ToProtoTokens(),
	})
	return resp, nil
}

func (handler *TokenServiceHandler) DeleteToken(ctx context.Context, req *connect.Request[registryv1alpha.DeleteTokenRequest]) (*connect.Response[registryv1alpha.DeleteTokenResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 删除token
	err := handler.tokenService.DeleteToken(userID, req.Msg.GetTokenId())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.DeleteTokenResponse{})
	return resp, nil
}

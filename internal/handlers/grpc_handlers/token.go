package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type TokenServiceHandler struct {
	tokenController *controllers.TokenController
}

func NewTokenServiceHandler() *TokenServiceHandler {
	return &TokenServiceHandler{
		tokenController: controllers.NewTokenController(),
	}
}

func (handler *TokenServiceHandler) CreateToken(ctx context.Context, req *connect.Request[registryv1alpha1.CreateTokenRequest]) (*connect.Response[registryv1alpha1.CreateTokenResponse], error) {
	resp, err := handler.tokenController.CreateToken(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *TokenServiceHandler) GetToken(ctx context.Context, req *connect.Request[registryv1alpha1.GetTokenRequest]) (*connect.Response[registryv1alpha1.GetTokenResponse], error) {
	resp, err := handler.tokenController.GetToken(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *TokenServiceHandler) ListTokens(ctx context.Context, req *connect.Request[registryv1alpha1.ListTokensRequest]) (*connect.Response[registryv1alpha1.ListTokensResponse], error) {
	resp, err := handler.tokenController.ListTokens(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *TokenServiceHandler) DeleteToken(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteTokenRequest]) (*connect.Response[registryv1alpha1.DeleteTokenResponse], error) {
	resp, err := handler.tokenController.DeleteToken(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

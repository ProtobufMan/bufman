package handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/services"
)

type AuthnServiceHandler struct {
	userService     services.UserService
	authnController controllers.AuthnController
}

func NewAuthnServiceHandler() *AuthnServiceHandler {
	return &AuthnServiceHandler{
		userService: services.NewUserService(),
	}
}

func (handler *AuthnServiceHandler) GetCurrentUser(ctx context.Context, req *connect.Request[registryv1alpha1.GetCurrentUserRequest]) (*connect.Response[registryv1alpha1.GetCurrentUserResponse], error) {
	resp, err := handler.authnController.GetCurrentUser(ctx, req.Msg)
	if err != nil {
		logger.Errorf("Error Get User: %v\n", err.Error())
		return nil, connect.NewError(err.Code(), err.Err())
	}

	return connect.NewResponse(resp), nil
}

func (handler *AuthnServiceHandler) GetCurrentUserSubject(ctx context.Context, req *connect.Request[registryv1alpha1.GetCurrentUserSubjectRequest]) (*connect.Response[registryv1alpha1.GetCurrentUserSubjectResponse], error) {
	//TODO implement me
	panic("implement me")
}

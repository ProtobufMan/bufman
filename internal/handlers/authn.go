package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type AuthnServiceHandler struct {
	userService services.UserService
}

func NewAuthnServiceHandler() *AuthnServiceHandler {
	return &AuthnServiceHandler{
		userService: services.NewUserService(),
	}
}

func (handler *AuthnServiceHandler) GetCurrentUser(ctx context.Context, req *connect.Request[registryv1alpha.GetCurrentUserRequest]) (*connect.Response[registryv1alpha.GetCurrentUserResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	user, err := handler.userService.GetUser(userID)
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.GetCurrentUserResponse{
		User: user.ToProtoUser(),
	})
	return resp, nil
}

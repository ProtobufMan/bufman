package handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
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

func (handler *AuthnServiceHandler) GetCurrentUser(ctx context.Context, req *connect.Request[registryv1alpha1.GetCurrentUserRequest]) (*connect.Response[registryv1alpha1.GetCurrentUserResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	user, err := handler.userService.GetUser(ctx, userID)
	if err != nil {
		logger.Errorf("Error Get User: %v\n", err.Error())
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.GetCurrentUserResponse{
		User: user.ToProtoUser(),
	})
	return resp, nil
}

func (handler *AuthnServiceHandler) GetCurrentUserSubject(ctx context.Context, req *connect.Request[registryv1alpha1.GetCurrentUserSubjectRequest]) (*connect.Response[registryv1alpha1.GetCurrentUserSubjectResponse], error) {
	//TODO implement me
	panic("implement me")
}

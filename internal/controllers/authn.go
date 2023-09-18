package controllers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
)

type AuthnController struct {
	userService services.UserService
}

func NewAuthnController() *AuthnController {
	return &AuthnController{
		userService: services.NewUserService(),
	}
}

func (controller *AuthnController) GetCurrentUser(ctx context.Context, req *registryv1alpha1.GetCurrentUserRequest) (*registryv1alpha1.GetCurrentUserResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 	get user by id
	user, err := controller.userService.GetUser(ctx, userID)
	if err != nil {
		logger.Errorf("Error Get User: %v\n", err.Error())
		return nil, err
	}

	resp := &registryv1alpha1.GetCurrentUserResponse{
		User: user.ToProtoUser(),
	}
	return resp, nil
}

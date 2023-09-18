/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type UserServiceHandler struct {
	userService services.UserService
	validator   validity.Validator
}

func NewUserServiceHandler() *UserServiceHandler {
	return &UserServiceHandler{
		userService: services.NewUserService(),
		validator:   validity.NewValidator(),
	}
}

func (handler *UserServiceHandler) CreateUser(ctx context.Context, req *connect.Request[registryv1alpha1.CreateUserRequest]) (*connect.Response[registryv1alpha1.CreateUserResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckUserName(req.Msg.GetUsername())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckPassword(req.Msg.GetPassword())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	user, err := handler.userService.CreateUser(ctx, req.Msg.GetUsername(), req.Msg.GetPassword()) // 创建用户
	if err != nil {
		logger.Errorf("Error create user: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err.Err())
	}

	// success
	resp := connect.NewResponse(&registryv1alpha1.CreateUserResponse{
		User: user.ToProtoUser(),
	})
	return resp, nil
}

func (handler *UserServiceHandler) GetUser(ctx context.Context, req *connect.Request[registryv1alpha1.GetUserRequest]) (*connect.Response[registryv1alpha1.GetUserResponse], error) {
	user, err := handler.userService.GetUser(ctx, req.Msg.GetId()) // 创建用户
	if err != nil {
		logger.Errorf("Error get user: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.GetUserResponse{
		User: user.ToProtoUser(),
	})
	return resp, nil
}
func (handler *UserServiceHandler) GetUserByUsername(ctx context.Context, req *connect.Request[registryv1alpha1.GetUserByUsernameRequest]) (*connect.Response[registryv1alpha1.GetUserByUsernameResponse], error) {
	user, err := handler.userService.GetUserByUsername(ctx, req.Msg.GetUsername()) // 创建用户
	if err != nil {
		logger.Errorf("Error get user: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.GetUserByUsernameResponse{
		User: user.ToProtoUser(),
	})
	return resp, nil
}

func (handler *UserServiceHandler) ListUsers(ctx context.Context, req *connect.Request[registryv1alpha1.ListUsersRequest]) (*connect.Response[registryv1alpha1.ListUsersResponse], error) {
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

	users, ListErr := handler.userService.ListUsers(ctx, pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse()) // 创建用户
	if err != nil {
		logger.Errorf("Error list users: %v\n", ListErr.Error())

		return nil, connect.NewError(ListErr.Code(), ListErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(users))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.ListUsersResponse{
		Users:         users.ToProtoUsers(),
		NextPageToken: nextPageToken,
	})

	return resp, nil
}

func (handler *UserServiceHandler) ListOrganizationUsers(ctx context.Context, req *connect.Request[registryv1alpha1.ListOrganizationUsersRequest]) (*connect.Response[registryv1alpha1.ListOrganizationUsersResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *UserServiceHandler) DeleteUser(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteUserRequest]) (*connect.Response[registryv1alpha1.DeleteUserResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *UserServiceHandler) DeactivateUser(ctx context.Context, req *connect.Request[registryv1alpha1.DeactivateUserRequest]) (*connect.Response[registryv1alpha1.DeactivateUserResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *UserServiceHandler) UpdateUserServerRole(ctx context.Context, req *connect.Request[registryv1alpha1.UpdateUserServerRoleRequest]) (*connect.Response[registryv1alpha1.UpdateUserServerRoleResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *UserServiceHandler) CountUsers(ctx context.Context, req *connect.Request[registryv1alpha1.CountUsersRequest]) (*connect.Response[registryv1alpha1.CountUsersResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *UserServiceHandler) UpdateUserSettings(ctx context.Context, req *connect.Request[registryv1alpha1.UpdateUserSettingsRequest]) (*connect.Response[registryv1alpha1.UpdateUserSettingsResponse], error) {
	//TODO implement me
	panic("implement me")
}

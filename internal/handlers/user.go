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
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type UserServiceHandler struct {
	userService services.UserService
}

func NewUserServiceHandler() *UserServiceHandler {
	return &UserServiceHandler{
		userService: services.NewUserService(),
	}
}

func (handler *UserServiceHandler) CreateUser(ctx context.Context, req *connect.Request[registryv1alpha.CreateUserRequest]) (*connect.Response[registryv1alpha.CreateUserResponse], error) {
	user, err := handler.userService.CreateUser(req.Msg.GetUsername(), req.Msg.GetPassword()) // 创建用户
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	// success
	resp := connect.NewResponse(&registryv1alpha.CreateUserResponse{
		User: user.ToProtoUser(),
	})
	return resp, nil
}

func (handler *UserServiceHandler) GetUser(ctx context.Context, req *connect.Request[registryv1alpha.GetUserRequest]) (*connect.Response[registryv1alpha.GetUserResponse], error) {
	user, err := handler.userService.GetUser(req.Msg.GetId()) // 创建用户
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.GetUserResponse{
		User: user.ToProtoUser(),
	})
	return resp, nil
}
func (handler *UserServiceHandler) GetUserByUsername(ctx context.Context, req *connect.Request[registryv1alpha.GetUserByUsernameRequest]) (*connect.Response[registryv1alpha.GetUserByUsernameResponse], error) {
	user, err := handler.userService.GetUserByUsername(req.Msg.GetUsername()) // 创建用户
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.GetUserByUsernameResponse{
		User: user.ToProtoUser(),
	})
	return resp, nil
}

func (handler *UserServiceHandler) ListUsers(ctx context.Context, req *connect.Request[registryv1alpha.ListUsersRequest]) (*connect.Response[registryv1alpha.ListUsersResponse], error) {
	users, err := handler.userService.ListUsers(int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse()) // 创建用户
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListUsersResponse{
		Users: users.ToProtoUsers(),
	})

	return resp, nil
}

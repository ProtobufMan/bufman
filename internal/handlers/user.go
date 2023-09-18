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
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type UserServiceHandler struct {
	userController *controllers.UserController
}

func NewUserServiceHandler() *UserServiceHandler {
	return &UserServiceHandler{
		userController: controllers.NewUserController(),
	}
}

func (handler *UserServiceHandler) CreateUser(ctx context.Context, req *connect.Request[registryv1alpha1.CreateUserRequest]) (*connect.Response[registryv1alpha1.CreateUserResponse], error) {
	resp, err := handler.userController.CreateUser(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *UserServiceHandler) GetUser(ctx context.Context, req *connect.Request[registryv1alpha1.GetUserRequest]) (*connect.Response[registryv1alpha1.GetUserResponse], error) {
	resp, err := handler.userController.GetUser(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}
func (handler *UserServiceHandler) GetUserByUsername(ctx context.Context, req *connect.Request[registryv1alpha1.GetUserByUsernameRequest]) (*connect.Response[registryv1alpha1.GetUserByUsernameResponse], error) {
	resp, err := handler.userController.GetUserByUsername(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *UserServiceHandler) ListUsers(ctx context.Context, req *connect.Request[registryv1alpha1.ListUsersRequest]) (*connect.Response[registryv1alpha1.ListUsersResponse], error) {
	resp, err := handler.userController.ListUsers(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

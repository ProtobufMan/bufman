package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/validity"
	"github.com/bufbuild/connect-go"
)

type RepositoryServiceHandler struct {
	repositoryService services.RepositoryService
	validator         validity.Validator
}

func NewRepositoryServiceHandler() *RepositoryServiceHandler {
	return &RepositoryServiceHandler{
		repositoryService: services.NewRepositoryService(),
		validator:         validity.NewValidator(),
	}
}

func (handler *RepositoryServiceHandler) GetRepository(ctx context.Context, req *connect.Request[registryv1alpha.GetRepositoryRequest]) (*connect.Response[registryv1alpha.GetRepositoryResponse], error) {
	// 查询
	repository, err := handler.repositoryService.GetRepository(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	repositoryCounts, err := handler.repositoryService.GetRepositoryCounts(repository.RepositoryID)
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	// 查询成功
	resp := connect.NewResponse(&registryv1alpha.GetRepositoryResponse{
		Repository: repository.ToProtoRepository(),
		Counts:     repositoryCounts.ToProtoRepositoryCounts(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) GetRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha.GetRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha.GetRepositoryByFullNameResponse], error) {
	// 验证参数
	userName, repositoryName, argErr := handler.validator.SplitFullName(req.Msg.GetFullName())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 查询
	repository, err := handler.repositoryService.GetRepositoryByUserNameAndRepositoryName(userName, repositoryName)
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	repositoryCounts, err := handler.repositoryService.GetRepositoryCounts(repository.RepositoryID)
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	// 查询成功
	resp := connect.NewResponse(&registryv1alpha.GetRepositoryByFullNameResponse{
		Repository: repository.ToProtoRepository(),
		Counts:     repositoryCounts.ToProtoRepositoryCounts(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) ListRepositories(ctx context.Context, req *connect.Request[registryv1alpha.ListRepositoriesRequest]) (*connect.Response[registryv1alpha.ListRepositoriesResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	repositories, err := handler.repositoryService.ListRepositories(int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.Reverse)
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListRepositoriesResponse{
		Repositories: repositories.ToProtoRepositories(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) ListUserRepositories(ctx context.Context, req *connect.Request[registryv1alpha.ListUserRepositoriesRequest]) (*connect.Response[registryv1alpha.ListUserRepositoriesResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	repositories, err := handler.repositoryService.ListUserRepositories(req.Msg.GetUserId(), int(req.Msg.PageOffset), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListUserRepositoriesResponse{
		Repositories: repositories.ToProtoRepositories(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) ListRepositoriesUserCanAccess(ctx context.Context, req *connect.Request[registryv1alpha.ListRepositoriesUserCanAccessRequest]) (*connect.Response[registryv1alpha.ListRepositoriesUserCanAccessResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	userID := ctx.Value(constant.UserIDKey).(string)
	repositories, err := handler.repositoryService.ListRepositoriesUserCanAccess(userID, int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListRepositoriesUserCanAccessResponse{
		Repositories: repositories.ToProtoRepositories(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) CreateRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha.CreateRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha.CreateRepositoryByFullNameResponse], error) {
	// 验证参数
	userName, repositoryName, argErr := handler.validator.SplitFullName(req.Msg.GetFullName())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}
	argErr = handler.validator.CheckRepositoryName(repositoryName)
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	userID := ctx.Value(constant.UserIDKey).(string)

	// 创建
	repository, err := handler.repositoryService.CreateRepositoryByUserNameAndRepositoryName(userID, userName, repositoryName, req.Msg.GetVisibility())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	// 成功
	resp := connect.NewResponse(&registryv1alpha.CreateRepositoryByFullNameResponse{
		Repository: repository.ToProtoRepository(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) DeleteRepository(ctx context.Context, req *connect.Request[registryv1alpha.DeleteRepositoryRequest]) (*connect.Response[registryv1alpha.DeleteRepositoryResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanDeleteByID(userID, req.Msg.GetId(), registryv1alphaconnect.RepositoryServiceDeleteRepositoryProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 查询repository，检查是否可以删除
	err := handler.repositoryService.DeleteRepository(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.DeleteRepositoryResponse{})
	return resp, nil
}

func (handler *RepositoryServiceHandler) DeleteRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha.DeleteRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha.DeleteRepositoryByFullNameResponse], error) {
	// 验证参数
	userName, repositoryName, argErr := handler.validator.SplitFullName(req.Msg.GetFullName())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanDelete(userID, userName, repositoryName, registryv1alphaconnect.RepositoryServiceDeleteRepositoryByFullNameProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 删除
	err := handler.repositoryService.DeleteRepositoryByUserNameAndRepositoryName(userName, repositoryName)
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.DeleteRepositoryByFullNameResponse{})
	return resp, nil
}

func (handler *RepositoryServiceHandler) DeprecateRepositoryByName(ctx context.Context, req *connect.Request[registryv1alpha.DeprecateRepositoryByNameRequest]) (*connect.Response[registryv1alpha.DeprecateRepositoryByNameResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEdit(userID, req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), registryv1alphaconnect.RepositoryServiceDeprecateRepositoryByNameProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 修改数据库
	updatedRepository, err := handler.repositoryService.DeprecateRepositoryByName(req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), req.Msg.GetDeprecationMessage())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.DeprecateRepositoryByNameResponse{
		Repository: updatedRepository.ToProtoRepository(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) UndeprecateRepositoryByName(ctx context.Context, req *connect.Request[registryv1alpha.UndeprecateRepositoryByNameRequest]) (*connect.Response[registryv1alpha.UndeprecateRepositoryByNameResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEdit(userID, req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), registryv1alphaconnect.RepositoryServiceUndeprecateRepositoryByNameProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 修改数据库
	updatedRepository, err := handler.repositoryService.UndeprecateRepositoryByName(req.Msg.GetOwnerName(), req.Msg.GetRepositoryName())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.UndeprecateRepositoryByNameResponse{
		Repository: updatedRepository.ToProtoRepository(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) UpdateRepositorySettingsByName(ctx context.Context, req *connect.Request[registryv1alpha.UpdateRepositorySettingsByNameRequest]) (*connect.Response[registryv1alpha.UpdateRepositorySettingsByNameResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEdit(userID, req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), registryv1alphaconnect.RepositoryServiceUpdateRepositorySettingsByNameProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 修改数据库
	err := handler.repositoryService.UpdateRepositorySettingsByName(req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), req.Msg.GetVisibility(), req.Msg.GetDescription())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.UpdateRepositorySettingsByNameResponse{})
	return resp, nil
}

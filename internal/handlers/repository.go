package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/util/security"
	"github.com/ProtobufMan/bufman/internal/util/validity"
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

func (handler *RepositoryServiceHandler) GetRepository(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryRequest]) (*connect.Response[registryv1alpha1.GetRepositoryResponse], error) {
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
	resp := connect.NewResponse(&registryv1alpha1.GetRepositoryResponse{
		Repository: repository.ToProtoRepository(),
		Counts:     repositoryCounts.ToProtoRepositoryCounts(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) GetRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha1.GetRepositoryByFullNameResponse], error) {
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
	resp := connect.NewResponse(&registryv1alpha1.GetRepositoryByFullNameResponse{
		Repository: repository.ToProtoRepository(),
		Counts:     repositoryCounts.ToProtoRepositoryCounts(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) ListRepositories(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoriesRequest]) (*connect.Response[registryv1alpha1.ListRepositoriesResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		return nil, e.NewInvalidArgumentError("page token")
	}

	repositories, listErr := handler.repositoryService.ListRepositories(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.Reverse)
	if listErr != nil {
		return nil, connect.NewError(listErr.Code(), listErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(repositories))
	if err != nil {
		return nil, e.NewInternalError("generate next page token")
	}

	resp := connect.NewResponse(&registryv1alpha1.ListRepositoriesResponse{
		Repositories:  repositories.ToProtoRepositories(),
		NextPageToken: nextPageToken,
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) ListUserRepositories(ctx context.Context, req *connect.Request[registryv1alpha1.ListUserRepositoriesRequest]) (*connect.Response[registryv1alpha1.ListUserRepositoriesResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		return nil, e.NewInvalidArgumentError("page token")
	}

	repositories, listErr := handler.repositoryService.ListUserRepositories(req.Msg.GetUserId(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if err != nil {
		return nil, connect.NewError(listErr.Code(), listErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(repositories))
	if err != nil {
		return nil, e.NewInternalError("generate next page token")
	}

	resp := connect.NewResponse(&registryv1alpha1.ListUserRepositoriesResponse{
		Repositories:  repositories.ToProtoRepositories(),
		NextPageToken: nextPageToken,
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) ListRepositoriesUserCanAccess(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoriesUserCanAccessRequest]) (*connect.Response[registryv1alpha1.ListRepositoriesUserCanAccessResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		return nil, e.NewInvalidArgumentError("page token")
	}

	userID := ctx.Value(constant.UserIDKey).(string)
	repositories, ListErr := handler.repositoryService.ListRepositoriesUserCanAccess(userID, pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if err != nil {
		return nil, connect.NewError(ListErr.Code(), ListErr)
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(repositories))
	if err != nil {
		return nil, e.NewInternalError("generate next page token")
	}

	resp := connect.NewResponse(&registryv1alpha1.ListRepositoriesUserCanAccessResponse{
		Repositories:  repositories.ToProtoRepositories(),
		NextPageToken: nextPageToken,
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) CreateRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha1.CreateRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha1.CreateRepositoryByFullNameResponse], error) {
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
	resp := connect.NewResponse(&registryv1alpha1.CreateRepositoryByFullNameResponse{
		Repository: repository.ToProtoRepository(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) DeleteRepository(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteRepositoryRequest]) (*connect.Response[registryv1alpha1.DeleteRepositoryResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanDeleteByID(userID, req.Msg.GetId(), registryv1alpha1connect.RepositoryServiceDeleteRepositoryProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 查询repository，检查是否可以删除
	err := handler.repositoryService.DeleteRepository(req.Msg.GetId())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.DeleteRepositoryResponse{})
	return resp, nil
}

func (handler *RepositoryServiceHandler) DeleteRepositoryByFullName(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteRepositoryByFullNameRequest]) (*connect.Response[registryv1alpha1.DeleteRepositoryByFullNameResponse], error) {
	// 验证参数
	userName, repositoryName, argErr := handler.validator.SplitFullName(req.Msg.GetFullName())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanDelete(userID, userName, repositoryName, registryv1alpha1connect.RepositoryServiceDeleteRepositoryByFullNameProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 删除
	err := handler.repositoryService.DeleteRepositoryByUserNameAndRepositoryName(userName, repositoryName)
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.DeleteRepositoryByFullNameResponse{})
	return resp, nil
}

func (handler *RepositoryServiceHandler) DeprecateRepositoryByName(ctx context.Context, req *connect.Request[registryv1alpha1.DeprecateRepositoryByNameRequest]) (*connect.Response[registryv1alpha1.DeprecateRepositoryByNameResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEdit(userID, req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryServiceDeprecateRepositoryByNameProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 修改数据库
	updatedRepository, err := handler.repositoryService.DeprecateRepositoryByName(req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), req.Msg.GetDeprecationMessage())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.DeprecateRepositoryByNameResponse{
		Repository: updatedRepository.ToProtoRepository(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) UndeprecateRepositoryByName(ctx context.Context, req *connect.Request[registryv1alpha1.UndeprecateRepositoryByNameRequest]) (*connect.Response[registryv1alpha1.UndeprecateRepositoryByNameResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEdit(userID, req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryServiceUndeprecateRepositoryByNameProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 修改数据库
	updatedRepository, err := handler.repositoryService.UndeprecateRepositoryByName(req.Msg.GetOwnerName(), req.Msg.GetRepositoryName())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.UndeprecateRepositoryByNameResponse{
		Repository: updatedRepository.ToProtoRepository(),
	})
	return resp, nil
}

func (handler *RepositoryServiceHandler) UpdateRepositorySettingsByName(ctx context.Context, req *connect.Request[registryv1alpha1.UpdateRepositorySettingsByNameRequest]) (*connect.Response[registryv1alpha1.UpdateRepositorySettingsByNameResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := handler.validator.CheckRepositoryCanEdit(userID, req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryServiceUpdateRepositorySettingsByNameProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 修改数据库
	err := handler.repositoryService.UpdateRepositorySettingsByName(req.Msg.GetOwnerName(), req.Msg.GetRepositoryName(), req.Msg.GetVisibility(), req.Msg.GetDescription())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.UpdateRepositorySettingsByNameResponse{})
	return resp, nil
}

func (handler *RepositoryServiceHandler) ListOrganizationRepositories(ctx context.Context, req *connect.Request[registryv1alpha1.ListOrganizationRepositoriesRequest]) (*connect.Response[registryv1alpha1.ListOrganizationRepositoriesResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *RepositoryServiceHandler) GetRepositoriesByFullName(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoriesByFullNameRequest]) (*connect.Response[registryv1alpha1.GetRepositoriesByFullNameResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *RepositoryServiceHandler) SetRepositoryContributor(ctx context.Context, req *connect.Request[registryv1alpha1.SetRepositoryContributorRequest]) (*connect.Response[registryv1alpha1.SetRepositoryContributorResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *RepositoryServiceHandler) ListRepositoryContributors(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryContributorsRequest]) (*connect.Response[registryv1alpha1.ListRepositoryContributorsResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *RepositoryServiceHandler) GetRepositoryContributor(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryContributorRequest]) (*connect.Response[registryv1alpha1.GetRepositoryContributorResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *RepositoryServiceHandler) GetRepositorySettings(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositorySettingsRequest]) (*connect.Response[registryv1alpha1.GetRepositorySettingsResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *RepositoryServiceHandler) GetRepositoriesMetadata(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoriesMetadataRequest]) (*connect.Response[registryv1alpha1.GetRepositoriesMetadataResponse], error) {
	//TODO implement me
	panic("implement me")
}

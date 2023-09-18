package controllers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
)

type RepositoryController struct {
	repositoryService services.RepositoryService
	validator         validity.Validator
}

func NewRepositoryController() *RepositoryController {
	return &RepositoryController{
		repositoryService: services.NewRepositoryService(),
		validator:         validity.NewValidator(),
	}
}

func (controller *RepositoryController) GetRepository(ctx context.Context, req *registryv1alpha1.GetRepositoryRequest) (*registryv1alpha1.GetRepositoryResponse, e.ResponseError) {
	// 查询
	repository, err := controller.repositoryService.GetRepository(ctx, req.GetId())
	if err != nil {
		logger.Errorf("Error get repo: %v\n", err.Error())

		return nil, err
	}

	repositoryCounts, err := controller.repositoryService.GetRepositoryCounts(ctx, repository.RepositoryID)
	if err != nil {
		logger.Errorf("Error get repo counts: %v\n", err.Error())

		return nil, err
	}

	// 查询成功
	resp := &registryv1alpha1.GetRepositoryResponse{
		Repository: repository.ToProtoRepository(),
		Counts:     repositoryCounts.ToProtoRepositoryCounts(),
	}
	return resp, nil
}

func (controller *RepositoryController) GetRepositoryByFullName(ctx context.Context, req *registryv1alpha1.GetRepositoryByFullNameRequest) (*registryv1alpha1.GetRepositoryByFullNameResponse, e.ResponseError) {
	// 验证参数
	userName, repositoryName, argErr := controller.validator.SplitFullName(req.GetFullName())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 查询
	repository, err := controller.repositoryService.GetRepositoryByUserNameAndRepositoryName(ctx, userName, repositoryName)
	if err != nil {
		logger.Errorf("Error get repo: %v\n", err.Error())

		return nil, err
	}

	repositoryCounts, err := controller.repositoryService.GetRepositoryCounts(ctx, repository.RepositoryID)
	if err != nil {
		logger.Errorf("Error get repo counts: %v\n", err.Error())

		return nil, err
	}

	// 查询成功
	resp := &registryv1alpha1.GetRepositoryByFullNameResponse{
		Repository: repository.ToProtoRepository(),
		Counts:     repositoryCounts.ToProtoRepositoryCounts(),
	}
	return resp, nil
}

func (controller *RepositoryController) ListRepositories(ctx context.Context, req *registryv1alpha1.ListRepositoriesRequest) (*registryv1alpha1.ListRepositoriesResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	repositories, listErr := controller.repositoryService.ListRepositories(ctx, pageTokenChaim.PageOffset, int(req.GetPageSize()), req.Reverse)
	if listErr != nil {
		logger.Errorf("Error list repos: %v\n", listErr.Error())

		return nil, listErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(repositories))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListRepositoriesResponse{
		Repositories:  repositories.ToProtoRepositories(),
		NextPageToken: nextPageToken,
	}
	return resp, nil
}

func (controller *RepositoryController) ListUserRepositories(ctx context.Context, req *registryv1alpha1.ListUserRepositoriesRequest) (*registryv1alpha1.ListUserRepositoriesResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	repositories, listErr := controller.repositoryService.ListUserRepositories(ctx, req.GetUserId(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if err != nil {
		logger.Errorf("Error list user repos: %v\n", listErr.Error())

		return nil, listErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(repositories))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListUserRepositoriesResponse{
		Repositories:  repositories.ToProtoRepositories(),
		NextPageToken: nextPageToken,
	}
	return resp, nil
}

func (controller *RepositoryController) ListRepositoriesUserCanAccess(ctx context.Context, req *registryv1alpha1.ListRepositoriesUserCanAccessRequest) (*registryv1alpha1.ListRepositoriesUserCanAccessResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error parse page token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	userID := ctx.Value(constant.UserIDKey).(string)
	repositories, ListErr := controller.repositoryService.ListRepositoriesUserCanAccess(ctx, userID, pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if err != nil {
		logger.Errorf("Error list repos user can access: %v\n", ListErr.Error())

		return nil, ListErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(repositories))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListRepositoriesUserCanAccessResponse{
		Repositories:  repositories.ToProtoRepositories(),
		NextPageToken: nextPageToken,
	}
	return resp, nil
}

func (controller *RepositoryController) CreateRepositoryByFullName(ctx context.Context, req *registryv1alpha1.CreateRepositoryByFullNameRequest) (*registryv1alpha1.CreateRepositoryByFullNameResponse, e.ResponseError) {
	// 验证参数
	userName, repositoryName, argErr := controller.validator.SplitFullName(req.GetFullName())
	if argErr != nil {
		logger.Errorf("Error check: %v", argErr.Error())

		return nil, argErr
	}
	argErr = controller.validator.CheckRepositoryName(repositoryName)
	if argErr != nil {
		logger.Errorf("Error check: %v", argErr.Error())

		return nil, argErr
	}

	userID := ctx.Value(constant.UserIDKey).(string)

	// 创建
	repository, err := controller.repositoryService.CreateRepositoryByUserNameAndRepositoryName(ctx, userID, userName, repositoryName, req.GetVisibility())
	if err != nil {
		logger.Errorf("Error create repo: %v", err.Error())

		return nil, err
	}

	// 成功
	resp := &registryv1alpha1.CreateRepositoryByFullNameResponse{
		Repository: repository.ToProtoRepository(),
	}
	return resp, nil
}

func (controller *RepositoryController) DeleteRepository(ctx context.Context, req *registryv1alpha1.DeleteRepositoryRequest) (*registryv1alpha1.DeleteRepositoryResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := controller.validator.CheckRepositoryCanDeleteByID(userID, req.GetId(), registryv1alpha1connect.RepositoryServiceDeleteRepositoryProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, permissionErr
	}

	// 查询repository，检查是否可以删除
	err := controller.repositoryService.DeleteRepository(ctx, req.GetId())
	if err != nil {
		logger.Errorf("Error delete repo: %v", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.DeleteRepositoryResponse{}
	return resp, nil
}

func (controller *RepositoryController) DeleteRepositoryByFullName(ctx context.Context, req *registryv1alpha1.DeleteRepositoryByFullNameRequest) (*registryv1alpha1.DeleteRepositoryByFullNameResponse, e.ResponseError) {
	// 验证参数
	userName, repositoryName, argErr := controller.validator.SplitFullName(req.GetFullName())
	if argErr != nil {
		logger.Errorf("Error check: %v\n", argErr.Error())

		return nil, argErr
	}

	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := controller.validator.CheckRepositoryCanDelete(userID, userName, repositoryName, registryv1alpha1connect.RepositoryServiceDeleteRepositoryByFullNameProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, permissionErr
	}

	// 删除
	err := controller.repositoryService.DeleteRepositoryByUserNameAndRepositoryName(ctx, userName, repositoryName)
	if err != nil {
		logger.Errorf("Error delete repo: %v", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.DeleteRepositoryByFullNameResponse{}
	return resp, nil
}

func (controller *RepositoryController) DeprecateRepositoryByName(ctx context.Context, req *registryv1alpha1.DeprecateRepositoryByNameRequest) (*registryv1alpha1.DeprecateRepositoryByNameResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := controller.validator.CheckRepositoryCanEdit(userID, req.GetOwnerName(), req.GetRepositoryName(), registryv1alpha1connect.RepositoryServiceDeprecateRepositoryByNameProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, permissionErr
	}

	// 修改数据库
	updatedRepository, err := controller.repositoryService.DeprecateRepositoryByName(ctx, req.GetOwnerName(), req.GetRepositoryName(), req.GetDeprecationMessage())
	if err != nil {
		logger.Errorf("Error deprecate repo: %v", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.DeprecateRepositoryByNameResponse{
		Repository: updatedRepository.ToProtoRepository(),
	}
	return resp, nil
}

func (controller *RepositoryController) UndeprecateRepositoryByName(ctx context.Context, req *registryv1alpha1.UndeprecateRepositoryByNameRequest) (*registryv1alpha1.UndeprecateRepositoryByNameResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := controller.validator.CheckRepositoryCanEdit(userID, req.GetOwnerName(), req.GetRepositoryName(), registryv1alpha1connect.RepositoryServiceUndeprecateRepositoryByNameProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, permissionErr
	}

	// 修改数据库
	updatedRepository, err := controller.repositoryService.UndeprecateRepositoryByName(ctx, req.GetOwnerName(), req.GetRepositoryName())
	if err != nil {
		logger.Errorf("Error undeprecate repo: %v", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.UndeprecateRepositoryByNameResponse{
		Repository: updatedRepository.ToProtoRepository(),
	}
	return resp, nil
}

func (controller *RepositoryController) UpdateRepositorySettingsByName(ctx context.Context, req *registryv1alpha1.UpdateRepositorySettingsByNameRequest) (*registryv1alpha1.UpdateRepositorySettingsByNameResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	_, permissionErr := controller.validator.CheckRepositoryCanEdit(userID, req.GetOwnerName(), req.GetRepositoryName(), registryv1alpha1connect.RepositoryServiceUpdateRepositorySettingsByNameProcedure)
	if permissionErr != nil {
		logger.Errorf("Error check permission: %v", permissionErr.Error())

		return nil, permissionErr
	}

	// 修改数据库
	err := controller.repositoryService.UpdateRepositorySettingsByName(ctx, req.GetOwnerName(), req.GetRepositoryName(), req.GetVisibility(), req.GetDescription())
	if err != nil {
		logger.Errorf("Error update repo settings: %v", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.UpdateRepositorySettingsByNameResponse{}
	return resp, nil
}

func (controller *RepositoryController) GetRepositoriesByFullName(ctx context.Context, req *registryv1alpha1.GetRepositoriesByFullNameRequest) (*registryv1alpha1.GetRepositoriesByFullNameResponse, e.ResponseError) {

	retRepos := make([]*registryv1alpha1.Repository, 0, len(req.FullNames))
	for _, fullName := range req.GetFullNames() {
		// 验证参数
		userName, repositoryName, argErr := controller.validator.SplitFullName(fullName)
		if argErr != nil {
			logger.Errorf("Error check: %v\n", argErr.Error())

			return nil, argErr
		}

		// 查询
		repository, err := controller.repositoryService.GetRepositoryByUserNameAndRepositoryName(ctx, userName, repositoryName)
		if err != nil {
			logger.Errorf("Error get repo: %v\n", err.Error())

			return nil, err
		}
		retRepos = append(retRepos, repository.ToProtoRepository())
	}

	resp := &registryv1alpha1.GetRepositoriesByFullNameResponse{
		Repositories: retRepos,
	}

	return resp, nil
}

package controllers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
)

type DockerRepoController struct {
	validator         validity.Validator
	dockerRepoService services.DockerRepoService
}

func NewDockerRepoController() *DockerRepoController {
	return &DockerRepoController{
		validator:         validity.NewValidator(),
		dockerRepoService: services.NewDockerRepoService(),
	}
}

func (controller *DockerRepoController) CreateDockerRepo(ctx context.Context, req *registryv1alpha1.CreateDockerRepoRequest) (*registryv1alpha1.CreateDockerRepoResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 检查docker repo name的合法性
	checkErr := controller.validator.CheckDockerRepoName(req.GetName())
	if checkErr != nil {
		logger.Errorf("Error Check Args: %v\n", checkErr.Error())
		return nil, checkErr
	}

	// 检查是否可以登录registry
	checkErr = controller.validator.CheckRegistryAuth(ctx, req.GetAddress(), req.GetUsername(), req.GetPassword())
	if checkErr != nil {
		logger.Errorf("Error Check Registry Login: %v\n", checkErr.Error())
		return nil, checkErr
	}

	// 在数据库中增加
	dockerRepo, err := controller.dockerRepoService.CreateDockerRepo(ctx, userID, req.GetName(), req.GetAddress(), req.GetUsername(), req.GetPassword(), req.GetNote())
	if err != nil {
		logger.Errorf("Error create docker repo: %v\n", err.Error())
		return nil, err
	}

	resp := &registryv1alpha1.CreateDockerRepoResponse{
		DockerRepo: dockerRepo.ToProtoDockerRepo(),
	}
	return resp, nil
}

func (controller *DockerRepoController) GetDockerRepo(ctx context.Context, req *registryv1alpha1.GetDockerRepoRequest) (*registryv1alpha1.GetDockerRepoResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询
	dockerRepo, err := controller.dockerRepoService.GetDockerRepoByID(ctx, req.GetId())
	if err != nil {
		logger.Errorf("Error get docker repo: %v\n", err.Error())
		return nil, err
	}

	// 检查权限
	if dockerRepo.UserID != userID {
		respErr := e.NewPermissionDeniedError("get docker repo")
		logger.Errorf("Error Check Permission: %v\n", respErr.Error())
		return nil, respErr
	}

	resp := &registryv1alpha1.GetDockerRepoResponse{
		DockerRepo: dockerRepo.ToProtoDockerRepo(),
	}

	return resp, nil
}

func (controller *DockerRepoController) GetDockerRepoByName(ctx context.Context, req *registryv1alpha1.GetDockerRepoByNameRequest) (*registryv1alpha1.GetDockerRepoByNameResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询
	dockerRepo, err := controller.dockerRepoService.GetDockerRepoByUserIDAndName(ctx, userID, req.GetName())
	if err != nil {
		logger.Errorf("Error get docker repo: %v\n", err.Error())
		return nil, err
	}

	resp := &registryv1alpha1.GetDockerRepoByNameResponse{
		DockerRepo: dockerRepo.ToProtoDockerRepo(),
	}

	return resp, nil
}

func (controller *DockerRepoController) ListDockerRepos(ctx context.Context, req *registryv1alpha1.ListDockerReposRequest) (*registryv1alpha1.ListDockerReposResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error Check Args: %v\n", argErr.Error())
		return nil, argErr
	}

	// 解析page token
	pageTokenChaim, pageTokenErr := security.ParsePageToken(req.GetPageToken())
	if pageTokenErr != nil {
		logger.Errorf("Error Parse Page Token: %v\n", pageTokenErr.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	// 查询
	dockerRepos, err := controller.dockerRepoService.ListDockerRepos(ctx, userID, pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if err != nil {
		logger.Errorf("Error list docker repo: %v\n", err.Error())
		return nil, err
	}

	// 生成下一页token
	nextPageToken, pageTokenErr := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(dockerRepos))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListDockerReposResponse{
		DockerRepos:   dockerRepos.ToProtoDockerRepos(),
		NextPageToken: nextPageToken,
	}

	return resp, nil
}

func (controller *DockerRepoController) UpdateDockerRepoByName(ctx context.Context, req *registryv1alpha1.UpdateDockerRepoByNameRequest) (*registryv1alpha1.UpdateDockerRepoByNameResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 检查是否可以登录registry
	checkErr := controller.validator.CheckRegistryAuth(ctx, req.GetAddress(), req.GetUsername(), req.GetPassword())
	if checkErr != nil {
		logger.Errorf("Error Check Registry Login: %v\n", checkErr.Error())
		return nil, checkErr
	}

	// 更新
	err := controller.dockerRepoService.UpdateDockerRepoByName(ctx, userID, req.GetName(), req.GetAddress(), req.Username, req.Password)
	if err != nil {
		logger.Errorf("Error update docker repo: %v\n", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.UpdateDockerRepoByNameResponse{}
	return resp, nil
}

func (controller *DockerRepoController) UpdateDockerRepoByID(ctx context.Context, req *registryv1alpha1.UpdateDockerRepoByIDRequest) (*registryv1alpha1.UpdateDockerRepoByIDResponse, e.ResponseError) {
	userID := ctx.Value(constant.UserIDKey).(string)

	dockerRepo, err := controller.dockerRepoService.GetDockerRepoByUserIDAndName(ctx, userID, req.GetId())
	if err != nil {
		logger.Errorf("Error get docker repo: %v\n", err.Error())

		return nil, err
	}

	// 检查权限
	if dockerRepo.UserID != userID {
		logger.Errorf("Error Check Permission: dockerRepo UserID is not equal to current User ID\n")

		respErr := e.NewPermissionDeniedError("update docker repo")
		return nil, respErr
	}

	// 检查是否可以登录registry
	checkErr := controller.validator.CheckRegistryAuth(ctx, req.GetAddress(), req.GetUsername(), req.GetPassword())
	if checkErr != nil {
		logger.Errorf("Error Check Registry Login: %v\n", checkErr.Error())
		return nil, checkErr
	}

	// 更新
	err = controller.dockerRepoService.UpdateDockerRepoByID(ctx, req.GetId(), req.GetAddress(), req.Username, req.Password)
	if err != nil {
		logger.Errorf("Error update docker repo: %v\n", err.Error())

		return nil, err
	}

	resp := &registryv1alpha1.UpdateDockerRepoByIDResponse{}
	return resp, nil
}

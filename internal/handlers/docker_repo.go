package handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/ProtobufMan/bufman/internal/core/security"
	"github.com/ProtobufMan/bufman/internal/core/validity"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type DockerRepoServiceHandler struct {
	validator         validity.Validator
	dockerRepoService services.DockerRepoService
}

func NewDockerRepoServiceHandler() *DockerRepoServiceHandler {
	return &DockerRepoServiceHandler{
		validator:         validity.NewValidator(),
		dockerRepoService: services.NewDockerRepoService(),
	}
}

func (handler *DockerRepoServiceHandler) CreateDockerRepo(ctx context.Context, req *connect.Request[registryv1alpha1.CreateDockerRepoRequest]) (*connect.Response[registryv1alpha1.CreateDockerRepoResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 检查docker repo name的合法性
	checkErr := handler.validator.CheckDockerRepoName(req.Msg.GetName())
	if checkErr != nil {
		logger.Errorf("Error Check Args: %v\n", checkErr.Error())
		return nil, connect.NewError(checkErr.Code(), checkErr)
	}

	// 在数据库中增加
	dockerRepo, err := handler.dockerRepoService.CreateDockerRepo(ctx, userID, req.Msg.GetName(), req.Msg.GetAddress(), req.Msg.GetUsername(), req.Msg.GetPassword(), req.Msg.GetNote())
	if err != nil {
		logger.Errorf("Error create docker repo: %v\n", err.Error())
		return nil, connect.NewError(err.Code(), err)
	}

	resp := connect.NewResponse(&registryv1alpha1.CreateDockerRepoResponse{
		DockerRepo: dockerRepo.ToProtoDockerRepo(),
	})
	return resp, nil
}

func (handler *DockerRepoServiceHandler) GetDockerRepo(ctx context.Context, req *connect.Request[registryv1alpha1.GetDockerRepoRequest]) (*connect.Response[registryv1alpha1.GetDockerRepoResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询
	dockerRepo, err := handler.dockerRepoService.GetDockerRepoByID(ctx, req.Msg.GetId())
	if err != nil {
		logger.Errorf("Error get docker repo: %v\n", err.Error())
		return nil, connect.NewError(err.Code(), err)
	}

	// 检查权限
	if dockerRepo.UserID != userID {
		respErr := e.NewPermissionDeniedError("get docker repo")
		logger.Errorf("Error Check Permission: %v\n", respErr.Error())
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetDockerRepoResponse{
		DockerRepo: dockerRepo.ToProtoDockerRepo(),
	})

	return resp, nil
}

func (handler *DockerRepoServiceHandler) GetDockerRepoByName(ctx context.Context, req *connect.Request[registryv1alpha1.GetDockerRepoByNameRequest]) (*connect.Response[registryv1alpha1.GetDockerRepoByNameResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 查询
	dockerRepo, err := handler.dockerRepoService.GetDockerRepoByUserIDAndName(ctx, userID, req.Msg.GetName())
	if err != nil {
		logger.Errorf("Error get docker repo: %v\n", err.Error())
		return nil, connect.NewError(err.Code(), err)
	}

	resp := connect.NewResponse(&registryv1alpha1.GetDockerRepoByNameResponse{
		DockerRepo: dockerRepo.ToProtoDockerRepo(),
	})

	return resp, nil
}

func (handler *DockerRepoServiceHandler) ListDockerRepos(ctx context.Context, req *connect.Request[registryv1alpha1.ListDockerReposRequest]) (*connect.Response[registryv1alpha1.ListDockerReposResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error Check Args: %v\n", argErr.Error())
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 解析page token
	pageTokenChaim, pageTokenErr := security.ParsePageToken(req.Msg.GetPageToken())
	if pageTokenErr != nil {
		logger.Errorf("Error Parse Page Token: %v\n", pageTokenErr.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 查询
	dockerRepos, err := handler.dockerRepoService.ListDockerRepos(ctx, userID, pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if err != nil {
		logger.Errorf("Error list docker repo: %v\n", err.Error())
		return nil, connect.NewError(err.Code(), err)
	}

	// 生成下一页token
	nextPageToken, pageTokenErr := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(dockerRepos))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate next page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.ListDockerReposResponse{
		DockerRepos:   dockerRepos.ToProtoDockerRepos(),
		NextPageToken: nextPageToken,
	})

	return resp, nil
}

func (handler *DockerRepoServiceHandler) UpdateDockerRepoByName(ctx context.Context, req *connect.Request[registryv1alpha1.UpdateDockerRepoByNameRequest]) (*connect.Response[registryv1alpha1.UpdateDockerRepoByNameResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	// 更新
	err := handler.dockerRepoService.UpdateDockerRepoByName(ctx, userID, req.Msg.GetName(), req.Msg.GetAddress(), req.Msg.Username, req.Msg.Password)
	if err != nil {
		logger.Errorf("Error update docker repo: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err)
	}

	resp := connect.NewResponse(&registryv1alpha1.UpdateDockerRepoByNameResponse{})
	return resp, nil
}

func (handler *DockerRepoServiceHandler) UpdateDockerRepoByID(ctx context.Context, req *connect.Request[registryv1alpha1.UpdateDockerRepoByIDRequest]) (*connect.Response[registryv1alpha1.UpdateDockerRepoByIDResponse], error) {
	userID := ctx.Value(constant.UserIDKey).(string)

	dockerRepo, err := handler.dockerRepoService.GetDockerRepoByUserIDAndName(ctx, userID, req.Msg.GetId())
	if err != nil {
		logger.Errorf("Error get docker repo: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err)
	}

	// 检查权限
	if dockerRepo.UserID != userID {
		logger.Errorf("Error Check Permission: dockerRepo UserID is not equal to current User ID\n")

		respErr := e.NewPermissionDeniedError("update docker repo")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 更新
	err = handler.dockerRepoService.UpdateDockerRepoByID(ctx, req.Msg.GetId(), req.Msg.GetAddress(), req.Msg.Username, req.Msg.Password)
	if err != nil {
		logger.Errorf("Error update docker repo: %v\n", err.Error())

		return nil, connect.NewError(err.Code(), err)
	}

	resp := connect.NewResponse(&registryv1alpha1.UpdateDockerRepoByIDResponse{})
	return resp, nil
}

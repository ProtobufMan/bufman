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

type CommitController struct {
	commitService services.CommitService
	validator     validity.Validator
}

func NewCommitController() *CommitController {
	return &CommitController{
		commitService: services.NewCommitService(),
		validator:     validity.NewValidator(),
	}
}

func (controller *CommitController) ListRepositoryCommitsByReference(ctx context.Context, req *registryv1alpha1.ListRepositoryCommitsByReferenceRequest) (*registryv1alpha1.ListRepositoryCommitsByReferenceResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error Check Args: %v\n", argErr.Error())
		return nil, argErr
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := controller.validator.CheckRepositoryCanAccess(userID, req.GetRepositoryOwner(), req.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure)
	if permissionErr != nil {
		logger.Errorf("Error Check Permission: %v\n", permissionErr.Error())
		return nil, permissionErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error Parse Page Token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	// 查询
	commits, respErr := controller.commitService.ListRepositoryCommitsByReference(ctx, repository.RepositoryID, req.GetReference(), pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error list repository commits: %v\n", respErr.Error())
		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(commits))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListRepositoryCommitsByReferenceResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
		NextPageToken:     nextPageToken,
	}
	return resp, nil
}

func (controller *CommitController) GetRepositoryCommitByReference(ctx context.Context, req *registryv1alpha1.GetRepositoryCommitByReferenceRequest) (*registryv1alpha1.GetRepositoryCommitByReferenceResponse, e.ResponseError) {
	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := controller.validator.CheckRepositoryCanAccess(userID, req.GetRepositoryOwner(), req.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure)
	if permissionErr != nil {
		logger.Errorf("Error Check Permission: %v\n", permissionErr.Error())
		return nil, permissionErr
	}

	// 查询
	commit, respErr := controller.commitService.GetRepositoryCommitByReference(ctx, repository.RepositoryID, req.GetReference())
	if respErr != nil {
		logger.Errorf("Error list repository commits: %v\n", respErr.Error())
		return nil, respErr
	}

	resp := &registryv1alpha1.GetRepositoryCommitByReferenceResponse{
		RepositoryCommit: commit.ToProtoRepositoryCommit(),
	}
	return resp, nil
}

func (controller *CommitController) ListRepositoryDraftCommits(ctx context.Context, req *registryv1alpha1.ListRepositoryDraftCommitsRequest) (*registryv1alpha1.ListRepositoryDraftCommitsResponse, e.ResponseError) {
	// 验证参数
	argErr := controller.validator.CheckPageSize(req.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error Check Args: %v\n", argErr.Error())
		return nil, argErr
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := controller.validator.CheckRepositoryCanAccess(userID, req.GetRepositoryOwner(), req.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceListRepositoryDraftCommitsProcedure)
	if permissionErr != nil {
		logger.Errorf("Error Check Permission: %v\n", permissionErr.Error())
		return nil, permissionErr
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.GetPageToken())
	if err != nil {
		logger.Errorf("Error Parse Page Token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, respErr
	}

	// 查询
	commits, respErr := controller.commitService.ListRepositoryDraftCommits(ctx, repository.RepositoryID, pageTokenChaim.PageOffset, int(req.GetPageSize()), req.GetReverse())
	if respErr != nil {
		logger.Errorf("Error list repository draft commits: %v\n", respErr.Error())
		return nil, respErr
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.GetPageSize()), len(commits))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate page token")
		return nil, respErr
	}

	resp := &registryv1alpha1.ListRepositoryDraftCommitsResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
		NextPageToken:     nextPageToken,
	}
	return resp, nil
}

func (controller *CommitController) DeleteRepositoryDraftCommit(ctx context.Context, req *registryv1alpha1.DeleteRepositoryDraftCommitRequest) (*registryv1alpha1.DeleteRepositoryDraftCommitResponse, e.ResponseError) {
	// 获取user ID
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := controller.validator.CheckRepositoryCanEdit(userID, req.GetRepositoryOwner(), req.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure)
	if permissionErr != nil {
		logger.Errorf("Error Check Permission: %v\n", permissionErr.Error())
		return nil, permissionErr
	}

	// 删除
	err := controller.commitService.DeleteRepositoryDraftCommit(ctx, repository.RepositoryID, req.GetDraftName())
	if err != nil {
		logger.Errorf("Error delete repository draft commits: %v\n", err.Error())
		return nil, err
	}

	resp := &registryv1alpha1.DeleteRepositoryDraftCommitResponse{}
	return resp, nil
}

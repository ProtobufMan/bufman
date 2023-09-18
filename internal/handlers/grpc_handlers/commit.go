package grpc_handlers

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
	"github.com/bufbuild/connect-go"
)

type CommitServiceHandler struct {
	commitService services.CommitService
	validator     validity.Validator
}

func NewCommitServiceHandler() *CommitServiceHandler {
	return &CommitServiceHandler{
		commitService: services.NewCommitService(),
		validator:     validity.NewValidator(),
	}
}

func (handler *CommitServiceHandler) ListRepositoryCommitsByReference(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryCommitsByReferenceRequest]) (*connect.Response[registryv1alpha1.ListRepositoryCommitsByReferenceResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error Check Args: %v\n", argErr.Error())
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure)
	if permissionErr != nil {
		logger.Errorf("Error Check Permission: %v\n", permissionErr.Error())
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error Parse Page Token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 查询
	commits, respErr := handler.commitService.ListRepositoryCommitsByReference(ctx, repository.RepositoryID, req.Msg.GetReference(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error list repository commits: %v\n", respErr.Error())
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(commits))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.ListRepositoryCommitsByReferenceResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
		NextPageToken:     nextPageToken,
	})
	return resp, nil
}

func (handler *CommitServiceHandler) GetRepositoryCommitByReference(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryCommitByReferenceRequest]) (*connect.Response[registryv1alpha1.GetRepositoryCommitByReferenceResponse], error) {
	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure)
	if permissionErr != nil {
		logger.Errorf("Error Check Permission: %v\n", permissionErr.Error())
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 查询
	commit, respErr := handler.commitService.GetRepositoryCommitByReference(ctx, repository.RepositoryID, req.Msg.GetReference())
	if respErr != nil {
		logger.Errorf("Error list repository commits: %v\n", respErr.Error())
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.GetRepositoryCommitByReferenceResponse{
		RepositoryCommit: commit.ToProtoRepositoryCommit(),
	})
	return resp, nil
}

func (handler *CommitServiceHandler) ListRepositoryDraftCommits(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryDraftCommitsRequest]) (*connect.Response[registryv1alpha1.ListRepositoryDraftCommitsResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		logger.Errorf("Error Check Args: %v\n", argErr.Error())
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceListRepositoryDraftCommitsProcedure)
	if permissionErr != nil {
		logger.Errorf("Error Check Permission: %v\n", permissionErr.Error())
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := security.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		logger.Errorf("Error Parse Page Token: %v\n", err.Error())

		respErr := e.NewInvalidArgumentError("page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	// 查询
	commits, respErr := handler.commitService.ListRepositoryDraftCommits(ctx, repository.RepositoryID, pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		logger.Errorf("Error list repository draft commits: %v\n", respErr.Error())
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	// 生成下一页token
	nextPageToken, err := security.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(commits))
	if err != nil {
		logger.Errorf("Error generate next page token: %v\n", err.Error())

		respErr := e.NewInternalError("generate page token")
		return nil, connect.NewError(respErr.Code(), respErr)
	}

	resp := connect.NewResponse(&registryv1alpha1.ListRepositoryDraftCommitsResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
		NextPageToken:     nextPageToken,
	})
	return resp, nil
}

func (handler *CommitServiceHandler) DeleteRepositoryDraftCommit(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteRepositoryDraftCommitRequest]) (*connect.Response[registryv1alpha1.DeleteRepositoryDraftCommitResponse], error) {
	// 获取user ID
	userID := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanEdit(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure)
	if permissionErr != nil {
		logger.Errorf("Error Check Permission: %v\n", permissionErr.Error())
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 删除
	err := handler.commitService.DeleteRepositoryDraftCommit(ctx, repository.RepositoryID, req.Msg.GetDraftName())
	if err != nil {
		logger.Errorf("Error delete repository draft commits: %v\n", err.Error())
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha1.DeleteRepositoryDraftCommitResponse{})
	return resp, nil
}

func (handler *CommitServiceHandler) ListRepositoryCommitsByBranch(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryCommitsByBranchRequest]) (*connect.Response[registryv1alpha1.ListRepositoryCommitsByBranchResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *CommitServiceHandler) GetRepositoryCommitBySequenceId(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryCommitBySequenceIdRequest]) (*connect.Response[registryv1alpha1.GetRepositoryCommitBySequenceIdResponse], error) {
	//TODO implement me
	panic("implement me")
}

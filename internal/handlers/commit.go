package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/ProtobufMan/bufman/internal/util"
	"github.com/ProtobufMan/bufman/internal/validity"
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
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := util.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询
	commits, respErr := handler.commitService.ListRepositoryCommitsByReference(repository.RepositoryID, req.Msg.GetReference(), pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	// 生成下一页token
	nextPageToken, err := util.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(commits))
	if err != nil {
		return nil, e.NewInternalError("generate next page token")
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
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 查询
	commit, respErr := handler.commitService.GetRepositoryCommitByReference(repository.RepositoryID, req.Msg.GetReference())
	if respErr != nil {
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
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alpha1connect.RepositoryCommitServiceListRepositoryDraftCommitsProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 解析page token
	pageTokenChaim, err := util.ParsePageToken(req.Msg.GetPageToken())
	if err != nil {
		return nil, e.NewInvalidArgumentError("page token")
	}

	// 查询
	commits, respErr := handler.commitService.ListRepositoryDraftCommits(repository.RepositoryID, pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	// 生成下一页token
	nextPageToken, err := util.GenerateNextPageToken(pageTokenChaim.PageOffset, int(req.Msg.GetPageSize()), len(commits))
	if err != nil {
		return nil, e.NewInternalError("generate next page token")
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
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 删除
	err := handler.commitService.DeleteRepositoryDraftCommit(repository.RepositoryID, req.Msg.GetDraftName())
	if err != nil {
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

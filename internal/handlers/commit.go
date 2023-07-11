package handlers

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	registryv1alpha "github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/services"
	"github.com/bufbuild/connect-go"
)

type CommitServiceHandler struct {
	commitService services.CommitService
}

func NewCommitServiceHandler() *CommitServiceHandler {
	return &CommitServiceHandler{
		commitService: services.NewCommitService(),
	}
}

func (handler *CommitServiceHandler) ListRepositoryCommitsByReference(ctx context.Context, req *connect.Request[registryv1alpha.ListRepositoryCommitsByReferenceRequest]) (*connect.Response[registryv1alpha.ListRepositoryCommitsByReferenceResponse], error) {
	var commits model.Commits
	var respErr e.ResponseError

	// 尝试获取user ID
	userID, ok := ctx.Value(constant.UserIDKey).(string)
	if ok {
		// 带有 user id
		commits, respErr = handler.commitService.ListRepositoryCommitsByReferenceWithUserID(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), req.Msg.GetReference(), int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	} else {
		// 不带 user id
		commits, respErr = handler.commitService.ListRepositoryCommitsByReference(req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), req.Msg.GetReference(), int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	}

	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListRepositoryCommitsByReferenceResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
	})
	return resp, nil
}

func (handler *CommitServiceHandler) GetRepositoryCommitByReference(ctx context.Context, req *connect.Request[registryv1alpha.GetRepositoryCommitByReferenceRequest]) (*connect.Response[registryv1alpha.GetRepositoryCommitByReferenceResponse], error) {
	var commit *model.Commit
	var respErr e.ResponseError

	// 尝试获取user ID
	userID, ok := ctx.Value(constant.UserIDKey).(string)
	if ok {
		commit, respErr = handler.commitService.GetRepositoryCommitByReferenceWithUserID(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), req.Msg.GetReference())
	} else {
		commit, respErr = handler.commitService.GetRepositoryCommitByReference(req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), req.Msg.GetReference())
	}

	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.GetRepositoryCommitByReferenceResponse{
		RepositoryCommit: commit.ToProtoRepositoryCommit(),
	})
	return resp, nil
}

func (handler *CommitServiceHandler) ListRepositoryDraftCommits(ctx context.Context, req *connect.Request[registryv1alpha.ListRepositoryDraftCommitsRequest]) (*connect.Response[registryv1alpha.ListRepositoryDraftCommitsResponse], error) {
	var commits model.Commits
	var respErr e.ResponseError

	// 尝试获取user ID
	userID, ok := ctx.Value(constant.UserIDKey).(string)
	if ok {
		// 带有 user id
		commits, respErr = handler.commitService.ListRepositoryDraftCommitsWithUserID(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	} else {
		// 不带 user id
		commits, respErr = handler.commitService.ListRepositoryDraftCommits(req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	}

	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListRepositoryDraftCommitsResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
	})
	return resp, nil
}

func (handler *CommitServiceHandler) DeleteRepositoryDraftCommit(ctx context.Context, req *connect.Request[registryv1alpha.DeleteRepositoryDraftCommitRequest]) (*connect.Response[registryv1alpha.DeleteRepositoryDraftCommitResponse], error) {
	// 获取user ID
	userID := ctx.Value(constant.UserIDKey).(string)

	// 删除
	err := handler.commitService.DeleteRepositoryDraftCommit(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), req.Msg.GetDraftName())
	if err != nil {
		return nil, connect.NewError(err.Code(), err.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.DeleteRepositoryDraftCommitResponse{})
	return resp, nil
}

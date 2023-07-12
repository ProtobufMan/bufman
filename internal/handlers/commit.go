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

func (handler *CommitServiceHandler) ListRepositoryCommitsByReference(ctx context.Context, req *connect.Request[registryv1alpha.ListRepositoryCommitsByReferenceRequest]) (*connect.Response[registryv1alpha.ListRepositoryCommitsByReferenceResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alphaconnect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 查询
	commits, respErr := handler.commitService.ListRepositoryCommitsByReference(repository.RepositoryID, req.Msg.GetReference(), int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.ListRepositoryCommitsByReferenceResponse{
		RepositoryCommits: commits.ToProtoRepositoryCommits(),
	})
	return resp, nil
}

func (handler *CommitServiceHandler) GetRepositoryCommitByReference(ctx context.Context, req *connect.Request[registryv1alpha.GetRepositoryCommitByReferenceRequest]) (*connect.Response[registryv1alpha.GetRepositoryCommitByReferenceResponse], error) {
	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alphaconnect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 查询
	commit, respErr := handler.commitService.GetRepositoryCommitByReference(repository.RepositoryID, req.Msg.GetReference())
	if respErr != nil {
		return nil, connect.NewError(respErr.Code(), respErr.Err())
	}

	resp := connect.NewResponse(&registryv1alpha.GetRepositoryCommitByReferenceResponse{
		RepositoryCommit: commit.ToProtoRepositoryCommit(),
	})
	return resp, nil
}

func (handler *CommitServiceHandler) ListRepositoryDraftCommits(ctx context.Context, req *connect.Request[registryv1alpha.ListRepositoryDraftCommitsRequest]) (*connect.Response[registryv1alpha.ListRepositoryDraftCommitsResponse], error) {
	// 验证参数
	argErr := handler.validator.CheckPageSize(req.Msg.GetPageSize())
	if argErr != nil {
		return nil, connect.NewError(argErr.Code(), argErr.Err())
	}

	// 尝试获取user ID
	userID, _ := ctx.Value(constant.UserIDKey).(string)

	// 验证用户权限
	repository, permissionErr := handler.validator.CheckRepositoryCanAccess(userID, req.Msg.GetRepositoryOwner(), req.Msg.GetRepositoryName(), registryv1alphaconnect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure)
	if permissionErr != nil {
		return nil, connect.NewError(permissionErr.Code(), permissionErr.Err())
	}

	// 查询
	commits, respErr := handler.commitService.ListRepositoryDraftCommits(repository.RepositoryID, int(req.Msg.GetPageOffset()), int(req.Msg.GetPageSize()), req.Msg.GetReverse())
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

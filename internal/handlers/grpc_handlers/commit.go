package grpc_handlers

import (
	"context"
	registryv1alpha1 "github.com/ProtobufMan/bufman-cli/private/gen/proto/go/bufman/alpha/registry/v1alpha1"
	"github.com/ProtobufMan/bufman/internal/controllers"
	"github.com/bufbuild/connect-go"
)

type CommitServiceHandler struct {
	commitController *controllers.CommitController
}

func NewCommitServiceHandler() *CommitServiceHandler {
	return &CommitServiceHandler{
		commitController: controllers.NewCommitController(),
	}
}

func (handler *CommitServiceHandler) ListRepositoryCommitsByReference(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryCommitsByReferenceRequest]) (*connect.Response[registryv1alpha1.ListRepositoryCommitsByReferenceResponse], error) {
	resp, err := handler.commitController.ListRepositoryCommitsByReference(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *CommitServiceHandler) GetRepositoryCommitByReference(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryCommitByReferenceRequest]) (*connect.Response[registryv1alpha1.GetRepositoryCommitByReferenceResponse], error) {
	resp, err := handler.commitController.GetRepositoryCommitByReference(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *CommitServiceHandler) ListRepositoryDraftCommits(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryDraftCommitsRequest]) (*connect.Response[registryv1alpha1.ListRepositoryDraftCommitsResponse], error) {
	resp, err := handler.commitController.ListRepositoryDraftCommits(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *CommitServiceHandler) DeleteRepositoryDraftCommit(ctx context.Context, req *connect.Request[registryv1alpha1.DeleteRepositoryDraftCommitRequest]) (*connect.Response[registryv1alpha1.DeleteRepositoryDraftCommitResponse], error) {
	resp, err := handler.commitController.DeleteRepositoryDraftCommit(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(err.Code(), err)
	}

	return connect.NewResponse(resp), nil
}

func (handler *CommitServiceHandler) ListRepositoryCommitsByBranch(ctx context.Context, req *connect.Request[registryv1alpha1.ListRepositoryCommitsByBranchRequest]) (*connect.Response[registryv1alpha1.ListRepositoryCommitsByBranchResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (handler *CommitServiceHandler) GetRepositoryCommitBySequenceId(ctx context.Context, req *connect.Request[registryv1alpha1.GetRepositoryCommitBySequenceIdRequest]) (*connect.Response[registryv1alpha1.GetRepositoryCommitBySequenceIdResponse], error) {
	//TODO implement me
	panic("implement me")
}

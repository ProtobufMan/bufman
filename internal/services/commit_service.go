package services

import (
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/gorm"
)

type CommitService interface {
	ListRepositoryCommitsByReference(repositoryID, reference string, offset, limit int, reverse bool) (model.Commits, e.ResponseError)
	GetRepositoryCommitByReference(repositoryID, reference string) (*model.Commit, e.ResponseError)
	ListRepositoryDraftCommits(repositoryID string, offset, limit int, reverse bool) (model.Commits, e.ResponseError)
	DeleteRepositoryDraftCommit(repositoryID, draftName string) e.ResponseError
}

type CommitServiceImpl struct {
	repositoryMapper mapper.RepositoryMapper
	commitMapper     mapper.CommitMapper
}

func NewCommitService() CommitService {
	return &CommitServiceImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
	}
}

func (commitService *CommitServiceImpl) ListRepositoryCommitsByReference(repositoryID, reference string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	// 查询commits
	commits, err := commitService.commitMapper.FindPageByRepositoryIDAndReference(repositoryID, reference, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure)
	}

	return commits, nil
}

func (commitService *CommitServiceImpl) GetRepositoryCommitByReference(repositoryID, reference string) (*model.Commit, e.ResponseError) {
	// 查询commit
	commit, err := commitService.commitMapper.FindByRepositoryIDAndReference(repositoryID, reference)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("commit")
		}

		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure)
	}

	return commit, nil
}

func (commitService *CommitServiceImpl) ListRepositoryDraftCommits(repositoryID string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	var commits model.Commits
	var err error
	// 查询draft
	commits, err = commitService.commitMapper.FindDraftPageByRepositoryID(repositoryID, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.RepositoryCommitServiceListRepositoryDraftCommitsProcedure)
	}

	return commits, nil
}

func (commitService *CommitServiceImpl) DeleteRepositoryDraftCommit(repositoryID, draftName string) e.ResponseError {
	// 删除
	err := commitService.commitMapper.DeleteByRepositoryIDAndDraftName(repositoryID, draftName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("draft")
		}

		return e.NewInternalError(registryv1alpha1connect.RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure)
	}

	return nil
}

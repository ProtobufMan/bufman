package services

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/gen/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/validity"
	"gorm.io/gorm"
)

type CommitService interface {
	ListRepositoryCommitsByReference(repositoryID, reference string, offset, limit int, reverse bool) (model.Commits, e.ResponseError)
	GetRepositoryCommitByReference(repositoryID, reference string) (*model.Commit, e.ResponseError)
	ListRepositoryDraftCommits(repositoryID string, offset, limit int, reverse bool) (model.Commits, e.ResponseError)
	DeleteRepositoryDraftCommit(userID, ownerName, repositoryName, draftName string) e.ResponseError
}

type CommitServiceImpl struct {
	repositoryMapper mapper.RepositoryMapper
	commitMapper     mapper.CommitMapper
	validator        validity.Validator
}

func NewCommitService() CommitService {
	return &CommitServiceImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
		validator:        validity.NewValidator(),
	}
}

func (commitService *CommitServiceImpl) ListRepositoryCommitsByReference(repositoryID, reference string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	return commitService.doListRepositoryCommitsByReference(repositoryID, reference, offset, limit, reverse)
}

func (commitService *CommitServiceImpl) GetRepositoryCommitByReference(repositoryID, reference string) (*model.Commit, e.ResponseError) {
	// 查询commit
	return commitService.doGetRepositoryCommitByReference(repositoryID, reference)
}

func (commitService *CommitServiceImpl) ListRepositoryDraftCommits(repositoryID string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	// 查询 draft
	return commitService.doListRepositoryDraftCommit(repositoryID, offset, limit, reverse)
}

func (commitService *CommitServiceImpl) DeleteRepositoryDraftCommit(userID, ownerName, repositoryName, draftName string) e.ResponseError {
	repository, err := commitService.repositoryMapper.FindByUserNameAndRepositoryName(ownerName, repositoryName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("repository")
		}

		return e.NewInternalError(registryv1alphaconnect.RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure)
	}

	// 不可以删除别人的draft
	if repository.UserID != userID {
		return e.NewPermissionDeniedError(registryv1alphaconnect.RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure)
	}

	// 删除
	err = commitService.commitMapper.DeleteByRepositoryIDAndDraftName(repository.RepositoryID, draftName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("draft")
		}

		return e.NewInternalError(registryv1alphaconnect.RepositoryCommitServiceDeleteRepositoryDraftCommitProcedure)
	}

	return nil
}

func (commitService *CommitServiceImpl) doListRepositoryCommitsByReference(repositoryID string, reference string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	var commits model.Commits
	var err error
	if reference == "" || reference == constant.DefaultBranch {
		commits, err = commitService.commitMapper.FindPageByRepositoryID(repositoryID, offset, limit, reverse)
	} else if len(reference) == constant.CommitLength {
		commits, err = commitService.commitMapper.FindPageByRepositoryIDAndCommitName(repositoryID, reference, offset, limit, reverse)
	} else {
		commits, err = commitService.commitMapper.FindPageByRepositoryIDAndTagName(repositoryID, reference, offset, limit, reverse)
	}

	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryCommitServiceListRepositoryCommitsByReferenceProcedure)
	}

	return commits, nil
}

func (commitService *CommitServiceImpl) doGetRepositoryCommitByReference(repositoryID string, reference string) (*model.Commit, e.ResponseError) {
	var commit *model.Commit
	var err error
	// repo is public
	if len(reference) == constant.CommitLength {
		// 查询commit
		commit, err = commitService.commitMapper.FindByRepositoryIDAndCommitName(repositoryID, reference)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, e.NewNotFoundError("commit")
			}
			return nil, e.NewInternalError(registryv1alphaconnect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure)
		}
	}

	// 查询tag
	commit, err = commitService.commitMapper.FindByRepositoryIDAndTagName(repositoryID, reference)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询draft
			commit, err = commitService.commitMapper.FindByRepositoryIDAndDraftName(repositoryID, reference)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, e.NewNotFoundError("commit")
			}

			return nil, e.NewInternalError(registryv1alphaconnect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure)
		}
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryCommitServiceGetRepositoryCommitByReferenceProcedure)
	}

	return commit, nil
}

func (commitService *CommitServiceImpl) doListRepositoryDraftCommit(repositoryID string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	var commits model.Commits
	var err error
	// 查询draft
	commits, err = commitService.commitMapper.FindDraftPageByRepositoryID(repositoryID, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(registryv1alphaconnect.RepositoryCommitServiceListRepositoryDraftCommitsProcedure)
	}

	return commits, nil
}

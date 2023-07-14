package services

import (
	"errors"
	"fmt"
	"github.com/ProtobufMan/bufman/internal/e"
	modulev1alpha "github.com/ProtobufMan/bufman/internal/gen/bufman/module/v1alpha"
	"github.com/ProtobufMan/bufman/internal/gen/bufman/registry/v1alpha/registryv1alphaconnect"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/gorm"
)

type ResolveService interface {
	GetModulePins(repositoryMap map[string]*model.Repository, moduleReferences []*modulev1alpha.ModuleReference) (model.Commits, e.ResponseError)
}

type ResolveServiceImpl struct {
	repositoryMapper mapper.RepositoryMapper
	commitMapper     mapper.CommitMapper
}

func NewResolveService() ResolveService {
	return &ResolveServiceImpl{
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
	}
}

func (resolveService *ResolveServiceImpl) GetModulePins(repositoryMap map[string]*model.Repository, moduleReferences []*modulev1alpha.ModuleReference) (model.Commits, e.ResponseError) {
	commits := make([]*model.Commit, 0, len(moduleReferences))
	for _, moduleReference := range moduleReferences {
		fullName := moduleReference.GetOwner() + "/" + moduleReference.GetRepository()
		repo, ok := repositoryMap[fullName]
		if !ok {
			return nil, e.NewInternalError(registryv1alphaconnect.ResolveServiceGetModulePinsProcedure)
		}

		// 查询reference
		commit, err := resolveService.commitMapper.FindByRepositoryIDAndReference(repo.RepositoryID, moduleReference.GetReference())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, e.NewNotFoundError(fmt.Sprintf("reference %s", moduleReference.GetReference()))
			}

			return nil, e.NewInternalError(registryv1alphaconnect.ResolveServiceGetModulePinsProcedure)
		}

		commits = append(commits, commit)
	}

	// TODO 加入依赖

	return commits, nil
}

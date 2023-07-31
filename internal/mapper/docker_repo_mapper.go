package mapper

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
)

type DockerRepoMapper interface {
	Create(dockerRepo *model.DockerRepo) error
	FindByUserIDAndName(userID string, name string) (*model.DockerRepo, error)
	FindPageByUserID(userID string, offset int, limit int, reverse bool) (model.DockerRepos, error)
	UpdateByUserIDAndName(userID, name string, dockerRepo *model.DockerRepo) error
}

type DockerRepoMapperImpl struct{}

func (d *DockerRepoMapperImpl) Create(dockerRepo *model.DockerRepo) error {
	return dal.DockerRepo.Create(dockerRepo)
}

func (d *DockerRepoMapperImpl) FindByUserIDAndName(userID string, name string) (*model.DockerRepo, error) {
	return dal.DockerRepo.Where(dal.DockerRepo.UserID.Eq(userID), dal.DockerRepo.Name.Eq(name)).First()
}

func (d *DockerRepoMapperImpl) FindPageByUserID(userID string, offset int, limit int, reverse bool) (model.DockerRepos, error) {
	stmt := dal.DockerRepo.Where(dal.DockerRepo.UserID.Eq(userID)).Offset(offset).Limit(limit)
	if reverse {
		stmt = stmt.Order(dal.DockerRepo.ID.Desc())
	}

	return stmt.Find()
}

func (d *DockerRepoMapperImpl) UpdateByUserIDAndName(userID, name string, dockerRepo *model.DockerRepo) error {
	_, err := dal.DockerRepo.Select(dal.DockerRepo.Address, dal.DockerRepo.UserName, dal.DockerRepo.Password).Where(dal.DockerRepo.UserID.Eq(userID), dal.DockerRepo.Name.Eq(name)).Updates(dockerRepo)
	return err
}

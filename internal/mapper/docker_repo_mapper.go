package mapper

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
)

type DockerRepoMapper interface {
	Create(dockerRepo *model.DockerRepo) error
	FindByDockerRepoID(dockerRepoID string) (*model.DockerRepo, error)
	FindByUserIDAndName(userID string, dockerRepoName string) (*model.DockerRepo, error)
	FindPageByUserID(userID string, offset int, limit int, reverse bool) (model.DockerRepos, error)
	UpdateByUserIDAndName(userID, dockerRepoName string, dockerRepo *model.DockerRepo) error
	UpdateByDockerRepoID(dockerRepoID string, dockerRepo *model.DockerRepo) error
}

type DockerRepoMapperImpl struct{}

func (d *DockerRepoMapperImpl) Create(dockerRepo *model.DockerRepo) error {
	return dal.DockerRepo.Create(dockerRepo)
}

func (d *DockerRepoMapperImpl) FindByDockerRepoID(dockerRepoID string) (*model.DockerRepo, error) {
	return dal.DockerRepo.Where(dal.DockerRepo.DockerRepoID.Eq(dockerRepoID)).First()
}

func (d *DockerRepoMapperImpl) FindByUserIDAndName(userID string, dockerRepoName string) (*model.DockerRepo, error) {
	return dal.DockerRepo.Where(dal.DockerRepo.UserID.Eq(userID), dal.DockerRepo.DockerRepoName.Eq(dockerRepoName)).First()
}

func (d *DockerRepoMapperImpl) FindPageByUserID(userID string, offset int, limit int, reverse bool) (model.DockerRepos, error) {
	stmt := dal.DockerRepo.Where(dal.DockerRepo.UserID.Eq(userID)).Offset(offset).Limit(limit)
	if reverse {
		stmt = stmt.Order(dal.DockerRepo.ID.Desc())
	}

	return stmt.Find()
}

func (d *DockerRepoMapperImpl) UpdateByUserIDAndName(userID, dockerRepoName string, dockerRepo *model.DockerRepo) error {
	_, err := dal.DockerRepo.Select(dal.DockerRepo.Address, dal.DockerRepo.UserName, dal.DockerRepo.Password).Where(dal.DockerRepo.UserID.Eq(userID), dal.DockerRepo.DockerRepoName.Eq(dockerRepoName)).Updates(dockerRepo)
	return err
}

func (d *DockerRepoMapperImpl) UpdateByDockerRepoID(dockerRepoID string, dockerRepo *model.DockerRepo) error {
	_, err := dal.DockerRepo.Select(dal.DockerRepo.Address, dal.DockerRepo.UserName, dal.DockerRepo.Password).Where(dal.DockerRepo.DockerRepoID.Eq(dockerRepoID)).Updates(dockerRepo)
	return err
}

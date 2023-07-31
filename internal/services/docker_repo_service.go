package services

import (
	"context"
	"errors"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DockerRepoService interface {
	CreateDockerRepo(ctx context.Context, userID, dockerRepoName, address, userName, password, note string) (*model.DockerRepo, e.ResponseError)
	GetDockerRepoByID(ctx context.Context, dockerRepoID string) (*model.DockerRepo, e.ResponseError)
	GetDockerRepoByUserIDAndName(ctx context.Context, userID, dockerRepoName string) (*model.DockerRepo, e.ResponseError)
	ListDockerRepos(ctx context.Context, userID string, offset, limit int, reverse bool) (model.DockerRepos, e.ResponseError)
	UpdateDockerRepoByID(ctx context.Context, dockerRepoID, address, userName, password string) e.ResponseError
	UpdateDockerRepoByName(ctx context.Context, userID, dockerRepoName, address, userName, password string) e.ResponseError
}

type DockerRepoServiceImpl struct {
	dockerRepoMapper mapper.DockerRepoMapper
}

func NewDockerRepoService() DockerRepoService {
	return &DockerRepoServiceImpl{
		dockerRepoMapper: &mapper.DockerRepoMapperImpl{},
	}
}

func (dockerRepoService *DockerRepoServiceImpl) CreateDockerRepo(ctx context.Context, userID, dockerRepoName, address, userName, password, note string) (*model.DockerRepo, e.ResponseError) {
	dockerRepo := &model.DockerRepo{
		UserID:         userID,
		DockerRepoID:   uuid.NewString(),
		DockerRepoName: dockerRepoName,
		Address:        address,
		UserName:       userName,
		Password:       password,
		Note:           note,
	}

	// 在数据库中记录
	err := dockerRepoService.dockerRepoMapper.Create(dockerRepo)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, e.NewAlreadyExistsError("docker repo name")
		}

		return nil, e.NewInternalError(err.Error())
	}

	return dockerRepo, nil
}

func (dockerRepoService *DockerRepoServiceImpl) GetDockerRepoByID(ctx context.Context, dockerRepoID string) (*model.DockerRepo, e.ResponseError) {
	dockerRepo, err := dockerRepoService.dockerRepoMapper.FindByDockerRepoID(dockerRepoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("docker repo")
		}

		return nil, e.NewInternalError(err.Error())
	}

	return dockerRepo, nil
}

func (dockerRepoService *DockerRepoServiceImpl) GetDockerRepoByUserIDAndName(ctx context.Context, userID, dockerRepoName string) (*model.DockerRepo, e.ResponseError) {
	dockerRepo, err := dockerRepoService.dockerRepoMapper.FindByUserIDAndName(userID, dockerRepoName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("docker repo")
		}

		return nil, e.NewInternalError(err.Error())
	}

	return dockerRepo, nil
}

func (dockerRepoService *DockerRepoServiceImpl) ListDockerRepos(ctx context.Context, userID string, offset, limit int, reverse bool) (model.DockerRepos, e.ResponseError) {
	dockerRepos, err := dockerRepoService.dockerRepoMapper.FindPageByUserID(userID, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return dockerRepos, nil
}

func (dockerRepoService *DockerRepoServiceImpl) UpdateDockerRepoByID(ctx context.Context, dockerRepoID, address, userName, password string) e.ResponseError {
	err := dockerRepoService.dockerRepoMapper.UpdateByDockerRepoID(dockerRepoID, &model.DockerRepo{
		Address:  address,
		UserName: userName,
		Password: password,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("docker repo")
		}

		return e.NewInternalError(err.Error())
	}

	return nil
}

func (dockerRepoService *DockerRepoServiceImpl) UpdateDockerRepoByName(ctx context.Context, userID, dockerRepoName, address, userName, password string) e.ResponseError {
	err := dockerRepoService.dockerRepoMapper.UpdateByUserIDAndName(userID, dockerRepoName, &model.DockerRepo{
		Address:  address,
		UserName: userName,
		Password: password,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.NewNotFoundError("docker repo")
		}

		return e.NewInternalError(err.Error())
	}

	return nil
}

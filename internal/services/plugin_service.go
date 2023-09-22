package services

import (
	"context"
	"errors"
	"github.com/ProtobufMan/bufman/internal/core/docker"
	"github.com/ProtobufMan/bufman/internal/core/storage"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/gorm"
)

type PluginService interface {
	ListPlugins(ctx context.Context, offset int, limit int, reverse bool, includeDeprecated bool) (model.Plugins, e.ResponseError)
	CreatePlugin(ctx context.Context, plugin *model.Plugin, dockerRepoName string) (*model.Plugin, e.ResponseError)
	GetLatestPlugin(ctx context.Context, owner string, name string) (*model.Plugin, e.ResponseError)
	GetLatestPluginWithVersion(ctx context.Context, owner string, name string, version string) (*model.Plugin, e.ResponseError)
	GetLatestPluginWithVersionAndReversion(ctx context.Context, owner string, name string, version string, reversion uint32) (*model.Plugin, e.ResponseError)
}

type PluginServiceImpl struct {
	pluginMapper     mapper.PluginMapper
	dockerRepoMapper mapper.DockerRepoMapper
	storageHelper    storage.StorageHelper
}

func NewPluginService() PluginService {
	return &PluginServiceImpl{
		pluginMapper:     &mapper.PluginMapperImpl{},
		dockerRepoMapper: &mapper.DockerRepoMapperImpl{},
		storageHelper:    storage.NewStorageHelper(),
	}
}

func (pluginService *PluginServiceImpl) ListPlugins(ctx context.Context, offset int, limit int, reverse bool, includeDeprecated bool) (model.Plugins, e.ResponseError) {
	plugins, err := pluginService.pluginMapper.FindPage(offset, limit, reverse, includeDeprecated)
	if err != nil {
		return nil, e.NewInternalError("ListPlugins")
	}

	return plugins, nil
}

func (pluginService *PluginServiceImpl) CreatePlugin(ctx context.Context, pluginModel *model.Plugin, dockerRepoName string) (*model.Plugin, e.ResponseError) {
	// 查询docker repo
	dockerRepo, err := pluginService.dockerRepoMapper.FindByUserIDAndName(pluginModel.UserID, dockerRepoName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("docker repo name")
		}

		return nil, e.NewInternalError(err.Error())
	}

	// try pull
	d, err := docker.NewDockerClient(dockerRepo.Address, dockerRepo.UserName, dockerRepo.Password)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}
	defer d.Close()

	err = d.TryPullImage(ctx, pluginModel.ImageName, pluginModel.ImageDigest)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	// 记录在数据库中
	pluginModel.DockerRepoID = dockerRepo.DockerRepoID
	err = pluginService.pluginMapper.Create(pluginModel)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return pluginModel, e.NewAlreadyExistsError("plugin")
		}
		return nil, e.NewInternalError(err.Error())
	}

	return pluginModel, nil
}

func (pluginService *PluginServiceImpl) GetLatestPlugin(ctx context.Context, owner string, name string) (*model.Plugin, e.ResponseError) {
	plugin, err := pluginService.pluginMapper.FindLastByName(owner, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("plugin")
		}

		return nil, e.NewInternalError("GetLatestPlugin")
	}

	return plugin, nil
}

func (pluginService *PluginServiceImpl) GetLatestPluginWithVersion(ctx context.Context, owner string, name string, version string) (*model.Plugin, e.ResponseError) {
	plugin, err := pluginService.pluginMapper.FindLastByNameAndVersion(owner, name, version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("plugin")
		}

		return nil, e.NewInternalError("GetLatestPluginWithVersion")
	}

	return plugin, nil
}

func (pluginService *PluginServiceImpl) GetLatestPluginWithVersionAndReversion(ctx context.Context, owner string, name string, version string, reversion uint32) (*model.Plugin, e.ResponseError) {
	plugin, err := pluginService.pluginMapper.FindByNameAndVersionReversion(owner, name, version, reversion)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("plugin")
		}

		return nil, e.NewInternalError("GetLatestPluginWithVersionAndReversion")
	}

	return plugin, nil
}

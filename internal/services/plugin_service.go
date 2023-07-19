package services

import (
	"errors"
	"github.com/ProtobufMan/bufman-cli/private/gen/proto/connect/bufman/alpha/registry/v1alpha1/registryv1alpha1connect"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util/storage"
	"gorm.io/gorm"
)

type PluginService interface {
	ListPlugins(offset int, limit int, reverse bool, includeDeprecated bool) (model.Plugins, e.ResponseError)
	CreatePlugin(plugin *model.Plugin, binaryData []byte) (*model.Plugin, e.ResponseError)
	GetLatestPlugin(owner string, name string) (*model.Plugin, e.ResponseError)
	GetLatestPluginWithVersion(owner string, name string, version string) (*model.Plugin, e.ResponseError)
	GetLatestPluginWithVersionAndReversion(owner string, name string, version string, reversion uint32) (*model.Plugin, e.ResponseError)
}

type PluginServiceImpl struct {
	pluginMapper  mapper.PluginMapper
	storageHelper storage.StorageHelper
}

func NewPluginService() PluginService {
	return &PluginServiceImpl{
		pluginMapper:  &mapper.PluginMapperImpl{},
		storageHelper: storage.NewStorageHelper(),
	}
}

func (pluginService *PluginServiceImpl) ListPlugins(offset int, limit int, reverse bool, includeDeprecated bool) (model.Plugins, e.ResponseError) {
	plugins, err := pluginService.pluginMapper.FindPage(offset, limit, reverse, includeDeprecated)
	if err != nil {
		return nil, e.NewInternalError("ListPlugins")
	}

	return plugins, nil
}

func (pluginService *PluginServiceImpl) CreatePlugin(plugin *model.Plugin, binaryData []byte) (*model.Plugin, e.ResponseError) {
	// 将二进制保存起来
	binaryName, err := pluginService.storageHelper.StorePlugin(plugin.PluginName, plugin.Version, plugin.Reversion, binaryData)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	// 记录在数据库中
	plugin.BinaryName = binaryName
	err = pluginService.pluginMapper.Create(plugin)
	if err != nil {
		return nil, e.NewInternalError(registryv1alpha1connect.PluginCurationServiceCreateCuratedPluginProcedure)
	}

	return plugin, nil
}

func (pluginService *PluginServiceImpl) GetLatestPlugin(owner string, name string) (*model.Plugin, e.ResponseError) {
	plugin, err := pluginService.pluginMapper.FindLastByName(owner, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("plugin")
		}

		return nil, e.NewInternalError("GetLatestPlugin")
	}

	return plugin, nil
}

func (pluginService *PluginServiceImpl) GetLatestPluginWithVersion(owner string, name string, version string) (*model.Plugin, e.ResponseError) {
	plugin, err := pluginService.pluginMapper.FindLastByNameAndVersion(owner, name, version)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("plugin")
		}

		return nil, e.NewInternalError("GetLatestPluginWithVersion")
	}

	return plugin, nil
}

func (pluginService *PluginServiceImpl) GetLatestPluginWithVersionAndReversion(owner string, name string, version string, reversion uint32) (*model.Plugin, e.ResponseError) {
	plugin, err := pluginService.pluginMapper.FindByNameAndVersionReversion(owner, name, version, reversion)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewNotFoundError("plugin")
		}

		return nil, e.NewInternalError("GetLatestPluginWithVersionAndReversion")
	}

	return plugin, nil
}

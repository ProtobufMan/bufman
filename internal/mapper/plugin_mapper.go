package mapper

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
)

type PluginMapper interface {
	Create(plugin *model.Plugin) error
	FindByNameAndVersion(userName, pluginName, version string, reversion uint32) (*model.Plugin, error)
}

type PluginMapperImpl struct{}

func (p *PluginMapperImpl) Create(plugin *model.Plugin) error {
	return dal.Plugin.Create(plugin)
}

func (p *PluginMapperImpl) FindByNameAndVersion(userName, pluginName, version string, reversion uint32) (*model.Plugin, error) {
	return dal.Plugin.Where(dal.Plugin.UserName.Eq(userName), dal.Plugin.PluginName.Eq(pluginName), dal.Plugin.Version.Eq(version), dal.Plugin.Reversion.Eq(reversion)).First()
}

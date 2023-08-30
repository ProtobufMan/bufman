package mapper

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
)

type PluginMapper interface {
	Create(plugin *model.Plugin) error
	FindByNameAndVersionReversion(userName, pluginName, version string, reversion uint32) (*model.Plugin, error)
	FindLastByName(userName, pluginName string) (*model.Plugin, error)
	FindLastByNameAndVersion(userName, pluginName string, version string) (*model.Plugin, error)
	FindPage(offset int, limit int, reverse bool, includeDeprecated bool) ([]*model.Plugin, error)
	FindPageByQuery(query string, offset int, limit int, reverse bool) ([]*model.Plugin, error)
}

type PluginMapperImpl struct{}

func (p *PluginMapperImpl) Create(plugin *model.Plugin) error {
	return dal.Plugin.Create(plugin)
}

func (p *PluginMapperImpl) FindByNameAndVersionReversion(userName, pluginName, version string, reversion uint32) (*model.Plugin, error) {
	return dal.Plugin.Where(dal.Plugin.UserName.Eq(userName), dal.Plugin.PluginName.Eq(pluginName), dal.Plugin.Version.Eq(version), dal.Plugin.Reversion.Eq(reversion)).First()
}

func (p *PluginMapperImpl) FindLastByName(userName, pluginName string) (*model.Plugin, error) {
	return dal.Plugin.Where(dal.Plugin.UserName.Eq(userName), dal.Plugin.PluginName.Eq(pluginName)).Last()
}

func (p *PluginMapperImpl) FindLastByNameAndVersion(userName, pluginName string, version string) (*model.Plugin, error) {
	return dal.Plugin.Where(dal.Plugin.UserName.Eq(userName), dal.Plugin.PluginName.Eq(pluginName), dal.Plugin.Version.Eq(version)).Last()

}

func (p *PluginMapperImpl) FindPage(offset int, limit int, reverse bool, includeDeprecated bool) ([]*model.Plugin, error) {
	stmt := dal.Plugin.Offset(offset).Limit(limit)
	if reverse {
		// 反转
		stmt = stmt.Order(dal.Plugin.ID.Desc())
	}

	if !includeDeprecated {
		// 不包括已经被丢弃的
		stmt = stmt.Where(dal.Plugin.Deprecated.Not())
	}

	return stmt.Find()
}

func (p *PluginMapperImpl) FindPageByQuery(query string, offset int, limit int, reverse bool) ([]*model.Plugin, error) {
	stmt := dal.Plugin.Where(dal.Plugin.PluginName.Like("%" + query + "%")).Or(dal.Plugin.Description.Like("%" + query + "%")).Offset(offset).Limit(limit)
	if reverse {
		// 反转
		stmt = stmt.Order(dal.Plugin.ID.Desc())
	}

	return stmt.Find()
}

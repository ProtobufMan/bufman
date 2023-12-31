// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dal

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/ProtobufMan/bufman/internal/model"
)

func newPlugin(db *gorm.DB, opts ...gen.DOOption) plugin {
	_plugin := plugin{}

	_plugin.pluginDo.UseDB(db, opts...)
	_plugin.pluginDo.UseModel(&model.Plugin{})

	tableName := _plugin.pluginDo.TableName()
	_plugin.ALL = field.NewAsterisk(tableName)
	_plugin.ID = field.NewInt64(tableName, "id")
	_plugin.UserID = field.NewString(tableName, "user_id")
	_plugin.UserName = field.NewString(tableName, "user_name")
	_plugin.PluginID = field.NewString(tableName, "plugin_id")
	_plugin.PluginName = field.NewString(tableName, "plugin_name")
	_plugin.Version = field.NewString(tableName, "version")
	_plugin.Reversion = field.NewUint32(tableName, "reversion")
	_plugin.ImageName = field.NewString(tableName, "image_name")
	_plugin.ImageDigest = field.NewString(tableName, "image_digest")
	_plugin.DockerRepoID = field.NewString(tableName, "docker_repo_id")
	_plugin.Description = field.NewString(tableName, "description")
	_plugin.Visibility = field.NewUint8(tableName, "visibility")
	_plugin.Deprecated = field.NewBool(tableName, "deprecated")
	_plugin.DeprecationMsg = field.NewString(tableName, "deprecation_msg")
	_plugin.CreatedTime = field.NewTime(tableName, "created_time")
	_plugin.UpdateTime = field.NewTime(tableName, "update_time")

	_plugin.fillFieldMap()

	return _plugin
}

type plugin struct {
	pluginDo

	ALL            field.Asterisk
	ID             field.Int64
	UserID         field.String
	UserName       field.String
	PluginID       field.String
	PluginName     field.String
	Version        field.String
	Reversion      field.Uint32
	ImageName      field.String
	ImageDigest    field.String
	DockerRepoID   field.String
	Description    field.String
	Visibility     field.Uint8
	Deprecated     field.Bool
	DeprecationMsg field.String
	CreatedTime    field.Time
	UpdateTime     field.Time

	fieldMap map[string]field.Expr
}

func (p plugin) Table(newTableName string) *plugin {
	p.pluginDo.UseTable(newTableName)
	return p.updateTableName(newTableName)
}

func (p plugin) As(alias string) *plugin {
	p.pluginDo.DO = *(p.pluginDo.As(alias).(*gen.DO))
	return p.updateTableName(alias)
}

func (p *plugin) updateTableName(table string) *plugin {
	p.ALL = field.NewAsterisk(table)
	p.ID = field.NewInt64(table, "id")
	p.UserID = field.NewString(table, "user_id")
	p.UserName = field.NewString(table, "user_name")
	p.PluginID = field.NewString(table, "plugin_id")
	p.PluginName = field.NewString(table, "plugin_name")
	p.Version = field.NewString(table, "version")
	p.Reversion = field.NewUint32(table, "reversion")
	p.ImageName = field.NewString(table, "image_name")
	p.ImageDigest = field.NewString(table, "image_digest")
	p.DockerRepoID = field.NewString(table, "docker_repo_id")
	p.Description = field.NewString(table, "description")
	p.Visibility = field.NewUint8(table, "visibility")
	p.Deprecated = field.NewBool(table, "deprecated")
	p.DeprecationMsg = field.NewString(table, "deprecation_msg")
	p.CreatedTime = field.NewTime(table, "created_time")
	p.UpdateTime = field.NewTime(table, "update_time")

	p.fillFieldMap()

	return p
}

func (p *plugin) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := p.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (p *plugin) fillFieldMap() {
	p.fieldMap = make(map[string]field.Expr, 16)
	p.fieldMap["id"] = p.ID
	p.fieldMap["user_id"] = p.UserID
	p.fieldMap["user_name"] = p.UserName
	p.fieldMap["plugin_id"] = p.PluginID
	p.fieldMap["plugin_name"] = p.PluginName
	p.fieldMap["version"] = p.Version
	p.fieldMap["reversion"] = p.Reversion
	p.fieldMap["image_name"] = p.ImageName
	p.fieldMap["image_digest"] = p.ImageDigest
	p.fieldMap["docker_repo_id"] = p.DockerRepoID
	p.fieldMap["description"] = p.Description
	p.fieldMap["visibility"] = p.Visibility
	p.fieldMap["deprecated"] = p.Deprecated
	p.fieldMap["deprecation_msg"] = p.DeprecationMsg
	p.fieldMap["created_time"] = p.CreatedTime
	p.fieldMap["update_time"] = p.UpdateTime
}

func (p plugin) clone(db *gorm.DB) plugin {
	p.pluginDo.ReplaceConnPool(db.Statement.ConnPool)
	return p
}

func (p plugin) replaceDB(db *gorm.DB) plugin {
	p.pluginDo.ReplaceDB(db)
	return p
}

type pluginDo struct{ gen.DO }

type IPluginDo interface {
	gen.SubQuery
	Debug() IPluginDo
	WithContext(ctx context.Context) IPluginDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IPluginDo
	WriteDB() IPluginDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IPluginDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IPluginDo
	Not(conds ...gen.Condition) IPluginDo
	Or(conds ...gen.Condition) IPluginDo
	Select(conds ...field.Expr) IPluginDo
	Where(conds ...gen.Condition) IPluginDo
	Order(conds ...field.Expr) IPluginDo
	Distinct(cols ...field.Expr) IPluginDo
	Omit(cols ...field.Expr) IPluginDo
	Join(table schema.Tabler, on ...field.Expr) IPluginDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IPluginDo
	RightJoin(table schema.Tabler, on ...field.Expr) IPluginDo
	Group(cols ...field.Expr) IPluginDo
	Having(conds ...gen.Condition) IPluginDo
	Limit(limit int) IPluginDo
	Offset(offset int) IPluginDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IPluginDo
	Unscoped() IPluginDo
	Create(values ...*model.Plugin) error
	CreateInBatches(values []*model.Plugin, batchSize int) error
	Save(values ...*model.Plugin) error
	First() (*model.Plugin, error)
	Take() (*model.Plugin, error)
	Last() (*model.Plugin, error)
	Find() ([]*model.Plugin, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Plugin, err error)
	FindInBatches(result *[]*model.Plugin, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.Plugin) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IPluginDo
	Assign(attrs ...field.AssignExpr) IPluginDo
	Joins(fields ...field.RelationField) IPluginDo
	Preload(fields ...field.RelationField) IPluginDo
	FirstOrInit() (*model.Plugin, error)
	FirstOrCreate() (*model.Plugin, error)
	FindByPage(offset int, limit int) (result []*model.Plugin, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IPluginDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (p pluginDo) Debug() IPluginDo {
	return p.withDO(p.DO.Debug())
}

func (p pluginDo) WithContext(ctx context.Context) IPluginDo {
	return p.withDO(p.DO.WithContext(ctx))
}

func (p pluginDo) ReadDB() IPluginDo {
	return p.Clauses(dbresolver.Read)
}

func (p pluginDo) WriteDB() IPluginDo {
	return p.Clauses(dbresolver.Write)
}

func (p pluginDo) Session(config *gorm.Session) IPluginDo {
	return p.withDO(p.DO.Session(config))
}

func (p pluginDo) Clauses(conds ...clause.Expression) IPluginDo {
	return p.withDO(p.DO.Clauses(conds...))
}

func (p pluginDo) Returning(value interface{}, columns ...string) IPluginDo {
	return p.withDO(p.DO.Returning(value, columns...))
}

func (p pluginDo) Not(conds ...gen.Condition) IPluginDo {
	return p.withDO(p.DO.Not(conds...))
}

func (p pluginDo) Or(conds ...gen.Condition) IPluginDo {
	return p.withDO(p.DO.Or(conds...))
}

func (p pluginDo) Select(conds ...field.Expr) IPluginDo {
	return p.withDO(p.DO.Select(conds...))
}

func (p pluginDo) Where(conds ...gen.Condition) IPluginDo {
	return p.withDO(p.DO.Where(conds...))
}

func (p pluginDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IPluginDo {
	return p.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (p pluginDo) Order(conds ...field.Expr) IPluginDo {
	return p.withDO(p.DO.Order(conds...))
}

func (p pluginDo) Distinct(cols ...field.Expr) IPluginDo {
	return p.withDO(p.DO.Distinct(cols...))
}

func (p pluginDo) Omit(cols ...field.Expr) IPluginDo {
	return p.withDO(p.DO.Omit(cols...))
}

func (p pluginDo) Join(table schema.Tabler, on ...field.Expr) IPluginDo {
	return p.withDO(p.DO.Join(table, on...))
}

func (p pluginDo) LeftJoin(table schema.Tabler, on ...field.Expr) IPluginDo {
	return p.withDO(p.DO.LeftJoin(table, on...))
}

func (p pluginDo) RightJoin(table schema.Tabler, on ...field.Expr) IPluginDo {
	return p.withDO(p.DO.RightJoin(table, on...))
}

func (p pluginDo) Group(cols ...field.Expr) IPluginDo {
	return p.withDO(p.DO.Group(cols...))
}

func (p pluginDo) Having(conds ...gen.Condition) IPluginDo {
	return p.withDO(p.DO.Having(conds...))
}

func (p pluginDo) Limit(limit int) IPluginDo {
	return p.withDO(p.DO.Limit(limit))
}

func (p pluginDo) Offset(offset int) IPluginDo {
	return p.withDO(p.DO.Offset(offset))
}

func (p pluginDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IPluginDo {
	return p.withDO(p.DO.Scopes(funcs...))
}

func (p pluginDo) Unscoped() IPluginDo {
	return p.withDO(p.DO.Unscoped())
}

func (p pluginDo) Create(values ...*model.Plugin) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Create(values)
}

func (p pluginDo) CreateInBatches(values []*model.Plugin, batchSize int) error {
	return p.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (p pluginDo) Save(values ...*model.Plugin) error {
	if len(values) == 0 {
		return nil
	}
	return p.DO.Save(values)
}

func (p pluginDo) First() (*model.Plugin, error) {
	if result, err := p.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Plugin), nil
	}
}

func (p pluginDo) Take() (*model.Plugin, error) {
	if result, err := p.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Plugin), nil
	}
}

func (p pluginDo) Last() (*model.Plugin, error) {
	if result, err := p.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Plugin), nil
	}
}

func (p pluginDo) Find() ([]*model.Plugin, error) {
	result, err := p.DO.Find()
	return result.([]*model.Plugin), err
}

func (p pluginDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Plugin, err error) {
	buf := make([]*model.Plugin, 0, batchSize)
	err = p.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (p pluginDo) FindInBatches(result *[]*model.Plugin, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return p.DO.FindInBatches(result, batchSize, fc)
}

func (p pluginDo) Attrs(attrs ...field.AssignExpr) IPluginDo {
	return p.withDO(p.DO.Attrs(attrs...))
}

func (p pluginDo) Assign(attrs ...field.AssignExpr) IPluginDo {
	return p.withDO(p.DO.Assign(attrs...))
}

func (p pluginDo) Joins(fields ...field.RelationField) IPluginDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Joins(_f))
	}
	return &p
}

func (p pluginDo) Preload(fields ...field.RelationField) IPluginDo {
	for _, _f := range fields {
		p = *p.withDO(p.DO.Preload(_f))
	}
	return &p
}

func (p pluginDo) FirstOrInit() (*model.Plugin, error) {
	if result, err := p.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Plugin), nil
	}
}

func (p pluginDo) FirstOrCreate() (*model.Plugin, error) {
	if result, err := p.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Plugin), nil
	}
}

func (p pluginDo) FindByPage(offset int, limit int) (result []*model.Plugin, count int64, err error) {
	result, err = p.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = p.Offset(-1).Limit(-1).Count()
	return
}

func (p pluginDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = p.Count()
	if err != nil {
		return
	}

	err = p.Offset(offset).Limit(limit).Scan(result)
	return
}

func (p pluginDo) Scan(result interface{}) (err error) {
	return p.DO.Scan(result)
}

func (p pluginDo) Delete(models ...*model.Plugin) (result gen.ResultInfo, err error) {
	return p.DO.Delete(models)
}

func (p *pluginDo) withDO(do gen.Dao) *pluginDo {
	p.DO = *do.(*gen.DO)
	return p
}

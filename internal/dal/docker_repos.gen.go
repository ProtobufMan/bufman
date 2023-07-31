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

func newDockerRepo(db *gorm.DB, opts ...gen.DOOption) dockerRepo {
	_dockerRepo := dockerRepo{}

	_dockerRepo.dockerRepoDo.UseDB(db, opts...)
	_dockerRepo.dockerRepoDo.UseModel(&model.DockerRepo{})

	tableName := _dockerRepo.dockerRepoDo.TableName()
	_dockerRepo.ALL = field.NewAsterisk(tableName)
	_dockerRepo.ID = field.NewInt64(tableName, "id")
	_dockerRepo.UserID = field.NewString(tableName, "user_id")
	_dockerRepo.DockerRepoID = field.NewString(tableName, "docker_repo_id")
	_dockerRepo.DockerRepoName = field.NewString(tableName, "docker_repo_name")
	_dockerRepo.Address = field.NewString(tableName, "address")
	_dockerRepo.UserName = field.NewString(tableName, "user_name")
	_dockerRepo.Password = field.NewString(tableName, "password")
	_dockerRepo.CreatedTime = field.NewTime(tableName, "created_time")
	_dockerRepo.UpdateTime = field.NewTime(tableName, "update_time")
	_dockerRepo.Note = field.NewString(tableName, "note")

	_dockerRepo.fillFieldMap()

	return _dockerRepo
}

type dockerRepo struct {
	dockerRepoDo

	ALL            field.Asterisk
	ID             field.Int64
	UserID         field.String
	DockerRepoID   field.String
	DockerRepoName field.String
	Address        field.String
	UserName       field.String
	Password       field.String
	CreatedTime    field.Time
	UpdateTime     field.Time
	Note           field.String

	fieldMap map[string]field.Expr
}

func (d dockerRepo) Table(newTableName string) *dockerRepo {
	d.dockerRepoDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d dockerRepo) As(alias string) *dockerRepo {
	d.dockerRepoDo.DO = *(d.dockerRepoDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *dockerRepo) updateTableName(table string) *dockerRepo {
	d.ALL = field.NewAsterisk(table)
	d.ID = field.NewInt64(table, "id")
	d.UserID = field.NewString(table, "user_id")
	d.DockerRepoID = field.NewString(table, "docker_repo_id")
	d.DockerRepoName = field.NewString(table, "docker_repo_name")
	d.Address = field.NewString(table, "address")
	d.UserName = field.NewString(table, "user_name")
	d.Password = field.NewString(table, "password")
	d.CreatedTime = field.NewTime(table, "created_time")
	d.UpdateTime = field.NewTime(table, "update_time")
	d.Note = field.NewString(table, "note")

	d.fillFieldMap()

	return d
}

func (d *dockerRepo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *dockerRepo) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 10)
	d.fieldMap["id"] = d.ID
	d.fieldMap["user_id"] = d.UserID
	d.fieldMap["docker_repo_id"] = d.DockerRepoID
	d.fieldMap["docker_repo_name"] = d.DockerRepoName
	d.fieldMap["address"] = d.Address
	d.fieldMap["user_name"] = d.UserName
	d.fieldMap["password"] = d.Password
	d.fieldMap["created_time"] = d.CreatedTime
	d.fieldMap["update_time"] = d.UpdateTime
	d.fieldMap["note"] = d.Note
}

func (d dockerRepo) clone(db *gorm.DB) dockerRepo {
	d.dockerRepoDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d dockerRepo) replaceDB(db *gorm.DB) dockerRepo {
	d.dockerRepoDo.ReplaceDB(db)
	return d
}

type dockerRepoDo struct{ gen.DO }

type IDockerRepoDo interface {
	gen.SubQuery
	Debug() IDockerRepoDo
	WithContext(ctx context.Context) IDockerRepoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IDockerRepoDo
	WriteDB() IDockerRepoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IDockerRepoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IDockerRepoDo
	Not(conds ...gen.Condition) IDockerRepoDo
	Or(conds ...gen.Condition) IDockerRepoDo
	Select(conds ...field.Expr) IDockerRepoDo
	Where(conds ...gen.Condition) IDockerRepoDo
	Order(conds ...field.Expr) IDockerRepoDo
	Distinct(cols ...field.Expr) IDockerRepoDo
	Omit(cols ...field.Expr) IDockerRepoDo
	Join(table schema.Tabler, on ...field.Expr) IDockerRepoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IDockerRepoDo
	RightJoin(table schema.Tabler, on ...field.Expr) IDockerRepoDo
	Group(cols ...field.Expr) IDockerRepoDo
	Having(conds ...gen.Condition) IDockerRepoDo
	Limit(limit int) IDockerRepoDo
	Offset(offset int) IDockerRepoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IDockerRepoDo
	Unscoped() IDockerRepoDo
	Create(values ...*model.DockerRepo) error
	CreateInBatches(values []*model.DockerRepo, batchSize int) error
	Save(values ...*model.DockerRepo) error
	First() (*model.DockerRepo, error)
	Take() (*model.DockerRepo, error)
	Last() (*model.DockerRepo, error)
	Find() ([]*model.DockerRepo, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.DockerRepo, err error)
	FindInBatches(result *[]*model.DockerRepo, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.DockerRepo) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IDockerRepoDo
	Assign(attrs ...field.AssignExpr) IDockerRepoDo
	Joins(fields ...field.RelationField) IDockerRepoDo
	Preload(fields ...field.RelationField) IDockerRepoDo
	FirstOrInit() (*model.DockerRepo, error)
	FirstOrCreate() (*model.DockerRepo, error)
	FindByPage(offset int, limit int) (result []*model.DockerRepo, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IDockerRepoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (d dockerRepoDo) Debug() IDockerRepoDo {
	return d.withDO(d.DO.Debug())
}

func (d dockerRepoDo) WithContext(ctx context.Context) IDockerRepoDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d dockerRepoDo) ReadDB() IDockerRepoDo {
	return d.Clauses(dbresolver.Read)
}

func (d dockerRepoDo) WriteDB() IDockerRepoDo {
	return d.Clauses(dbresolver.Write)
}

func (d dockerRepoDo) Session(config *gorm.Session) IDockerRepoDo {
	return d.withDO(d.DO.Session(config))
}

func (d dockerRepoDo) Clauses(conds ...clause.Expression) IDockerRepoDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d dockerRepoDo) Returning(value interface{}, columns ...string) IDockerRepoDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d dockerRepoDo) Not(conds ...gen.Condition) IDockerRepoDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d dockerRepoDo) Or(conds ...gen.Condition) IDockerRepoDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d dockerRepoDo) Select(conds ...field.Expr) IDockerRepoDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d dockerRepoDo) Where(conds ...gen.Condition) IDockerRepoDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d dockerRepoDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IDockerRepoDo {
	return d.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (d dockerRepoDo) Order(conds ...field.Expr) IDockerRepoDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d dockerRepoDo) Distinct(cols ...field.Expr) IDockerRepoDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d dockerRepoDo) Omit(cols ...field.Expr) IDockerRepoDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d dockerRepoDo) Join(table schema.Tabler, on ...field.Expr) IDockerRepoDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d dockerRepoDo) LeftJoin(table schema.Tabler, on ...field.Expr) IDockerRepoDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d dockerRepoDo) RightJoin(table schema.Tabler, on ...field.Expr) IDockerRepoDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d dockerRepoDo) Group(cols ...field.Expr) IDockerRepoDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d dockerRepoDo) Having(conds ...gen.Condition) IDockerRepoDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d dockerRepoDo) Limit(limit int) IDockerRepoDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d dockerRepoDo) Offset(offset int) IDockerRepoDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d dockerRepoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IDockerRepoDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d dockerRepoDo) Unscoped() IDockerRepoDo {
	return d.withDO(d.DO.Unscoped())
}

func (d dockerRepoDo) Create(values ...*model.DockerRepo) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d dockerRepoDo) CreateInBatches(values []*model.DockerRepo, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d dockerRepoDo) Save(values ...*model.DockerRepo) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d dockerRepoDo) First() (*model.DockerRepo, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.DockerRepo), nil
	}
}

func (d dockerRepoDo) Take() (*model.DockerRepo, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.DockerRepo), nil
	}
}

func (d dockerRepoDo) Last() (*model.DockerRepo, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.DockerRepo), nil
	}
}

func (d dockerRepoDo) Find() ([]*model.DockerRepo, error) {
	result, err := d.DO.Find()
	return result.([]*model.DockerRepo), err
}

func (d dockerRepoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.DockerRepo, err error) {
	buf := make([]*model.DockerRepo, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d dockerRepoDo) FindInBatches(result *[]*model.DockerRepo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d dockerRepoDo) Attrs(attrs ...field.AssignExpr) IDockerRepoDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d dockerRepoDo) Assign(attrs ...field.AssignExpr) IDockerRepoDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d dockerRepoDo) Joins(fields ...field.RelationField) IDockerRepoDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d dockerRepoDo) Preload(fields ...field.RelationField) IDockerRepoDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d dockerRepoDo) FirstOrInit() (*model.DockerRepo, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.DockerRepo), nil
	}
}

func (d dockerRepoDo) FirstOrCreate() (*model.DockerRepo, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.DockerRepo), nil
	}
}

func (d dockerRepoDo) FindByPage(offset int, limit int) (result []*model.DockerRepo, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d dockerRepoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d dockerRepoDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d dockerRepoDo) Delete(models ...*model.DockerRepo) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *dockerRepoDo) withDO(do gen.Dao) *dockerRepoDo {
	d.DO = *do.(*gen.DO)
	return d
}

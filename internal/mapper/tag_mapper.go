package mapper

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
)

type TagMapper interface {
	Create(tag *model.Tag) error
	GetCountsByRepositoryID(repositoryID string) (int64, error)
	FindPageByRepositoryID(repositoryID string, offset, limit int, reverse bool) (model.Tags, error)
}

type TagMapperImpl struct{}

func (t *TagMapperImpl) Create(tag *model.Tag) error {
	return dal.Tag.Create(tag)
}

func (t *TagMapperImpl) GetCountsByRepositoryID(repositoryID string) (int64, error) {
	return dal.Tag.Where(dal.Tag.RepositoryID.Eq(repositoryID)).Count()
}

func (t *TagMapperImpl) FindPageByRepositoryID(repositoryID string, offset, limit int, reverse bool) (model.Tags, error) {
	stmt := dal.Tag.Where(dal.Tag.RepositoryID.Eq(repositoryID)).Offset(offset).Limit(limit)
	if reverse {
		stmt = stmt.Order(dal.Tag.ID.Desc())
	}

	return stmt.Find()
}

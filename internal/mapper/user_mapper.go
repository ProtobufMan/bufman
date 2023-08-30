package mapper

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
)

type UserMapper interface {
	Create(user *model.User) error
	FindByUserID(userID string) (*model.User, error)
	FindByUserName(userName string) (*model.User, error)
	FindPage(offset int, limit int, reverse bool) (model.Users, error)
	FindPageByQuery(query string, offset int, limit int, reverse bool) (model.Users, error)
}

type UserMapperImpl struct{}

func (u *UserMapperImpl) Create(user *model.User) error {
	return dal.User.Create(user)
}

func (u *UserMapperImpl) FindByUserID(userID string) (*model.User, error) {
	return dal.User.Where(dal.User.UserID.Eq(userID)).First()
}

func (u *UserMapperImpl) FindByUserName(userName string) (*model.User, error) {
	return dal.User.Where(dal.User.UserName.Eq(userName)).First()
}

func (u *UserMapperImpl) FindPage(offset int, limit int, reverse bool) (model.Users, error) {
	stmt := dal.User
	if reverse {
		stmt.Order(dal.User.ID.Desc())
	}

	users, _, err := stmt.FindByPage(offset, limit)
	return users, err
}

func (u *UserMapperImpl) FindPageByQuery(query string, offset int, limit int, reverse bool) (model.Users, error) {
	stmt := dal.User.Where(dal.User.UserName.Like("%" + query + "%"))
	if reverse {
		stmt.Order(dal.User.ID.Desc())
	}

	users, _, err := stmt.FindByPage(offset, limit)
	return users, err
}

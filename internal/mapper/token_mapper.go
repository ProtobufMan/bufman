package mapper

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
	"time"
)

type TokenMapper interface {
	Create(token *model.Token) error
	FindAvailableByTokenID(tokenID string) (*model.Token, error)
	FindAvailableByTokenName(tokenName string) (*model.Token, error)
	FindAvailablePageByUserID(userID string, offset int, limit int, reverse bool) (model.Tokens, error)
	DeleteByTokenID(tokenID string) error
}

type TokenMapperImpl struct{}

func (t *TokenMapperImpl) Create(token *model.Token) error {
	return dal.Token.Create(token)
}

func (t *TokenMapperImpl) FindAvailableByTokenID(tokenID string) (*model.Token, error) {
	return dal.Token.Where(dal.Token.TokenID.Eq(tokenID), dal.Token.ExpireTime.Gt(time.Now())).First()
}

func (t *TokenMapperImpl) FindAvailableByTokenName(tokenName string) (*model.Token, error) {
	return dal.Token.Where(dal.Token.TokenName.Eq(tokenName), dal.Token.ExpireTime.Gt(time.Now())).First()
}

func (t *TokenMapperImpl) FindAvailablePageByUserID(userID string, offset int, limit int, reverse bool) (model.Tokens, error) {
	stmt := dal.Token.Where(dal.Token.UserID.Eq(userID), dal.Token.ExpireTime.Gt(time.Now()))
	if reverse {
		stmt.Order(dal.Token.ID.Desc())
	}

	tokens, _, err := stmt.FindByPage(offset, limit)

	return tokens, err
}

func (t *TokenMapperImpl) DeleteByTokenID(tokenID string) error {
	token := &model.Token{}
	_, err := dal.Token.Where(dal.Token.TokenID.Eq(tokenID), dal.Token.ExpireTime.Gt(time.Now())).Delete(token)
	return err
}

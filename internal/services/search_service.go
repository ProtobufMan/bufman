package services

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
)

type SearchService interface {
	SearchUser(ctx context.Context, query string, offset, limit int, reverse bool) (model.Users, e.ResponseError)
}

type SearchServiceImpl struct {
	userMapper mapper.UserMapper
}

func (searchService *SearchServiceImpl) SearchUser(ctx context.Context, query string, offset, limit int, reverse bool) (model.Users, e.ResponseError) {
	users, err := searchService.userMapper.FindPageByQuery(query, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return users, nil
}

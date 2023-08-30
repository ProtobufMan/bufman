package services

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
)

type SearchService interface {
	SearchUser(ctx context.Context, query string, offset, limit int, reverse bool) (model.Users, e.ResponseError)
	SearchRepository(ctx context.Context, query string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	SearchCurationPlugin(ctx context.Context, query string, offset, limit int, reverse bool) (model.Plugins, e.ResponseError)
	SearchTag(ctx context.Context, repositoryID, query string, offset, limit int, reverse bool) (model.Tags, e.ResponseError)
	SearchDraft(ctx context.Context, repositoryID, query string, offset, limit int, reverse bool) (model.Commits, e.ResponseError)
}

type SearchServiceImpl struct {
	userMapper       mapper.UserMapper
	repositoryMapper mapper.RepositoryMapper
	commitMapper     mapper.CommitMapper
	tagMapper        mapper.TagMapper
	pluginMapper     mapper.PluginMapper
}

func (searchService *SearchServiceImpl) SearchUser(ctx context.Context, query string, offset, limit int, reverse bool) (model.Users, e.ResponseError) {
	users, err := searchService.userMapper.FindPageByQuery(query, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return users, nil
}

func (searchService *SearchServiceImpl) SearchRepository(ctx context.Context, query string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError) {
	repositories, err := searchService.repositoryMapper.FindPageByQuery(query, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return repositories, nil
}

func (searchService *SearchServiceImpl) SearchCurationPlugin(ctx context.Context, query string, offset, limit int, reverse bool) (model.Plugins, e.ResponseError) {
	plugins, err := searchService.pluginMapper.FindPageByQuery(query, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return plugins, nil
}

func (searchService *SearchServiceImpl) SearchTag(ctx context.Context, repositoryID, query string, offset, limit int, reverse bool) (model.Tags, e.ResponseError) {
	tags, err := searchService.tagMapper.FindPageByRepositoryIDAndQuery(repositoryID, query, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return tags, nil
}

func (searchService *SearchServiceImpl) SearchDraft(ctx context.Context, repositoryID, query string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	commits, err := searchService.commitMapper.FindDraftPageByRepositoryIDAndQuery(repositoryID, query, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return commits, nil
}

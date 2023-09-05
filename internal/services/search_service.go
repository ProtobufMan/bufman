package services

import (
	"context"
	"encoding/json"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/es"
	"github.com/ProtobufMan/bufman/internal/core/lru"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
)

type SearchService interface {
	SearchUser(ctx context.Context, query string, offset, limit int, reverse bool) (model.Users, e.ResponseError)
	SearchRepository(ctx context.Context, query string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	SearchLastCommitByContent(ctx context.Context, query string, offset, limit int, reverse bool) (model.Commits, e.ResponseError)
	SearchCurationPlugin(ctx context.Context, query string, offset, limit int, reverse bool) (model.Plugins, e.ResponseError)
	SearchTag(ctx context.Context, repositoryID, query string, offset, limit int, reverse bool) (model.Tags, e.ResponseError)
	SearchDraft(ctx context.Context, repositoryID, query string, offset, limit int, reverse bool) (model.Commits, e.ResponseError)
}

func NewSearchService() SearchService {
	return &SearchServiceImpl{
		userMapper:       &mapper.UserMapperImpl{},
		repositoryMapper: &mapper.RepositoryMapperImpl{},
		commitMapper:     &mapper.CommitMapperImpl{},
		tagMapper:        &mapper.TagMapperImpl{},
		pluginMapper:     &mapper.PluginMapperImpl{},
	}
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

func (searchService *SearchServiceImpl) SearchLastCommitByContent(ctx context.Context, query string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	// 连接elastic search
	esClient, err := es.NewEsClient()
	if err != nil || esClient == nil {
		return nil, e.NewInternalError(err.Error())
	}
	defer esClient.Close()

	// 在 ElasticSearch 中查询数据
	results, err := esClient.Query(ctx, constant.ESFileBlobIndex, query, offset, limit)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	lruQueue := lru.NewLru(len(results))
	for i := 0; i < len(results); i++ {
		data := results[i]
		fileBlob := &model.FileBlob{}
		err = json.Unmarshal(data, &fileBlob)
		if err != nil {
			continue
		}

		identity := fileBlob.RepositoryID
		if v, ok := lruQueue.Get(identity); !ok || (ok && v.(*model.FileBlob).CreatedTime.Before(fileBlob.CreatedTime)) {
			// 同一个repo下，只记录最晚的匹配commit
			_ = lruQueue.Add(identity, fileBlob)
		}
	}

	// 转为commit
	// 在LRU队列上，越靠近后方，在es中的查询结果位置越靠前，所以倒序遍历
	commits := make([]*model.Commit, 0, lruQueue.Len())
	_ = lruQueue.RangeValue(true, func(key, value interface{}) error {
		fileBlob, ok := value.(*model.FileBlob)
		if !ok {
			return nil
		}

		commit := &model.Commit{
			CommitID:       fileBlob.CommitID,
			CommitName:     fileBlob.CommitName,
			UserID:         fileBlob.UserID,
			UserName:       fileBlob.UserName,
			RepositoryID:   fileBlob.RepositoryID,
			RepositoryName: fileBlob.RepositoryName,
		}

		commits = append(commits, commit)

		return nil
	})

	return commits, nil
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

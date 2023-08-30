package services

import (
	"context"
	"encoding/json"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/e"
	"github.com/ProtobufMan/bufman/internal/mapper"
	"github.com/ProtobufMan/bufman/internal/model"
	"github.com/ProtobufMan/bufman/internal/util/es"
	"sort"
)

type SearchService interface {
	SearchUser(ctx context.Context, query string, offset, limit int, reverse bool) (model.Users, e.ResponseError)
	SearchRepository(ctx context.Context, query string, offset, limit int, reverse bool) (model.Repositories, e.ResponseError)
	SearchLastCommitByContent(ctx context.Context, query string, offset, limit int, reverse bool) (model.Commits, e.ResponseError)
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
	esClient         es.Client
}

func (searchService *SearchServiceImpl) SearchUser(ctx context.Context, query string, offset, limit int, reverse bool) (model.Users, e.ResponseError) {
	users, err := searchService.userMapper.FindPageByQuery(query, offset, limit, reverse)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	return users, nil
}

func (searchService *SearchServiceImpl) SearchLastCommitByContent(ctx context.Context, query string, offset, limit int, reverse bool) (model.Commits, e.ResponseError) {
	if searchService.esClient == nil {
		return nil, e.NewInternalError("not implement search commit by content")
	}

	// 在 ElasticSearch 中查询数据
	results, err := searchService.esClient.Query(ctx, constant.ESFileBlobIndex, query, offset, limit)
	if err != nil {
		return nil, e.NewInternalError(err.Error())
	}

	fileBlobSet := make(map[string]*model.FileBlob)
	for i := 0; i < len(results); i++ {
		data := results[i]
		fileBlob := &model.FileBlob{}
		err = json.Unmarshal(data, &fileBlob)
		if err != nil {
			continue
		}

		identity := fileBlob.RepositoryID
		if b, ok := fileBlobSet[identity]; !ok || (ok && b.CreatedTime.Before(fileBlob.CreatedTime)) {
			// 同一个repo下，只记录最晚的匹配commit
			fileBlobSet[identity] = b
		}
	}

	// 查询顺序调整
	keys := make([]string, 0, len(fileBlobSet))
	for key := range fileBlobSet {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 转为对应的commit
	commits := make([]*model.Commit, 0, len(fileBlobSet))
	for i := 0; i < len(keys); i++ {
		fileBlob, ok := fileBlobSet[keys[i]]
		if !ok {
			continue
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
	}

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

package es

import (
	"context"
	"github.com/ProtobufMan/bufman/internal/config"
	"github.com/olivere/elastic/v7"
)

type Client interface {
	Create(ctx context.Context, index string, id string, body interface{}) error
	Delete(ctx context.Context, index string, id string) error
	Find(ctx context.Context, index string, id string) ([]byte, error)
	Query(ctx context.Context, index string, query string, offset, limit int) ([][]byte, error)
	Close() error
}

func NewEsClient() (Client, error) {
	v, err := config.EsCliPool.Get()
	if err != nil {
		return nil, err
	}
	c := v.(*elastic.Client)

	return &clientImpl{
		client: c,
	}, nil
}

type clientImpl struct {
	client *elastic.Client
}

func (c *clientImpl) Create(ctx context.Context, index string, id string, body interface{}) error {
	_, err := c.client.Index().
		Index(index).
		Id(id).
		BodyJson(body).
		Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *clientImpl) Delete(ctx context.Context, index string, id string) error {
	_, err := c.client.Delete().Index(index).Id(id).Do(ctx)
	return err
}

func (c *clientImpl) Find(ctx context.Context, index string, id string) ([]byte, error) {
	res, err := c.client.Get().Index(index).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}

	return res.Source.MarshalJSON()
}

func (c *clientImpl) Query(ctx context.Context, index string, query string, offset, limit int) ([][]byte, error) {
	q := elastic.NewMultiMatchQuery(query, "RepositoryName", "Content")
	res, err := c.client.Search(index).Query(q).Size(limit).From(offset).Do(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([][]byte, 0, res.Hits.TotalHits.Value)
	for _, hit := range res.Hits.Hits {
		data, err := hit.Source.MarshalJSON()
		if err != nil {
			continue
		}

		ret = append(ret, data)
	}

	return ret, nil
}

func (c *clientImpl) Close() error {
	if c.client != nil {
		return config.EsCliPool.Put(c.client)
	}

	return nil
}

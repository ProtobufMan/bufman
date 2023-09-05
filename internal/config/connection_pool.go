package config

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/olivere/elastic/v7"
	"github.com/silenceper/pool"
)

func NewDockerCliPool() (pool.Pool, error) {
	c := &pool.Config{
		InitialCap: 1,
		MaxCap:     Properties.Docker.MaxOpenConnections,
		MaxIdle:    Properties.Docker.MaxIdleConnections,
		Factory: func() (interface{}, error) {
			dockerCli, err := NewDockerClient()
			if err != nil {
				return nil, err
			}

			return dockerCli, nil
		},
		Close: func(v interface{}) error {
			return v.(*client.Client).Close()
		},
		Ping: func(v interface{}) error {
			dockerCli := v.(*client.Client)
			_, err := dockerCli.Ping(context.Background())

			return err
		},
		IdleTimeout: Properties.Docker.MaxIdleTime,
	}

	return pool.NewChannelPool(c)
}

func NewElasticSearchCliPool() (pool.Pool, error) {
	c := &pool.Config{
		InitialCap: 1,
		MaxCap:     Properties.ElasticSearch.MaxOpenConnections,
		MaxIdle:    Properties.ElasticSearch.MaxIdleConnections,
		Factory: func() (interface{}, error) {
			esClient, err := NewEsClient()
			if err != nil {
				return nil, err
			}

			return esClient, nil
		},
		Close: func(v interface{}) error {
			return nil
		},
		Ping: func(v interface{}) error {
			elasticCli := v.(*elastic.Client)
			for i := 0; i < len(Properties.ElasticSearch.Urls); i++ {
				url := Properties.ElasticSearch.Urls[i]
				_, _, err := elasticCli.Ping(url).Do(context.Background())
				if err != nil {
					return err
				}
			}

			return nil
		},
		IdleTimeout: Properties.ElasticSearch.MaxIdleTime,
	}

	return pool.NewChannelPool(c)
}

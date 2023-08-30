package config

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/docker/docker/client"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"os"
	"time"
)

/*
TODO ！！！这是临时的，之后改为读取配置文件
*/
const (
	mysqlDSNKey   = "IDL_MGR_MYSQL_DSN"
	serverHostKey = "IDL_MGR_SERVER_HOST"
)

type Config struct {
	BufMan        BufMan
	Docker        Docker
	ElasticSearch ElasticSearch
}

type BufMan struct {
	ServerHost string
	MysqlDsn   string

	PageTokenExpireTime time.Duration
	PageTokenSecret     string
}

type Docker struct {
	Host       string
	CACertPath string
	CertPath   string
	KeyPath    string
}

type ElasticSearch struct {
	Urls     []string
	Username string
	Password string
}

var (
	DataBase   *gorm.DB
	Properties *Config
)

func LoadConfig() {
	Properties = &Config{
		BufMan: BufMan{
			ServerHost: "bufman.io",

			PageTokenExpireTime: time.Minute * 10, // 默认过期时间为10分钟
			PageTokenSecret:     "123456",
		},
		Docker: Docker{
			Host:       client.DefaultDockerHost,
			CACertPath: "",
			CertPath:   "",
			KeyPath:    "",
		},
		ElasticSearch: ElasticSearch{
			Urls:     []string{elastic.DefaultURL},
			Username: "",
			Password: "",
		},
	}

	if mysqlDSNENV := os.Getenv(mysqlDSNKey); mysqlDSNENV != "" {
		Properties.BufMan.MysqlDsn = mysqlDSNENV
	}
	if serverHostENV := os.Getenv(serverHostKey); serverHostENV != "" {
		Properties.BufMan.ServerHost = serverHostENV
	}

	if err := os.MkdirAll(constant.FileSavaDir, 0666); err != nil {
		panic(err)
	}

	// test docker config
	cli, err := NewDockerClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
}

func NewDockerClient() (*client.Client, error) {
	options := make([]client.Opt, 0, 4)
	options = append(options, client.WithAPIVersionNegotiation())
	options = append(options, client.WithHost(Properties.Docker.Host))
	if Properties.Docker.CACertPath != "" {
		options = append(options, client.WithTLSClientConfig(Properties.Docker.CACertPath, Properties.Docker.CertPath, Properties.Docker.KeyPath))
	}
	cli, err := client.NewClientWithOpts(options...)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func NewEsClient() (*elastic.Client, error) {
	if len(Properties.ElasticSearch.Urls) == 0 {
		return nil, errors.New("no elastic search config")
	}

	c, err := elastic.NewClient(elastic.SetURL(Properties.ElasticSearch.Urls...), elastic.SetBasicAuth(Properties.ElasticSearch.Username, Properties.ElasticSearch.Password))
	if err != nil {
		return nil, err
	}

	return c, nil
}

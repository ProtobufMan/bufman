package config

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/ProtobufMan/bufman/internal/core/logger"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/silenceper/pool"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"time"
)

const (
	mysqlDSNKey        = "BUFMAN_MYSQL_DSN"
	pageTokenSecretKey = "BUFMAN_PAGE_TOKEN_SECRET"
	esUsernameKey      = "BUFMAN_ES_USERNAME"
	esPasswordKey      = "BUFMAN_ES_PASSWORD"
)

const (
	configFileName = "config.yaml"
	configFileType = "yaml"
)

const (
	storageFSMode = "fs"
	storageESMode = "elasticsearch"
)

type Config struct {
	BufMan        BufMan        `mapstructure:"bufman"`
	MySQL         MySQL         `mapstructure:"mysql"`
	Docker        Docker        `mapstructure:"docker"`
	ElasticSearch ElasticSearch `mapstructure:"elastic_search"`
}

type BufMan struct {
	Mode       string `mapstructure:"mode"`
	ServerHost string `mapstructure:"server_host"`
	Port       int    `mapstructure:"port"`

	PageTokenExpireTime time.Duration `mapstructure:"page_token_expire_time"`
	PageTokenSecret     string        `mapstructure:"page_token_secret"`

	StorageMode  string `mapstructure:"storage_mode"`
	UseFSStorage bool   `mapstructure:"-"`
}

type MySQL struct {
	MysqlDsn           string        `mapstructure:"mysql_dsn"`
	MaxOpenConnections int           `mapstructure:"max_open_connections"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections"`
	MaxLifeTime        time.Duration `mapstructure:"max_life_time"`
	MaxIdleTime        time.Duration `mapstructure:"max_idle_time"`
}

type Docker struct {
	Host               string        `mapstructure:"host"`
	CACertPath         string        `mapstructure:"ca_cert_path"`
	CertPath           string        `mapstructure:"cert_path"`
	KeyPath            string        `mapstructure:"key_path"`
	MaxOpenConnections int           `mapstructure:"max_open_connections"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections"`
	MaxIdleTime        time.Duration `mapstructure:"max_idle_time"`
}

type ElasticSearch struct {
	Urls               []string      `mapstructure:"urls"`
	Username           string        `mapstructure:"username"`
	Password           string        `mapstructure:"password"`
	MaxOpenConnections int           `mapstructure:"max_open_connections"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections"`
	MaxIdleTime        time.Duration `mapstructure:"max_idle_time"`
}

var (
	DataBase   *gorm.DB
	Properties = &Config{}

	DockerCliPool pool.Pool
	EsCliPool     pool.Pool
)

func LoadConfig() {
	// 默认值
	Properties = &Config{
		BufMan: BufMan{
			Mode:                gin.DebugMode,
			ServerHost:          "bufman.io",
			Port:                8080,
			PageTokenExpireTime: time.Minute * 10, // 默认过期时间为10分钟
			PageTokenSecret:     "123456",
			StorageMode:         storageFSMode,
		},
		MySQL: MySQL{
			MysqlDsn:           "root:12345678@tcp(127.0.0.1:3306)/bufman?charset=utf8mb4&parseTime=True&loc=Local",
			MaxOpenConnections: 10,
			MaxIdleConnections: 10,
			MaxLifeTime:        15 * time.Minute,
			MaxIdleTime:        15 * time.Minute,
		},
		Docker: Docker{
			Host:               client.DefaultDockerHost,
			CACertPath:         "",
			CertPath:           "",
			KeyPath:            "",
			MaxOpenConnections: 10,
			MaxIdleConnections: 10,
			MaxIdleTime:        15 * time.Minute,
		},
		ElasticSearch: ElasticSearch{
			Urls:               []string{elastic.DefaultURL},
			Username:           "",
			Password:           "",
			MaxOpenConnections: 10,
			MaxIdleConnections: 10,
			MaxIdleTime:        15 * time.Minute,
		},
	}

	// 从配置文件中读取
	loadFromFile()

	// 从环境变量中读取
	loadFromENV()

	if Properties.BufMan.StorageMode != storageESMode {
		Properties.BufMan.UseFSStorage = true
	}

	// gin、logger设置level
	gin.SetMode(Properties.BufMan.Mode)
	err := logger.SetLevel(Properties.BufMan.Mode)
	if err != nil {
		panic(err)
	}

	DockerCliPool, err = NewDockerCliPool()
	if err != nil {
		panic(err)
	}

	if Properties.BufMan.UseFSStorage {
		if err := os.MkdirAll(constant.FileSavaDir, 0666); err != nil {
			panic(err)
		}
	} else {
		EsCliPool, err = NewElasticSearchCliPool()
		if err != nil {
			panic(err)
		}
	}
}

func loadFromFile() {
	viper.SetConfigFile(configFileName)
	viper.SetConfigType(configFileType)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&Properties); err != nil {
		panic(err)
	}
}

func loadFromENV() {
	if mysqlDSNENV := os.Getenv(mysqlDSNKey); mysqlDSNENV != "" {
		Properties.MySQL.MysqlDsn = mysqlDSNENV
	}
	if pageTokenSecretENV := os.Getenv(pageTokenSecretKey); pageTokenSecretENV != "" {
		Properties.BufMan.PageTokenSecret = pageTokenSecretENV
	}
	if esUsernameENV := os.Getenv(esUsernameKey); esUsernameENV != "" {
		Properties.ElasticSearch.Username = esUsernameENV
	}
	if esPasswordENV := os.Getenv(esPasswordKey); esPasswordENV != "" {
		Properties.ElasticSearch.Password = esPasswordENV
	}
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
	if Properties.BufMan.UseFSStorage {
		return nil, errors.New("use fs as storage")
	}

	c, err := elastic.NewClient(elastic.SetURL(Properties.ElasticSearch.Urls...), elastic.SetBasicAuth(Properties.ElasticSearch.Username, Properties.ElasticSearch.Password))
	if err != nil {
		return nil, err
	}

	return c, nil
}

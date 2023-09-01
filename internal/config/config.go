package config

import (
	"errors"
	"github.com/ProtobufMan/bufman/internal/constant"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"time"
)

/*
TODO ！！！这是临时的，之后改为读取配置文件
*/
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

type Config struct {
	BufMan        BufMan        `mapstructure:"bufman"`
	MySQL         MySQL         `mapstructure:"mysql"`
	Docker        Docker        `mapstructure:"docker"`
	ElasticSearch ElasticSearch `mapstructure:"elastic_search"`
}

type BufMan struct {
	Mode       string `mapstructure:"mode"`
	ServerHost string `mapstructure:"server_host"`

	PageTokenExpireTime time.Duration `mapstructure:"page_token_expire_time"`
	PageTokenSecret     string        `mapstructure:"page_token_secret"`

	UseFSStorage bool `mapstructure:"use_fs_storage"`
}

type MySQL struct {
	MysqlDsn           string        `mapstructure:"mysql_dsn"`
	MaxOpenConnections int           `mapstructure:"max_open_connections"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections"`
	MaxLifeTime        time.Duration `mapstructure:"max_life_time"`
	MaxIdleTime        time.Duration `mapstructure:"max_idle_time"`
}

type Docker struct {
	Host       string `mapstructure:"host"`
	CACertPath string `mapstructure:"ca_cert_path"`
	CertPath   string `mapstructure:"cert_path"`
	KeyPath    string `mapstructure:"key_path"`
}

type ElasticSearch struct {
	Urls     []string `mapstructure:"urls"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
}

var (
	DataBase   *gorm.DB
	Properties = &Config{}
)

func LoadConfig() {
	// 默认值
	Properties = &Config{
		BufMan: BufMan{
			Mode:       gin.DebugMode,
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

	// 从配置文件中读取
	loadFromFile()

	// 从环境变量中读取
	loadFromENV()

	gin.SetMode(Properties.BufMan.Mode)

	if Properties.BufMan.UseFSStorage {
		if err := os.MkdirAll(constant.FileSavaDir, 0666); err != nil {
			panic(err)
		}
	}

	// test docker config
	dockerCli, err := NewDockerClient()
	if err != nil {
		panic(err)
	}
	defer dockerCli.Close()

	// test es config
	_, err = NewEsClient()
	if err != nil {
		panic(err)
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

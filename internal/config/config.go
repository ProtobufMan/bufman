package config

import (
	"github.com/ProtobufMan/bufman/internal/constant"
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
	BufMan BufMan
}

type BufMan struct {
	ServerHost string
	MysqlDsn   string

	PageTokenExpireTime time.Duration
	PageTokenSecret     string
}

var (
	DataBase   *gorm.DB
	Properties *Config
)

func LoadConfig() {
	Properties = &Config{
		BufMan: BufMan{
			PageTokenExpireTime: time.Minute * 10, // 默认过期时间为10分钟
			PageTokenSecret:     "123456",
		},
	}
	Properties.BufMan.MysqlDsn = os.Getenv(mysqlDSNKey)
	Properties.BufMan.ServerHost = os.Getenv(serverHostKey)

	if err := os.MkdirAll(constant.FileSavaDir, 0666); err != nil {
		panic(err)
	}
}

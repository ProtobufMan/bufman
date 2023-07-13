package config

import (
	"gorm.io/gorm"
	"os"
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
}

var (
	DataBase   *gorm.DB
	Properties *Config
)

func LoadConfig() {
	Properties = &Config{}
	Properties.BufMan.MysqlDsn = os.Getenv(mysqlDSNKey)
	Properties.BufMan.ServerHost = os.Getenv(serverHostKey)
}

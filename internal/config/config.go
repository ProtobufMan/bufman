package config

import (
	"github.com/ProtobufMan/bufman/internal/dal"
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/driver/mysql"
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
	IDLManager IDLManager
}

type IDLManager struct {
	ServerHost string
	MysqlDsn   string
}

var (
	DataBase   *gorm.DB
	Properties *Config
)

func LoadConfig() {
	Properties = &Config{}
	Properties.IDLManager.MysqlDsn = os.Getenv(mysqlDSNKey)
	Properties.IDLManager.ServerHost = os.Getenv(serverHostKey)

	loadDatabaseConfig(Properties.IDLManager.MysqlDsn)
}

// load database for idl-manager
func loadDatabaseConfig(dsn string) {
	var db *gorm.DB
	var err error
	if dsn == "" {
		panic("MySQL DSN is empty")
	} else {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	}
	if err != nil {
		panic(err)
	} else {
		DataBase = db
		// init table
		initErr := DataBase.AutoMigrate(
			&model.Repository{},
			&model.Commit{},
			&model.Tag{},
			&model.User{},
			&model.Token{},
			&model.FileManifest{},
		)
		if initErr != nil {
			panic(initErr)
		}
	}

	dal.SetDefault(db)
}

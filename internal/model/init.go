package model

import (
	"github.com/ProtobufMan/bufman/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	dsn := config.Properties.MySQL.MysqlDsn
	var DB *gorm.DB
	var err error
	if dsn == "" {
		panic("MySQL DSN is empty")
	} else {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	}

	if err != nil {
		panic(err)
	} else {
		config.DataBase = DB
		// init table
		initErr := DB.AutoMigrate(
			&Repository{},
			&Commit{},
			&Tag{},
			&User{},
			&Token{},
			&FileManifest{},
			&FileBlob{},
			&Plugin{},
			&DockerRepo{},
		)
		if initErr != nil {
			panic(initErr)
		}
	}

	db, err := config.DataBase.DB()
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.Properties.MySQL.MaxOpenConnections)
	db.SetMaxIdleConns(config.Properties.MySQL.MaxIdleConnections)
	db.SetConnMaxLifetime(config.Properties.MySQL.MaxLifeTime)
	db.SetConnMaxIdleTime(config.Properties.MySQL.MaxIdleTime)
}

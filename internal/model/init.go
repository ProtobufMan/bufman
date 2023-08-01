package model

import (
	"github.com/ProtobufMan/bufman/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	dsn := config.Properties.BufMan.MysqlDsn
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
		config.DataBase = db
		// init table
		initErr := db.AutoMigrate(
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
}

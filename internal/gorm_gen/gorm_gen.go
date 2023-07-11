package main

import (
	"github.com/ProtobufMan/bufman/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:qwas617.@tcp(127.0.0.1:3306)/idl_manager?charset=utf8mb4&parseTime=True&loc=Local"

	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "internal/dal", // output directory
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable: true,
	})

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(model.User{}, model.Token{}, model.Repository{}, model.Tag{}, model.Commit{}, model.FileManifest{})

	// Execute the generator
	g.Execute()
}

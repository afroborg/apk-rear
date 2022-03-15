package db

import (
	"log"

	"github.com/afroborg/apk-rear/pkg/models"
	"github.com/afroborg/apk-rear/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	dbUser := utils.GetEnvVariable("POSTGRES_USER", "")
	dbPwd := utils.GetEnvVariable("POSTGRES_PASSWORD", "")
	dbHost := utils.GetEnvVariable("POSTGRES_HOST", "")
	dbName := utils.GetEnvVariable("POSTGRES_DB", "")

	dbURL := "postgres://" + dbUser + ":" + dbPwd + "@" + dbHost + "/" + dbName + "?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Alcohol{})
	db.AutoMigrate(&models.Status{})

	return db
}

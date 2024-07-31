package database

import (
	"log"

	"github.com/pedro-hos/guess-who-web/models"
	"github.com/pedro-hos/guess-who-web/util"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Connect(config util.Config, isProd bool) {

	if isProd {
		DB, err = gorm.Open(postgres.Open(config.DbUrl), &gorm.Config{})
	} else {
		DB, err = gorm.Open(sqlite.Open(config.DbUrl), &gorm.Config{})
	}

	if err != nil {
		log.Panic("failed to connect database")
	}

	// Migrate the schema
	DB.AutoMigrate(&models.State{}, &models.City{}, &models.Clue{}, &models.Card{})

}

package database

import (
	"log"

	"github.com/pedro-hos/guess-who-web/models"
	"github.com/pedro-hos/guess-who-web/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Connect() {

	config, errConfig := util.LoadConfig(".")
	if errConfig != nil {
		log.Fatal("cannot load config:", errConfig)
	}

	if config.Environment == "development" {
		DB, err = gorm.Open(sqlite.Open("database/test.db"), &gorm.Config{})
		if err != nil {
			log.Panic("failed to connect database")
		}
	}

	// Migrate the schema
	DB.AutoMigrate(&models.Card{}, &models.City{}, &models.State{}, &models.Clue{})

}

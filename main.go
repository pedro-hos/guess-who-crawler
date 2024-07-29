package main

import (
	"github.com/pedro-hos/guess-who-web/controllers"
	"github.com/pedro-hos/guess-who-web/database"
)

func main() {

	database.Connect()

	database.DB.Exec("DELETE FROM states")
	database.DB.Exec("DELETE FROM cities")
	database.DB.Exec("DELETE FROM cards")
	database.DB.Exec("DELETE FROM clues")

	controllers.RunScraper()

	// state := models.State{}
	// state.Name = "São Paulo"

	// sjc := models.City{Name: "São José dos Campos"}
	// jacarei := models.City{Name: "Jacareí"}

	// state.Cities = append(state.Cities, sjc, jacarei)

	// database.DB.Create(&state)
	// fmt.Println(state.ID)
}

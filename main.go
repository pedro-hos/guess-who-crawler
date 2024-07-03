package main

import (
	"github.com/pedro-hos/guess-who-web/controllers"
	"github.com/pedro-hos/guess-who-web/database"
)

func main() {

	database.Connect()
	controllers.RunScraper()

	// state := models.State{}
	// state.Name = "São Paulo"

	// sjc := models.City{Name: "São José dos Campos"}
	// jacarei := models.City{Name: "Jacareí"}

	// state.Cities = append(state.Cities, sjc, jacarei)

	// database.DB.Create(&state)
	// fmt.Println(state.ID)
}

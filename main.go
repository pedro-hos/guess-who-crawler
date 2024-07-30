package main

import (
	"github.com/pedro-hos/guess-who-web/cmd"
	"github.com/pedro-hos/guess-who-web/database"
)

func main() {
	database.Connect()
	cmd.Execute()
}

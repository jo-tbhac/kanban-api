package main

import (
	"os"

	"github.com/jo-tbhac/kanban-api/controllers"
	"github.com/jo-tbhac/kanban-api/db"
)

func main() {
	db := db.Get()
	defer db.Close()

	if os.Getenv("GIN_MODE") == "debug" {
		db.LogMode(true)
	}

	controllers.StartServer()
}

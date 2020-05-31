package main

import (
	"os"

	"github.com/jo-tbhac/kanban-api/controllers"
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/migration"
)

func main() {
	db := db.Get()
	defer db.Close()

	if os.Getenv("GIN_MODE") == "debug" {
		db.LogMode(true)
	}

	migration.Migrate()

	controllers.StartServer()
}

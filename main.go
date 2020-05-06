package main

import (
	"github.com/jo-tbhac/kanban-api/controllers"
	"github.com/jo-tbhac/kanban-api/db"
)

func main() {
	db := db.Get()
	defer db.Close()

	controllers.StartServer()
}

package db

import (
	"log"

	"github.com/jo-tbhac/kanban-api/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db  *gorm.DB
	err error
)

func init() {
	db, err = gorm.Open(config.Config.Database.Driver, config.Config.Database.Name)
	if err != nil {
		log.Fatalf("failed db connection: %v", err)
	}
}

func Get() *gorm.DB {
	return db
}

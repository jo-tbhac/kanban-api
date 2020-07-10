package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"local.packages/config"
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

// Get returns an instance of gorm.DB.
func Get() *gorm.DB {
	return db
}

package db

import (
	"fmt"
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
	c := fmt.Sprintf(
		"%s:%s@%s/%s?parseTime=True",
		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Name)

	db, err = gorm.Open(config.Config.Database.Driver, c)
	if err != nil {
		log.Fatalf("failed db connection: %v", err)
	}
}

// Get returns an instance of gorm.DB.
func Get() *gorm.DB {
	return db
}

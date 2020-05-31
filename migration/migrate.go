package migration

import (
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/models"
)

func Migrate() {
	db := db.Get()

	db.AutoMigrate(
		&models.User{},
		&models.Board{},
		&models.List{},
		&models.Card{},
		&models.Label{},
		&models.CardLabel{},
	)

	db.Model(&models.Board{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Label{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.List{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CardLabel{}).AddForeignKey("card_id", "cards(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.CardLabel{}).AddForeignKey("label_id", "labels(id)", "RESTRICT", "RESTRICT")
}

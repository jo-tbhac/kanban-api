package migration

import (
	"local.packages/db"
	"local.packages/entity"
)

func Migrate() {
	db := db.Get()

	db.AutoMigrate(
		&entity.User{},
		&entity.Board{},
		&entity.List{},
		&entity.Card{},
		&entity.Label{},
		&entity.CardLabel{},
	)

	db.Model(&entity.Board{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.Label{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.List{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.CardLabel{}).AddForeignKey("card_id", "cards(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.CardLabel{}).AddForeignKey("label_id", "labels(id)", "RESTRICT", "RESTRICT")
}

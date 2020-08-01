package migration

import (
	"local.packages/db"
	"local.packages/entity"
)

// Migrate ...
func Migrate() {
	db := db.Get()

	db.AutoMigrate(
		&entity.User{},
		&entity.Board{},
		&entity.List{},
		&entity.Card{},
		&entity.Label{},
		&entity.CardLabel{},
		&entity.CheckList{},
		&entity.CheckListItem{},
		&entity.File{},
	)

	db.Model(&entity.Board{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.Label{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.List{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.CardLabel{}).AddForeignKey("card_id", "cards(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.CardLabel{}).AddForeignKey("label_id", "labels(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.CardLabel{}).AddForeignKey("label_id", "labels(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.CardLabel{}).AddForeignKey("label_id", "labels(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.CheckList{}).AddForeignKey("card_id", "cards(id)", "RESTRICT", "RESTRICT")
	db.Model(&entity.CheckListItem{}).AddForeignKey("check_list_id", "check_lists(id)", "CASCADE", "RESTRICT")
	db.Model(&entity.File{}).AddForeignKey("card_id", "cards(id)", "CASCADE", "RESTRICT")
}

package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
)

type CardLabel struct {
	ID        uint
	CreatedAt time.Time
	CardID    uint
	LabelID   uint
}

func init() {
	db := db.Get()
	db.AutoMigrate(&CardLabel{})
	db.Model(&CardLabel{}).AddForeignKey("card_id", "cards(id)", "RESTRICT", "RESTRICT")
	db.Model(&CardLabel{}).AddForeignKey("label_id", "labels(id)", "RESTRICT", "RESTRICT")
}

func selectCardLabelColumn(db *gorm.DB) *gorm.DB {
	return db.Select("id, card_id, label_id")
}

func (cl *CardLabel) ValidateUID(uid uint) bool {
	db := db.Get()
	var b Board

	db.Joins("Join lists ON boards.id = lists.board_id").
		Joins("Join labels ON boards.id = labels.board_id").
		Joins("Join cards ON lists.id = cards.list_id").
		Select("user_id").
		Where("labels.id = ?", cl.LabelID).
		Where("cards.id = ?", cl.CardID).
		Find(&b)

	return b.UserID == uid
}

func (cl *CardLabel) Create() (Label, []validator.ValidationError) {
	db := db.Get()
	var l Label

	if err := db.Create(cl).Error; err != nil {
		return l, validator.FormattedValidationError(err)
	}

	db.Model(cl).Related(&l)

	return l, nil
}

func (cl *CardLabel) Find(uid uint) *gorm.DB {
	db := db.Get()

	return db.Scopes(selectCardLabelColumn).
		Joins("Join labels ON card_labels.label_id = labels.id").
		Joins("Join boards ON labels.board_id = boards.id").
		Where("boards.user_id = ?", uid).
		Where("card_labels.label_id = ?", cl.LabelID).
		Where("card_labels.card_id = ?", cl.CardID).
		First(cl)
}

func (cl *CardLabel) Delete() []validator.ValidationError {
	db := db.Get()

	if err := db.Where("label_id = ? AND card_id = ?", cl.LabelID, cl.CardID).Delete(cl).Error; err != nil {
		log.Printf("fail to delete card_label: %v", err)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

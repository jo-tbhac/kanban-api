package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
)

type Label struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-" gorm:"not null"`
	UpdatedAt time.Time  `json:"-" gorm:"not null"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" validate:"required,max=50" gorm:"not null;size:50"`
	Color     string     `json:"color" validate:"required,hexcolor" gorm:"not null;size:7"`
	BoardID   uint       `json:"-" gorm:"not null"`
}

type Labels []Label

func selectLabelColumn(db *gorm.DB) *gorm.DB {
	return db.Select("labels.id, labels.name, labels.color, labels.board_id")
}

func selectWithLabelAssociationKey(db *gorm.DB) *gorm.DB {
	return db.Select("labels.id, labels.name, labels.color, labels.board_id, card_labels.card_id")
}

func (l *Label) BeforeSave() error {
	return validator.Validate(l)
}

func (l *Label) Find(id, uid uint) *gorm.DB {
	db := db.Get()

	return db.Joins("Join boards on boards.id = labels.board_id").
		Where("boards.user_id = ?", uid).
		First(l, id)
}

func (l *Label) Create() []validator.ValidationError {
	db := db.Get()

	if err := db.Create(l).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (l *Label) Update() []validator.ValidationError {
	db := db.Get()

	if err := db.Save(l).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (l *Label) Delete() []validator.ValidationError {
	db := db.Get()

	if err := db.Delete(l).Error; err != nil {
		log.Printf("fail to delete label: %v", err)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

func (ls *Labels) GetAll(bid, uid uint) {
	db := db.Get()

	db.Scopes(selectLabelColumn).
		Joins("Join boards on boards.id = labels.board_id").
		Where("boards.user_id = ?", uid).
		Where("labels.board_id = ?", bid).
		Find(ls)
}

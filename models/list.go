package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
)

type List struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-" gorm:"not null"`
	UpdatedAt time.Time  `json:"-" gorm:"not null"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" validate:"required,max=50" gorm:"not null;size:50"`
	BoardID   uint       `json:"board_id" gorm:"not null"`
	Cards     []Card     `json:"cards"`
}

func selectListColumn(db *gorm.DB) *gorm.DB {
	return db.Select("lists.id, lists.name, lists.board_id")
}

func (l *List) BeforeSave() error {
	return validator.Validate(l)
}

func (l *List) Find(id, uid uint) *gorm.DB {
	db := db.Get()

	return db.Joins("Join boards on boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		First(l, id)
}

func (l *List) GetBoardID() {
	db := db.Get()

	db.Select("board_id").First(l, l.ID)
}

func (l *List) Create() []validator.ValidationError {
	db := db.Get()

	if err := db.Create(l).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (l *List) Update() []validator.ValidationError {
	db := db.Get()

	if err := db.Save(l).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (l *List) Delete() []validator.ValidationError {
	db := db.Get()

	if err := db.Delete(l).Error; err != nil {
		log.Printf("fail to delete list: %v", err)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

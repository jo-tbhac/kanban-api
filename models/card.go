package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
)

type Card struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Title       string     `json:"title" validate:"required,max=50"`
	Description string     `json:"description"`
	ListID      uint       `json:"list_id"`
	Labels      []Label    `json:"labels" gorm:"many2many:card_labels;"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&Card{})
}

func (c *Card) ValidateUID(uid uint) bool {
	db := db.Get()
	var b Board

	db.Joins("Join lists ON boards.id = lists.board_id").
		Select("user_id").
		Where("lists.id = ?", c.ListID).
		First(&b)

	return b.UserID == uid
}

func (c *Card) BeforeSave() error {
	return validator.Validate(c)
}

func (c *Card) Find(id, uid uint) *gorm.DB {
	db := db.Get()

	return db.Joins("Join lists ON lists.id = cards.list_id").
		Joins("Join boards ON boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		First(c, id)
}

func (c *Card) Create() []validator.ValidationError {
	db := db.Get()

	if err := db.Create(c).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (c *Card) Update() []validator.ValidationError {
	db := db.Get()

	if err := db.Save(c).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (c *Card) Delete() []validator.ValidationError {
	db := db.Get()

	if err := db.Delete(c).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

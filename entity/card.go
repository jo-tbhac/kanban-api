package entity

import (
	"time"

	"local.packages/validator"
)

type Card struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"-" gorm:"not null"`
	UpdatedAt   time.Time  `json:"-" gorm:"not null"`
	DeletedAt   *time.Time `json:"-"`
	Title       string     `json:"title" validate:"required,max=50" gorm:"not null;size:50"`
	Description string     `json:"description"`
	ListID      uint       `json:"list_id" gorm:"not null"`
	Labels      []Label    `json:"labels" gorm:"many2many:card_labels;"`
}

func (c *Card) BeforeSave() error {
	return validator.Validate(c)
}

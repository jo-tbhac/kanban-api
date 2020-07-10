package entity

import (
	"time"

	"local.packages/validator"
)

// Card is model of cards table.
type Card struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"-" gorm:"not null"`
	UpdatedAt   time.Time  `json:"-" gorm:"not null"`
	DeletedAt   *time.Time `json:"-"`
	Title       string     `json:"title" validate:"required,max=50" gorm:"not null;size:50"`
	Description string     `json:"description" gorm:"type:varchar(20000)"`
	ListID      uint       `json:"list_id" gorm:"not null"`
	Labels      []Label    `json:"labels" gorm:"many2many:card_labels;"`
	Index       int        `json:"index"`
}

// BeforeSave called before create/update a record of cards table.
// validate a field of struct and return an error if there is an invalid value
func (c *Card) BeforeSave() error {
	return validator.Validate(c)
}

package entity

import (
	"time"

	"local.packages/validator"
)

// Board is model of boards table.
type Board struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-" gorm:"not null"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" validate:"required,max=50" gorm:"size:50;not null"`
	UserID    uint       `json:"-" gorm:"not null"`
	Lists     []List     `json:"lists"`
}

// BeforeSave called before create/update a record of boards table.
// validate a field of struct and return an error if there is an invalid value
func (b *Board) BeforeSave() error {
	return validator.Validate(b)
}

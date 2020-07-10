package entity

import (
	"time"

	"local.packages/validator"
)

// Label is model of labels table.
type Label struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-" gorm:"not null"`
	UpdatedAt time.Time  `json:"-" gorm:"not null"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" validate:"required,max=50" gorm:"not null;size:50"`
	Color     string     `json:"color" validate:"required,hexcolor" gorm:"not null;size:7"`
	BoardID   uint       `json:"-" gorm:"not null"`
}

// BeforeSave called before create/update a record of labels table.
// validate a field of struct and return an error if there is an invalid value
func (l *Label) BeforeSave() error {
	return validator.Validate(l)
}

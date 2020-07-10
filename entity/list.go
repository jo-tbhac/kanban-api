package entity

import (
	"time"

	"local.packages/validator"
)

// List is model of lists table.
type List struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-" gorm:"not null"`
	UpdatedAt time.Time  `json:"-" gorm:"not null"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" validate:"required,max=50" gorm:"not null;size:50"`
	BoardID   uint       `json:"board_id" gorm:"not null"`
	Cards     []Card     `json:"cards"`
	Index     int        `json:"index"`
}

// BeforeSave called before create/update a record of lists table.
// validate a field of struct and return an error if there is an invalid value
func (l *List) BeforeSave() error {
	return validator.Validate(l)
}

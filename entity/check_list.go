package entity

import (
	"time"

	"local.packages/validator"
)

// CheckList is model of check_lists table.
type CheckList struct {
	ID        uint            `json:"id"`
	CreatedAt time.Time       `json:"-" gorm:"not null"`
	UpdatedAt time.Time       `json:"-" gorm:"not null"`
	Title     string          `json:"title" validate:"required,max=50" gorm:"not null;size:50"`
	CardID    uint            `json:"card_id" gorm:"not null"`
	Items     []CheckListItem `json:"items"`
}

// BeforeSave called before create/update a record of check_lists table.
// validate a field of struct and return an error if there is an invalid value
func (c *CheckList) BeforeSave() error {
	return validator.Validate(c)
}

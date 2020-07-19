package entity

import (
	"time"
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

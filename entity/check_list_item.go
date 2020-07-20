package entity

import (
	"time"
)

// CheckListItem is model of check_list_items table.
type CheckListItem struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"-" gorm:"not null"`
	UpdatedAt   time.Time `json:"-" gorm:"not null"`
	Name        string    `json:"name" validate:"required,max=50" gorm:"not null;size:50"`
	Check       bool      `json:"check" gorm:"default:false"`
	CheckListID uint      `json:"check_list_id" gorm:"not null"`
}

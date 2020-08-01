package entity

import (
	"time"
)

// File is model of files table.
type File struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"-" gorm:"not null"`
	UpdatedAt   time.Time `json:"-" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	URL         string    `json:"url" gorm:"not null"`
	ContentType string    `json:"content_type" gorm:"not null"`
	CardID      uint      `json:"card_id" gorm:"not null"`
}

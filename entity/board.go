package entity

import (
	"time"

	"local.packages/validator"
)

type Board struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-" gorm:"not null"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"not null"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" validate:"required,max=50" gorm:"size:50;not null"`
	UserID    uint       `json:"-" gorm:"not null"`
	Lists     []List     `json:"lists"`
}

func (b *Board) BeforeSave() error {
	return validator.Validate(b)
}

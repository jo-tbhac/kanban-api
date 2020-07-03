package entity

import (
	"time"

	"local.packages/validator"
)

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

func (l *List) BeforeSave() error {
	return validator.Validate(l)
}

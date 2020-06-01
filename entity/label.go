package entity

import (
	"time"

	"github.com/jo-tbhac/kanban-api/validator"
)

type Label struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-" gorm:"not null"`
	UpdatedAt time.Time  `json:"-" gorm:"not null"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" validate:"required,max=50" gorm:"not null;size:50"`
	Color     string     `json:"color" validate:"required,hexcolor" gorm:"not null;size:7"`
	BoardID   uint       `json:"-" gorm:"not null"`
}

func (l *Label) BeforeSave() error {
	return validator.Validate(l)
}

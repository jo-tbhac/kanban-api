package models

import (
	"time"

	"github.com/jo-tbhac/kanban-api/db"
)

type Label struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" binding:"required,max=50"`
	Color     string     `json:"color" binding:"required,max=7"`
	BoardID   uint       `json:"board_id" binding:"required"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&Label{})
	db.Model(&Label{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
}

func (l *Label) Create() error {
	db := db.Get()

	if err := db.Create(&l).Error; err != nil {
		return err
	}

	return nil
}

package models

import (
	"time"

	"github.com/jo-tbhac/kanban-api/db"
)

type Card struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Title       string     `json:"title" binding:"required,max=50"`
	Description string     `json:"description"`
	ListID      uint       `json:"list_id" binding:"required"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&Card{})
}

func (c *Card) GetBoardID() uint {
	db := db.Get()

	var l List

	db.First(&l, c.ListID)

	return l.BoardID
}

func (c *Card) Create() error {
	db := db.Get()

	if err := db.Create(c).Error; err != nil {
		return err
	}

	return nil
}

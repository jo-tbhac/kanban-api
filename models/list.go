package models

import (
	"time"

	"github.com/jo-tbhac/kanban-api/db"
)

type List struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" binding:"required,max=50"`
	BoardID   uint       `json:"board_id"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&List{})
	db.Model(&List{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
}

func (l *List) GetBoardID() {
	db := db.Get()

	db.Select("board_id").First(l, l.ID)
}

func (l *List) Create() error {
	db := db.Get()

	if err := db.Create(l).Error; err != nil {
		return err
	}

	return nil
}

func (l *List) Update() error {
	db := db.Get()

	if err := db.Save(l).Error; err != nil {
		return err
	}

	return nil
}

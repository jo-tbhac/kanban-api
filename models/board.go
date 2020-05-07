package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/db"
)

type Board struct {
	gorm.Model
	Name   string `json:"name" binding:"required,max=50"`
	UserID uint   `json:"user_id"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&Board{})
	db.Model(&Board{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
}

func (b *Board) Create() error {
	db := db.Get()

	if err := db.Create(&b).Error; err != nil {
		return err
	}

	return nil
}

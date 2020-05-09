package models

import (
	"errors"
	"log"
	"time"

	"github.com/jo-tbhac/kanban-api/db"
)

type Board struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" binding:"required,max=50"`
	UserID    uint       `json:"user_id"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&Board{})
	db.Model(&Board{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
}

func (b *Board) Create() error {
	db := db.Get()

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

func (b *Board) Get(uid uint) error {
	db := db.Get()

	db.Scopes(BoardOwnerValidation(uid)).Where("id = ?", b.ID).First(b)

	if b.UserID == UserDoesNotExist {
		log.Println("failed get board. does not match uid and board.user_id.")
		return errors.New("invalid parameters")
	}

	return nil
}

func GetAllBoard(b *[]Board, u *User) {
	db := db.Get()

	db.Model(u).Related(b)
}

package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
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

func BoardOwnerValidation(uid uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("boards.user_id = ?", uid)
	}
}

func JoinBoardTableTo(tn string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := fmt.Sprintf("left join boards on boards.id = %s.board_id", tn)
		return db.Table(tn).Joins(query)
	}
}

func RelatedBoardOwnerIsValid(bid, uid uint) bool {
	db := db.Get()

	var b Board

	db.First(&b, bid)

	return b.UserID == uid
}

func (b *Board) Create() error {
	db := db.Get()

	if err := db.Create(b).Error; err != nil {
		return err
	}

	return nil
}

func (b *Board) Update() error {
	db := db.Get()

	if err := db.Omit("user_id").Save(b).Error; err != nil {
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

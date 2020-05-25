package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
)

type Board struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" validate:"required,max=50"`
	UserID    uint       `json:"user_id"`
	Lists     []List     `json:"lists"`
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

func ValidateUID(id, uid uint) bool {
	db := db.Get()

	var b Board

	db.Select("user_id").First(&b, id)

	return b.UserID == uid
}

func (b *Board) BeforeSave() error {
	return validator.Validate(b)
}

func (b *Board) Find(id, uid uint) {
	db := db.Get()

	db.Where("user_id = ?", uid).First(b, id)
}

func (b *Board) Create() []validator.ValidationError {
	db := db.Get()

	if err := db.Create(b).Error; err != nil {
		return validator.ValidationMessages(err)
	}

	return nil
}

func (b *Board) Update() []validator.ValidationError {
	db := db.Get()

	if err := db.Save(b).Error; err != nil {
		return validator.ValidationMessages(err)
	}

	return nil
}

func (b *Board) Delete() []validator.ValidationError {
	db := db.Get()

	if r := db.Where("user_id = ?", b.UserID).Delete(b).RowsAffected; r == 0 {
		return validator.MakeErrors("invalid request")
	}

	return nil
}

func GetAllBoard(b *[]Board, u *User) {
	db := db.Get()

	db.Preload("Lists").Model(u).Related(b)
}

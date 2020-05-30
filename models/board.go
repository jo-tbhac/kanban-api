package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
)

type Board struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" validate:"required,max=50"`
	UserID    uint       `json:"-"`
	Lists     []List     `json:"lists"`
}

type Boards []Board

func init() {
	db := db.Get()
	db.AutoMigrate(&Board{})
	db.Model(&Board{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
}

func selectBoardColumn(db *gorm.DB) *gorm.DB {
	return db.Select("id, updated_at, name, user_id")
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

func (b *Board) Find(id, uid uint) *gorm.DB {
	db := db.Get()

	return db.Scopes(selectBoardColumn).
		Preload("Lists", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectListColumn)
		}).
		Preload("Lists.Cards", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectCardColumn)
		}).
		Preload("Lists.Cards.Labels", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectWithLabelAssociationKey)
		}).
		Where("user_id = ?", uid).
		First(b, id)
}

func (b *Board) Create() []validator.ValidationError {
	db := db.Get()

	if err := db.Create(b).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (b *Board) Update() []validator.ValidationError {
	db := db.Get()

	if err := db.Save(b).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (b *Board) Delete() []validator.ValidationError {
	db := db.Get()

	if r := db.Where("user_id = ?", b.UserID).Delete(b).RowsAffected; r == 0 {
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

func (bs *Boards) GetAll(u *User) {
	db := db.Get()

	db.Scopes(selectBoardColumn).Model(u).Related(bs)
}

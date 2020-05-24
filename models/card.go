package models

import (
	"time"

	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
)

type Card struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Title       string     `json:"title" validate:"required,max=50"`
	Description string     `json:"description"`
	ListID      uint       `json:"list_id" validate:"required"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&Card{})
}

func (c *Card) BeforeSave() error {
	return validator.Validate(c)
}

func (c *Card) GetBoardID() uint {
	db := db.Get()

	var l List

	db.First(&l, c.ListID)

	return l.BoardID
}

func (c *Card) Find(id, uid uint) {
	db := db.Get()

	db.Joins("Join lists ON lists.id = cards.id").
		Joins("Join boards ON boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		First(c, id)
}

func (c *Card) Create() error {
	db := db.Get()

	if err := db.Create(c).Error; err != nil {
		return err
	}

	return nil
}

func (c *Card) Update() []validator.ValidationError {
	db := db.Get()

	if err := db.Save(c).Error; err != nil {
		return validator.ValidationMessages(err)
	}

	return nil
}

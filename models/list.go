package models

import (
	"log"
	"time"

	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
)

type List struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `json:"name" validate:"required,max=50"`
	BoardID   uint       `json:"board_id"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&List{})
	db.Model(&List{}).AddForeignKey("board_id", "boards(id)", "RESTRICT", "RESTRICT")
}

func (l *List) BeforeSave() error {
	return validator.Validate(l)
}

func (l *List) Find(id, uid uint) {
	db := db.Get()

	db.Joins("Join boards on boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		First(l, id)
}

func (l *List) GetBoardID() {
	db := db.Get()

	db.Select("board_id").First(l, l.ID)
}

func (l *List) Create() []validator.ValidationError {
	db := db.Get()

	if err := db.Create(l).Error; err != nil {
		return validator.ValidationMessages(err)
	}

	return nil
}

func (l *List) Update() []validator.ValidationError {
	db := db.Get()

	if err := db.Save(l).Error; err != nil {
		return validator.ValidationMessages(err)
	}

	return nil
}

func (l *List) Delete() []validator.ValidationError {
	db := db.Get()

	if err := db.Delete(l).Error; err != nil {
		log.Printf("fail to delete list: %v", err)
		return validator.MakeErrors("invalid request")
	}

	return nil
}

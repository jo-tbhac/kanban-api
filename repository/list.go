package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/entity"
	"github.com/jo-tbhac/kanban-api/validator"
)

type ListRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) *ListRepository {
	return &ListRepository{
		db: db,
	}
}

func selectListColumn(db *gorm.DB) *gorm.DB {
	return db.Select("lists.id, lists.name, lists.board_id")
}

func (r *ListRepository) ValidateUID(id, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Select("user_id").Where("user_id = ?", uid).First(&b, id).RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

func (r *ListRepository) Find(id, uid uint) (*entity.List, []validator.ValidationError) {
	var l entity.List

	rslt := r.db.Joins("Join boards on boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		First(&l, id)

	if rslt.RecordNotFound() {
		return &l, validator.NewValidationErrors("invalid parameters")
	}

	return &l, nil
}

func (r *ListRepository) Create(name string, bid uint) (*entity.List, []validator.ValidationError) {
	l := &entity.List{
		Name:    name,
		BoardID: bid,
	}

	if err := r.db.Create(l).Error; err != nil {
		return l, validator.FormattedValidationError(err)
	}

	return l, nil
}

func (r *ListRepository) Update(l *entity.List, name string) []validator.ValidationError {
	l.Name = name

	if err := r.db.Save(l).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (r *ListRepository) Delete(l *entity.List) []validator.ValidationError {
	if err := r.db.Delete(l).Error; err != nil {
		log.Printf("fail to delete list: %v", err)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

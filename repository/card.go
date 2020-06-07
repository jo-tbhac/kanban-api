package repository

import (
	"github.com/jinzhu/gorm"

	"local.packages/entity"
	"local.packages/validator"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func selectCardColumn(db *gorm.DB) *gorm.DB {
	return db.Select("cards.id, cards.title, cards.description, cards.list_id")
}

func (r *CardRepository) ValidateUID(lid, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Joins("Join lists ON boards.id = lists.board_id").
		Select("user_id").
		Where("lists.id = ?", lid).
		Where("boards.user_id = ?", uid).
		First(&b).
		RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

func (r *CardRepository) Find(id, uid uint) (*entity.Card, []validator.ValidationError) {
	var c entity.Card

	rslt := r.db.Joins("Join lists ON lists.id = cards.list_id").
		Joins("Join boards ON boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		First(&c, id)

	if rslt.RecordNotFound() {
		return &c, validator.NewValidationErrors("invalid parameters")
	}

	return &c, nil
}

func (r *CardRepository) Create(title string, lid uint) (*entity.Card, []validator.ValidationError) {
	c := &entity.Card{
		Title:  title,
		ListID: lid,
	}

	if err := r.db.Create(c).Error; err != nil {
		return c, validator.FormattedValidationError(err)
	}

	return c, nil
}

func (r *CardRepository) UpdateTitle(c *entity.Card, title string) []validator.ValidationError {
	if err := r.db.Model(c).Updates(map[string]interface{}{"title": title}).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (r *CardRepository) UpdateDescription(c *entity.Card, description string) []validator.ValidationError {
	if err := r.db.Model(c).Updates(map[string]interface{}{"description": description}).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (r *CardRepository) Delete(c *entity.Card) []validator.ValidationError {
	if err := r.db.Delete(c).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

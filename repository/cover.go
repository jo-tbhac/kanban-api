package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"local.packages/entity"
	"local.packages/validator"
)

// CoverRepository ...
type CoverRepository struct {
	db *gorm.DB
}

// NewCoverRepository is constructor for NewCoverRepository.
func NewCoverRepository(db *gorm.DB) *CoverRepository {
	return &CoverRepository{
		db: db,
	}
}

// ValidateUID validates whether a cardID received as args was created by the login user.
func (r *CoverRepository) ValidateUID(cid, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Joins("Join lists ON boards.id = lists.board_id").
		Joins("Join cards ON lists.id = cards.list_id").
		Select("user_id").
		Where("cards.id = ?", cid).
		Where("boards.user_id = ?", uid).
		First(&b).
		RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

// Find returns a record of CheckList that found by id.
func (r *CoverRepository) Find(cid, uid uint) (*entity.Cover, []validator.ValidationError) {
	var c entity.Cover

	if r.db.Joins("Join cards ON covers.card_id = cards.id").
		Joins("Join lists ON cards.list_id = lists.id").
		Joins("Join boards ON lists.board_id = boards.id").
		Where("boards.user_id = ?", uid).
		Where("covers.card_id = ?", cid).
		First(&c).
		RecordNotFound() {
		return &c, validator.NewValidationErrors("invalid parameters")
	}

	return &c, nil
}

// Create insert a new record to a card_labels table.
func (r *CoverRepository) Create(cid, fid uint) (*entity.Cover, []validator.ValidationError) {
	c := &entity.Cover{
		CardID: cid,
		FileID: fid,
	}

	if err := r.db.Create(c).Error; err != nil {
		return c, validator.FormattedMySQLError(err)
	}

	return c, nil
}

// Update update a record's file_id in a covers table
func (r *CoverRepository) Update(c *entity.Cover, newID uint) []validator.ValidationError {
	if err := r.db.Model(c).Update("file_id", newID).Error; err != nil {
		return validator.FormattedMySQLError(err)
	}

	return nil
}

// Delete delete a record from a covers table
func (r *CoverRepository) Delete(c *entity.Cover) []validator.ValidationError {
	if rslt := r.db.Delete(c); rslt.RowsAffected == 0 {
		log.Printf("fail to delete cover: %v", rslt.Error)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

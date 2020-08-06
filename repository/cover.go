package repository

import (
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

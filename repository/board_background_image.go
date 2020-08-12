package repository

import (
	"github.com/jinzhu/gorm"
	"local.packages/entity"
	"local.packages/validator"
)

// BoardBackgroundImageRepository ...
type BoardBackgroundImageRepository struct {
	db *gorm.DB
}

// NewBoardBackgroundImageRepository is constructor for BoardBackgroundImageRepository.
func NewBoardBackgroundImageRepository(db *gorm.DB) *BoardBackgroundImageRepository {
	return &BoardBackgroundImageRepository{
		db: db,
	}
}

// Find returns a record of BoardBackgroundImage that found by id.
func (r *BoardBackgroundImageRepository) Find(id, uid uint) (*entity.BoardBackgroundImage, []validator.ValidationError) {
	var b entity.BoardBackgroundImage

	if r.db.Joins("Join boards ON board_background_images.board_id = boards.id").
		Where("boards.user_id = ?", uid).
		Where("boards.id = ?", id).
		First(&b).
		RecordNotFound() {
		return &b, validator.NewValidationErrors(ErrorInvalidSession)
	}

	return &b, nil
}

// Update update a record's background_image_id in a board_background_images table
func (r *BoardBackgroundImageRepository) Update(b *entity.BoardBackgroundImage, newID uint) []validator.ValidationError {
	if err := r.db.Model(b).Update("background_image_id", newID).Error; err != nil {
		return validator.FormattedMySQLError(err)
	}

	return nil
}

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

// ValidateUID validates whether a boardID received as args was created by the login user.
func (r *BoardBackgroundImageRepository) ValidateUID(id, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Select("user_id").Where("user_id = ?", uid).First(&b, id).RecordNotFound() {
		return validator.NewValidationErrors(ErrorInvalidSession)
	}

	return nil
}

// Create insert a new record to a board_background_images table.
func (r *BoardBackgroundImageRepository) Create(bid, iid uint) (*entity.BoardBackgroundImage, []validator.ValidationError) {
	b := &entity.BoardBackgroundImage{
		BoardID:           bid,
		BackgroundImageID: iid,
	}

	if err := r.db.Create(b).Error; err != nil {
		return b, validator.FormattedMySQLError(err)
	}

	return b, nil
}

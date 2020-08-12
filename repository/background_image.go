package repository

import (
	"github.com/jinzhu/gorm"
	"local.packages/entity"
)

// BackgroundImageRepository ...
type BackgroundImageRepository struct {
	db *gorm.DB
}

// NewBackgroundImageRepository is constructor for BackgroundImageRepository.
func NewBackgroundImageRepository(db *gorm.DB) *BackgroundImageRepository {
	return &BackgroundImageRepository{
		db: db,
	}
}

// GetAll returns slice of BackgroundImage's record.
func (r *BackgroundImageRepository) GetAll() *[]entity.BackgroundImage {
	var bs []entity.BackgroundImage

	r.db.Find(&bs)

	return &bs
}

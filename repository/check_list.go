package repository

import (
	"github.com/jinzhu/gorm"
	"local.packages/entity"
	"local.packages/validator"
)

// CheckListRepository ...
type CheckListRepository struct {
	db *gorm.DB
}

// NewCheckListRepository is constructor for CheckListRepository.
func NewCheckListRepository(db *gorm.DB) *CheckListRepository {
	return &CheckListRepository{
		db: db,
	}
}

// Create insert a new record to a check_lists table.
func (r *CheckListRepository) Create(title string, cid uint) (*entity.CheckList, []validator.ValidationError) {
	c := &entity.CheckList{
		Title:  title,
		CardID: cid,
	}

	if err := r.db.Create(c).Error; err != nil {
		return c, validator.FormattedValidationError(err)
	}

	return c, nil
}

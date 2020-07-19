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

// Update update a record's title in a check_lists table
func (r *CheckListRepository) Update(c *entity.CheckList, title string) []validator.ValidationError {
	if err := r.db.Model(c).Update("title", title).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

// Delete delete a record from a checl_lists table.
func (r *CheckListRepository) Delete(c *entity.CheckList) []validator.ValidationError {
	if err := r.db.Delete(c).Error; err != nil {
		return validator.FormattedMySQLError(err)
	}

	return nil
}

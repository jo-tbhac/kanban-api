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

// ValidateUID validates whether a cardID received as args was created by the login user.
func (r *CheckListRepository) ValidateUID(cid, uid uint) []validator.ValidationError {
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
// func (r *CheckListRepository) Find(id, uid) *entity.CheckList {

// }

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
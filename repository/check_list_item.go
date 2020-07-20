package repository

import (
	"log"

	"github.com/jinzhu/gorm"
	"local.packages/entity"
	"local.packages/validator"
)

// CheckListItemRepository ...
type CheckListItemRepository struct {
	db *gorm.DB
}

// NewCheckListItemRepository is constructor for CheckListItemRepository.
func NewCheckListItemRepository(db *gorm.DB) *CheckListItemRepository {
	return &CheckListItemRepository{
		db: db,
	}
}

// ValidateUID validates whether a checkListID received as args was created by the login user.
func (r *CheckListItemRepository) ValidateUID(cid, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Joins("Join lists ON boards.id = lists.board_id").
		Joins("Join cards ON lists.id = cards.list_id").
		Joins("Join check_lists ON cards.id = check_lists.card_id").
		Select("user_id").
		Where("check_lists.id = ?", cid).
		Where("boards.user_id = ?", uid).
		First(&b).
		RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

// Find returns a record of CheckListItem that found by id.
func (r *CheckListItemRepository) Find(id, uid uint) (*entity.CheckListItem, []validator.ValidationError) {
	var item entity.CheckListItem

	if r.db.Joins("Join check_lists ON check_list_items.check_list_id = check_lists.id").
		Joins("Join cards ON check_lists.card_id = cards.id").
		Joins("Join lists ON cards.list_id = lists.id").
		Joins("Join boards ON lists.board_id = boards.id").
		Where("boards.user_id = ?", uid).
		First(&item, id).
		RecordNotFound() {
		return &item, validator.NewValidationErrors("invalid parameters")
	}

	return &item, nil
}

// Create insert a new record to a check_lists table.
func (r *CheckListItemRepository) Create(name string, cid uint) (*entity.CheckListItem, []validator.ValidationError) {
	item := &entity.CheckListItem{
		Name:        name,
		CheckListID: cid,
	}

	if err := r.db.Create(item).Error; err != nil {
		return item, validator.FormattedValidationError(err)
	}

	return item, nil
}

// Update update a record's name in a check_list_items table.
func (r *CheckListItemRepository) Update(item *entity.CheckListItem, name string) []validator.ValidationError {
	if err := r.db.Model(item).Update("name", name).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

// Check update a record's column of check in a check_list_items table.
func (r *CheckListItemRepository) Check(item *entity.CheckListItem, check bool) []validator.ValidationError {
	if err := r.db.Model(item).Update("check", check).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}


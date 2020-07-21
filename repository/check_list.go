package repository

import (
	"log"

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

func selectCheckListColumn(db *gorm.DB) *gorm.DB {
	return db.Select("check_lists.id, check_lists.title, check_lists.card_id")
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
func (r *CheckListRepository) Find(id, uid uint) (*entity.CheckList, []validator.ValidationError) {
	var c entity.CheckList

	if r.db.Joins("Join cards ON check_lists.card_id = cards.id").
		Joins("Join lists ON cards.list_id = lists.id").
		Joins("Join boards ON lists.board_id = boards.id").
		Where("boards.user_id = ?", uid).
		First(&c, id).
		RecordNotFound() {
		return &c, validator.NewValidationErrors("invalid parameters")
	}

	return &c, nil
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

// Delete delete a record from a check_lists table.
func (r *CheckListRepository) Delete(c *entity.CheckList) []validator.ValidationError {
	if rslt := r.db.Delete(c); rslt.RowsAffected == 0 {
		log.Printf("fail to delete check_list: %v", rslt.Error)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

// GetAll returns slice of CheckList's record.
func (r *CheckListRepository) GetAll(bid, uid uint) *[]entity.CheckList {
	var cs []entity.CheckList

	r.db.Scopes(selectCheckListColumn).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectCheckListItemColumn)
		}).
		Joins("Join cards ON check_lists.card_id = cards.id").
		Joins("Join lists ON cards.list_id = lists.id").
		Joins("Join boards ON lists.board_id = boards.id").
		Where("boards.id = ?", bid).
		Where("boards.user_id = ?", uid).
		Find(&cs)

	return &cs
}

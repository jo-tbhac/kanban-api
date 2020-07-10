package repository

import (
	"log"

	"github.com/jinzhu/gorm"

	"local.packages/entity"
	"local.packages/validator"
)

// LabelRepository ...
type LabelRepository struct {
	db *gorm.DB
}

// NewLabelRepository is constructor for LabelRepository.
func NewLabelRepository(db *gorm.DB) *LabelRepository {
	return &LabelRepository{
		db: db,
	}
}

func selectLabelColumn(db *gorm.DB) *gorm.DB {
	return db.Select("labels.id, labels.name, labels.color, labels.board_id")
}

func selectWithLabelAssociationKey(db *gorm.DB) *gorm.DB {
	return db.Select("labels.id, labels.name, labels.color, labels.board_id, card_labels.card_id")
}

// ValidateUID validates whether a boardID received as args was created by the login user.
func (r *LabelRepository) ValidateUID(id, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Select("user_id").Where("user_id = ?", uid).First(&b, id).RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

// Find returns a record of Label that found by id.
func (r *LabelRepository) Find(id, uid uint) (*entity.Label, []validator.ValidationError) {
	var l entity.Label

	rslt := r.db.Joins("Join boards on boards.id = labels.board_id").
		Where("boards.user_id = ?", uid).
		First(&l, id)

	if rslt.RecordNotFound() {
		return &l, validator.NewValidationErrors("invalid parameters")
	}

	return &l, nil
}

// Create insert a new record to a labels table.
func (r *LabelRepository) Create(name, color string, bid uint) (*entity.Label, []validator.ValidationError) {
	l := &entity.Label{
		Name:    name,
		Color:   color,
		BoardID: bid,
	}

	if err := r.db.Create(l).Error; err != nil {
		return l, validator.FormattedValidationError(err)
	}

	return l, nil
}

// Update update a record in a labels table.
func (r *LabelRepository) Update(l *entity.Label, name, color string) []validator.ValidationError {
	if err := r.db.Model(l).Updates(map[string]interface{}{"name": name, "color": color}).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

// Delete delete a record from a labels table.
// use soft delete.
func (r *LabelRepository) Delete(l *entity.Label) []validator.ValidationError {
	if rslt := r.db.Delete(l); rslt.RowsAffected == 0 {
		log.Printf("fail to delete label: %v", rslt.Error)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

// GetAll returns slice of Label's record.
func (r *LabelRepository) GetAll(bid, uid uint) *[]entity.Label {
	var ls []entity.Label

	r.db.Scopes(selectLabelColumn).
		Joins("Join boards on boards.id = labels.board_id").
		Where("boards.user_id = ?", uid).
		Where("labels.board_id = ?", bid).
		Find(&ls)

	return &ls
}

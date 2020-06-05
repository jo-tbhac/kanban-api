package repository

import (
	"log"

	"github.com/jinzhu/gorm"

	"local.packages/entity"
	"local.packages/validator"
)

type LabelRepository struct {
	db *gorm.DB
}

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

func (r *LabelRepository) ValidateUID(id, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Select("user_id").Where("user_id = ?", uid).First(&b, id).RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

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

func (r *LabelRepository) Update(l *entity.Label, name, color string) []validator.ValidationError {
	if err := r.db.Model(l).Updates(map[string]interface{}{"name": name, "color": color}).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (r *LabelRepository) Delete(l *entity.Label) []validator.ValidationError {
	if err := r.db.Delete(l).Error; err != nil {
		log.Printf("fail to delete label: %v", err)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

func (r *LabelRepository) GetAll(bid, uid uint) *[]entity.Label {
	var ls []entity.Label

	r.db.Scopes(selectLabelColumn).
		Joins("Join boards on boards.id = labels.board_id").
		Where("boards.user_id = ?", uid).
		Where("labels.board_id = ?", bid).
		Find(&ls)

	return &ls
}

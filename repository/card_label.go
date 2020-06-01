package repository

import (
	"log"

	"github.com/jinzhu/gorm"

	"local.packages/entity"
	"local.packages/validator"
)

type CardLabelRepository struct {
	db *gorm.DB
}

func NewCardLabelRepository(db *gorm.DB) *CardLabelRepository {
	return &CardLabelRepository{
		db: db,
	}
}

func (r *CardLabelRepository) ValidateUID(lid, cid, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Joins("Join lists ON boards.id = lists.board_id").
		Joins("Join labels ON boards.id = labels.board_id").
		Joins("Join cards ON lists.id = cards.list_id").
		Select("user_id").
		Where("labels.id = ?", lid).
		Where("cards.id = ?", cid).
		Where("boards.user_id = ?", uid).
		Find(&b).
		RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

func (r *CardLabelRepository) Create(lid, cid uint) (*entity.Label, []validator.ValidationError) {
	var l entity.Label
	cl := &entity.CardLabel{
		LabelID: lid,
		CardID:  cid,
	}

	if err := r.db.Create(cl).Error; err != nil {
		return &l, validator.FormattedMySQLError(err)
	}

	r.db.Model(cl).Related(&l)

	return &l, nil
}

func (r *CardLabelRepository) Find(lid, cid, uid uint) (*entity.CardLabel, []validator.ValidationError) {
	var cl entity.CardLabel

	if r.db.Joins("Join labels ON card_labels.label_id = labels.id").
		Joins("Join boards ON labels.board_id = boards.id").
		Where("boards.user_id = ?", uid).
		Where("card_labels.label_id = ?", lid).
		Where("card_labels.card_id = ?", cid).
		First(&cl).
		RecordNotFound() {
		return &cl, validator.NewValidationErrors("invalid parameters")
	}

	return &cl, nil
}

func (r *CardLabelRepository) Delete(cl *entity.CardLabel) []validator.ValidationError {
	if err := r.db.Delete(cl).Error; err != nil {
		log.Printf("fail to delete card_label: %v", err)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

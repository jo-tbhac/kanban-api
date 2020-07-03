package repository

import (
	"github.com/jinzhu/gorm"

	"local.packages/entity"
	"local.packages/validator"
)

type BoardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{
		db: db,
	}
}

func selectBoardColumn(db *gorm.DB) *gorm.DB {
	return db.Select("id, updated_at, name, user_id")
}

func (r *BoardRepository) Find(id, uid uint) (*entity.Board, []validator.ValidationError) {
	var b entity.Board

	rslt := r.db.Scopes(selectBoardColumn).
		Preload("Lists", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectListColumn).Order("lists.index asc")
		}).
		Preload("Lists.Cards", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectCardColumn)
		}).
		Preload("Lists.Cards.Labels", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectWithLabelAssociationKey)
		}).
		Where("user_id = ?", uid).
		First(&b, id)

	if rslt.RecordNotFound() {
		return &b, validator.NewValidationErrors("invalid parameters")
	}

	return &b, nil
}

func (r *BoardRepository) FindWithoutPreload(id, uid uint) (*entity.Board, []validator.ValidationError) {
	var b entity.Board

	if r.db.Scopes(selectBoardColumn).Where("user_id = ?", uid).First(&b, id).RecordNotFound() {
		return &b, validator.NewValidationErrors("invalid parameters")
	}

	return &b, nil
}

func (r *BoardRepository) Create(name string, uid uint) (*entity.Board, []validator.ValidationError) {
	b := &entity.Board{
		Name:   name,
		UserID: uid,
	}

	if err := r.db.Create(b).Error; err != nil {
		return b, validator.FormattedValidationError(err)
	}

	return b, nil
}

func (r *BoardRepository) Update(b *entity.Board, name string) []validator.ValidationError {
	if err := r.db.Set("gorm:association_autoupdate", false).Model(b).Update("name", name).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (r *BoardRepository) Delete(id, uid uint) []validator.ValidationError {
	if rslt := r.db.Where("id = ? AND user_id = ?", id, uid).Delete(&entity.Board{}).RowsAffected; rslt == 0 {
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

func (r *BoardRepository) GetAll(uid uint) *[]entity.Board {
	var bs []entity.Board

	r.db.Scopes(selectBoardColumn).Where("user_id = ?", uid).Find(&bs)

	return &bs
}

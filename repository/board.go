package repository

import (
	"reflect"

	"github.com/jinzhu/gorm"

	"local.packages/entity"
	"local.packages/validator"
)

// BoardRepository ...
type BoardRepository struct {
	db *gorm.DB
}

// NewBoardRepository is constructor for BoardRepository.
func NewBoardRepository(db *gorm.DB) *BoardRepository {
	return &BoardRepository{
		db: db,
	}
}

func selectBoardColumn(db *gorm.DB) *gorm.DB {
	return db.Select("id, updated_at, name, user_id")
}

// Find returns a record of Board that contains related model's records.
func (r *BoardRepository) Find(id, uid uint) (*entity.Board, []validator.ValidationError) {
	var b entity.Board

	rslt := r.db.Scopes(selectBoardColumn).
		Preload("BackgroundImage").
		Preload("Lists", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectListColumn).Order("lists.index asc")
		}).
		Preload("Lists.Cards", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectCardColumn).Order("cards.index asc")
		}).
		Preload("Lists.Cards.Labels", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(selectWithLabelAssociationKey)
		}).
		Preload("Lists.Cards.Cover").
		Where("user_id = ?", uid).
		First(&b, id)

	if rslt.RecordNotFound() {
		return &b, validator.NewValidationErrors(ErrorRecordNotFound)
	}

	return &b, nil
}

// FindWithoutPreload returns a record of Board without related model's records.
func (r *BoardRepository) FindWithoutPreload(id, uid uint) (*entity.Board, []validator.ValidationError) {
	var b entity.Board

	if r.db.Scopes(selectBoardColumn).Where("user_id = ?", uid).First(&b, id).RecordNotFound() {
		return &b, validator.NewValidationErrors(ErrorRecordNotFound)
	}

	return &b, nil
}

// Create insert a new record to a boards table.
func (r *BoardRepository) Create(name string, iid, uid uint) (*entity.Board, []validator.ValidationError) {
	b := &entity.Board{
		Name:   name,
		UserID: uid,
	}

	i := &entity.BoardBackgroundImage{
		BackgroundImageID: iid,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(b).Error; err != nil {
			return err
		}

		i.BoardID = b.ID

		if err := tx.Create(i).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t := reflect.TypeOf(err)

		if t.String() == "*mysql.MySQLError" {
			return b, validator.FormattedMySQLError(err)
		}
		return b, validator.FormattedValidationError(err)
	}

	b.BackgroundImage = i

	return b, nil
}

// Update update a record in a boards table.
func (r *BoardRepository) Update(b *entity.Board, name string) []validator.ValidationError {
	if err := r.db.Set("gorm:association_autoupdate", false).Model(b).Update("name", name).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

// Delete delete a record from a boards table.
// use soft delete.
func (r *BoardRepository) Delete(id, uid uint) []validator.ValidationError {
	if rslt := r.db.Where("id = ? AND user_id = ?", id, uid).Delete(&entity.Board{}).RowsAffected; rslt == 0 {
		return validator.NewValidationErrors(ErrorInvalidRequest)
	}

	return nil
}

// GetAll returns slice of Board's record.
func (r *BoardRepository) GetAll(uid uint) *[]entity.Board {
	var bs []entity.Board

	r.db.Scopes(selectBoardColumn).Preload("BackgroundImage").Where("user_id = ?", uid).Find(&bs)

	return &bs
}

// Search returns ids of Board that found by Board's name.
func (r *BoardRepository) Search(name string, uid uint) []uint {
	var ids []uint

	r.db.Model(&entity.Board{}).
		Where("user_id = ?", uid).
		Where("name LIKE ?", "%"+name+"%").
		Pluck("id", &ids)

	return ids
}

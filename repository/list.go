package repository

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"local.packages/entity"
	"local.packages/validator"
)

type ListRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) *ListRepository {
	return &ListRepository{
		db: db,
	}
}

func selectListColumn(db *gorm.DB) *gorm.DB {
	return db.Select("lists.id, lists.name, lists.board_id, lists.index")
}

func (r *ListRepository) ValidateUID(id, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Select("user_id").Where("user_id = ?", uid).First(&b, id).RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

func (r *ListRepository) Find(id, uid uint) (*entity.List, []validator.ValidationError) {
	var l entity.List

	rslt := r.db.Joins("Join boards on boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		First(&l, id)

	if rslt.RecordNotFound() {
		return &l, validator.NewValidationErrors("invalid parameters")
	}

	return &l, nil
}

func (r *ListRepository) Create(name string, bid uint) (*entity.List, []validator.ValidationError) {
	pl := &entity.List{}
	if r := r.db.Select("`index`").Where("board_id = ?", bid).Order("`index` desc").Take(pl).RowsAffected; r > 0 {
		pl.Index = pl.Index + 1
	}

	l := &entity.List{
		Name:    name,
		Index:   pl.Index,
		BoardID: bid,
	}

	if err := r.db.Create(l).Error; err != nil {
		return l, validator.FormattedValidationError(err)
	}

	return l, nil
}

func (r *ListRepository) Update(l *entity.List, name string) []validator.ValidationError {
	if err := r.db.Model(l).Updates(map[string]interface{}{"name": name}).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

func (r *ListRepository) UpdateIndex(params []struct {
	ID    uint
	Index int
}) []validator.ValidationError {
	ids := make([]string, 0, len(params))
	values := make([]string, 0, len(params))

	for _, p := range params {
		ids = append(ids, strconv.Itoa(int(p.ID)))
		values = append(values, strconv.Itoa(p.Index))
	}

	joinedIDs := strings.Join(ids, ",")
	joinedValues := strings.Join(values, ",")
	q := fmt.Sprintf("UPDATE `lists` SET `index` = ELT(FIELD(id,%s),%s) WHERE id IN (%s)", joinedIDs, joinedValues, joinedIDs)

	if err := r.db.Exec(q).Error; err != nil {
		return validator.FormattedValidationError(err)
	}
	return nil
}

func (r *ListRepository) Delete(l *entity.List) []validator.ValidationError {
	if rslt := r.db.Model(l).UpdateColumns(map[string]interface{}{"deleted_at": time.Now(), "index": 0}); rslt.RowsAffected == 0 {
		log.Printf("fail to delete list: %v", rslt.Error)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

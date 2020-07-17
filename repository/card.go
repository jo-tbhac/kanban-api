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

// CardRepository ...
type CardRepository struct {
	db *gorm.DB
}

// NewCardRepository is constructor for CardRepository.
func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func selectCardColumn(db *gorm.DB) *gorm.DB {
	return db.Select("cards.id, cards.title, cards.description, cards.list_id, cards.index")
}

// ValidateUID validates whether a listID received as args was created by the login user.
func (r *CardRepository) ValidateUID(lid, uid uint) []validator.ValidationError {
	var b entity.Board

	if r.db.Joins("Join lists ON boards.id = lists.board_id").
		Select("user_id").
		Where("lists.id = ?", lid).
		Where("boards.user_id = ?", uid).
		First(&b).
		RecordNotFound() {
		return validator.NewValidationErrors("invalid parameters")
	}

	return nil
}

// Find returns a record of Card that found by id.
func (r *CardRepository) Find(id, uid uint) (*entity.Card, []validator.ValidationError) {
	var c entity.Card

	rslt := r.db.Joins("Join lists ON lists.id = cards.list_id").
		Joins("Join boards ON boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		First(&c, id)

	if rslt.RecordNotFound() {
		return &c, validator.NewValidationErrors("invalid parameters")
	}

	return &c, nil
}

// Create insert a new record to a cards table.
func (r *CardRepository) Create(title string, lid uint) (*entity.Card, []validator.ValidationError) {
	pc := &entity.Card{}
	if r := r.db.Select("`index`").Where("list_id = ?", lid).Order("`index` desc").Take(pc).RowsAffected; r > 0 {
		pc.Index = pc.Index + 1
	}

	c := &entity.Card{
		Title:  title,
		ListID: lid,
		Index:  pc.Index,
	}

	if err := r.db.Create(c).Error; err != nil {
		return c, validator.FormattedValidationError(err)
	}

	return c, nil
}

// UpdateTitle update a record's title in a cards table.
func (r *CardRepository) UpdateTitle(c *entity.Card, title string) []validator.ValidationError {
	if err := r.db.Model(c).Updates(map[string]interface{}{"title": title}).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

// UpdateDescription update a record's description in a cards table.
func (r *CardRepository) UpdateDescription(c *entity.Card, description string) []validator.ValidationError {
	if err := r.db.Model(c).Updates(map[string]interface{}{"description": description}).Error; err != nil {
		return validator.FormattedValidationError(err)
	}

	return nil
}

// UpdateIndex update Card's order that recieved as args.
func (r *CardRepository) UpdateIndex(params []struct {
	ID     uint `json:"id"`
	Index  int  `json:"index"`
	ListID uint `json:"list_id"`
}) []validator.ValidationError {
	ids := make([]string, 0, len(params))
	listIds := make([]string, 0, len(params))
	values := make([]string, 0, len(params))

	for _, p := range params {
		ids = append(ids, strconv.Itoa(int(p.ID)))
		listIds = append(listIds, strconv.Itoa(int(p.ListID)))
		values = append(values, strconv.Itoa(p.Index))
	}

	joinedIDs := strings.Join(ids, ",")
	joinedListIDs := strings.Join(listIds, ",")
	joinedValues := strings.Join(values, ",")
	q := fmt.Sprintf(
		"UPDATE `cards` SET `index` = ELT(FIELD(id,%s),%s), `list_id` = ELT(FIELD(id,%s),%s) WHERE id IN (%s)",
		joinedIDs,
		joinedValues,
		joinedIDs,
		joinedListIDs,
		joinedIDs)

	if err := r.db.Exec(q).Error; err != nil {
		return validator.FormattedValidationError(err)
	}
	return nil
}

// Delete delete a record from a cards table.
// use soft delete.
func (r *CardRepository) Delete(c *entity.Card) []validator.ValidationError {
	if rslt := r.db.Model(c).UpdateColumns(map[string]interface{}{"deleted_at": time.Now(), "index": 0}); rslt.RowsAffected == 0 {
		log.Printf("fail to delete card: %v", rslt.Error)
		return validator.NewValidationErrors("invalid request")
	}

	return nil
}

// Search returns ids of Card that found by Card's title.
func (r *CardRepository) Search(bid, uid uint, title string) []uint {
	var ids []uint

	r.db.Model(&entity.Card{}).
		Joins("Join lists ON lists.id = cards.list_id").
		Joins("Join boards ON boards.id = lists.board_id").
		Where("boards.user_id = ?", uid).
		Where("boards.id = ?", bid).
		Where("cards.title LIKE ?", "%"+title+"%").
		Where("lists.deleted_at IS NULL").
		Order("cards.list_id asc").
		Pluck("cards.id", &ids)

	return ids
}

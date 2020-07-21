package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"local.packages/entity"
	"local.packages/utils"
)

func TestShouldSuccessfullyValidateUIDOnCardLabelRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardLabelRepository(db)

	userID := uint(1)
	cardID := uint(2)
	labelID := uint(3)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		Join lists ON boards.id = lists.board_id
		Join labels ON boards.id = labels.board_id
		Join cards ON lists.id = cards.list_id
		WHERE 'boards'.'deleted_at' IS NULL AND ((labels.id = ?) AND (cards.id = ?) AND (boards.user_id = ?))`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(labelID, cardID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userID))

	if err := r.ValidateUID(labelID, cardID, userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureValidateUIDOnCardLabelRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardLabelRepository(db)

	userID := uint(1)
	cardID := uint(2)
	labelID := uint(3)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		Join lists ON boards.id = lists.board_id
		Join labels ON boards.id = labels.board_id
		Join cards ON lists.id = cards.list_id
		WHERE 'boards'.'deleted_at' IS NULL AND ((labels.id = ?) AND (cards.id = ?) AND (boards.user_id = ?))`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(labelID, cardID, userID).
		WillReturnError(gorm.ErrRecordNotFound)

	err := r.ValidateUID(labelID, cardID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyCreateCardLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardLabelRepository(db)

	labelID := uint(1)
	cardID := uint(2)

	insertQuery := "INSERT INTO `card_labels` (`card_id`,`label_id`) VALUES (?,?)"
	findQuery := "SELECT * FROM `labels` WHERE `labels`.`deleted_at` IS NULL AND ((`id` = ?))"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).
		WithArgs(cardID, labelID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta(findQuery)).
		WithArgs(labelID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(labelID))

	l, err := r.Create(labelID, cardID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, l.ID, labelID)
}

func TestShouldNotCreateCardLabelWhenDuplicatePrimaryKey(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardLabelRepository(db)

	labelID := uint(1)
	cardID := uint(2)

	query := "INSERT INTO `card_labels` (`card_id`,`label_id`) VALUES (?,?)"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(cardID, labelID).
		WillReturnError(fmt.Errorf("Error 1062: Duplicate entry '%d-%d' for key 'email'", cardID, labelID))

	mock.ExpectRollback()

	_, err := r.Create(labelID, cardID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, fmt.Sprintf("%d-%d has already been taken", cardID, labelID))
}

func TestShouldSuccessfullyFindCardLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardLabelRepository(db)

	labelID := uint(1)
	cardID := uint(2)
	userID := uint(3)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'card_labels'.*
		FROM 'card_labels'
		Join labels ON card_labels.label_id = labels.id
		Join boards ON labels.board_id = boards.id
		WHERE (boards.user_id = ?) AND (card_labels.label_id = ?) AND (card_labels.card_id = ?)
		ORDER BY 'card_labels'.'card_id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, labelID, cardID).
		WillReturnRows(sqlmock.NewRows([]string{"card_id", "label_id"}).AddRow(cardID, labelID))

	cl, err := r.Find(labelID, cardID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, cl.CardID, cardID)
	assert.Equal(t, cl.LabelID, labelID)
}

func TestShouldNotFindCardLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardLabelRepository(db)

	labelID := uint(1)
	cardID := uint(2)
	userID := uint(3)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'card_labels'.*
		FROM 'card_labels'
		Join labels ON card_labels.label_id = labels.id
		Join boards ON labels.board_id = boards.id
		WHERE (boards.user_id = ?) AND (card_labels.label_id = ?) AND (card_labels.card_id = ?)
		ORDER BY 'card_labels'.'card_id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, labelID, cardID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(labelID, cardID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyDeleteCardLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardLabelRepository(db)

	cl := &entity.CardLabel{
		CardID:  uint(1),
		LabelID: uint(2),
	}

	query := utils.ReplaceQuotationForQuery(`
		DELETE FROM 'card_labels'
		WHERE 'card_labels'.'card_id' = ? AND 'card_labels'.'label_id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(cl.CardID, cl.LabelID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Delete(cl); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldNotDeleteCardLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardLabelRepository(db)

	cl := &entity.CardLabel{
		CardID:  uint(1),
		LabelID: uint(2),
	}

	query := utils.ReplaceQuotationForQuery(`
		DELETE FROM 'card_labels'
		WHERE 'card_labels'.'card_id' = ? AND 'card_labels'.'label_id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(cl.CardID, cl.LabelID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err := r.Delete(cl)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid request")
}

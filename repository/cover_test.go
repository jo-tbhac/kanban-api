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

func TestShouldSuccessfullyValidateUIDOnCoverRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCoverRepository(db)

	userID := uint(1)
	cardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id FROM 'boards'
		Join lists ON boards.id = lists.board_id
		Join cards ON lists.id = cards.list_id
		WHERE 'boards'.'deleted_at' IS NULL AND ((cards.id = ?) AND (boards.user_id = ?))
		ORDER BY 'boards'.'id' ASC LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cardID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userID))

	if err := r.ValidateUID(cardID, userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureValidateUIDOnCoverRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCoverRepository(db)

	userID := uint(1)
	cardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id FROM 'boards'
		Join lists ON boards.id = lists.board_id
		Join cards ON lists.id = cards.list_id
		WHERE 'boards'.'deleted_at' IS NULL AND ((cards.id = ?) AND (boards.user_id = ?))
		ORDER BY 'boards'.'id' ASC LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(cardID, userID).
		WillReturnError(gorm.ErrRecordNotFound)

	err := r.ValidateUID(cardID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyCreateCover(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	cardID := uint(1)
	fileID := uint(2)

	r := NewCoverRepository(db)

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'covers' ('card_id','file_id')
		VALUES (?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(cardID, fileID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	c, err := r.Create(cardID, fileID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, c.CardID, cardID)
	assert.Equal(t, c.FileID, fileID)
}

func TestShoulNotCreateCover(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	cardID := uint(1)
	fileID := uint(2)

	r := NewCoverRepository(db)

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'covers' ('card_id','file_id')
		VALUES (?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(cardID, fileID).
		WillReturnError(fmt.Errorf("Error 1062: Duplicate entry '%d' for key", cardID))

	mock.ExpectRollback()

	_, err := r.Create(cardID, fileID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, fmt.Sprintf("%d has already been taken", cardID))
}

func TestShouldSuccessfullyFindCover(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCoverRepository(db)

	cardID := uint(1)
	userID := uint(2)
	fileID := uint(3)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'covers'.*
		FROM 'covers'
		Join cards ON covers.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.user_id = ?) AND (covers.card_id = ?) AND (covers.file_id = ?)
		ORDER BY 'covers'.'card_id' ASC LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, cardID, fileID).
		WillReturnRows(sqlmock.NewRows([]string{"card_id", "file_id"}).AddRow(cardID, fileID))

	c, err := r.Find(cardID, fileID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, c.CardID, cardID)
	assert.Equal(t, c.FileID, fileID)
}

func TestShouldNotFindCover(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCoverRepository(db)

	cardID := uint(1)
	userID := uint(2)
	fileID := uint(3)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'covers'.*
		FROM 'covers'
		Join cards ON covers.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.user_id = ?) AND (covers.card_id = ?) AND (covers.file_id = ?)
		ORDER BY 'covers'.'card_id' ASC LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, cardID, fileID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(cardID, fileID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyUpdateCover(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	newFileID := uint(3)

	c := &entity.Cover{
		CardID: uint(1),
		FileID: uint(3),
	}

	r := NewCoverRepository(db)

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'covers'
		SET 'file_id' = ?
		WHERE 'covers'.'card_id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(c.FileID, c.CardID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Update(c, newFileID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, c.FileID, newFileID)
}

func TestShouldFailureUpdateCover(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	newFileID := uint(3)

	c := &entity.Cover{
		CardID: uint(1),
		FileID: uint(3),
	}

	r := NewCoverRepository(db)

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'covers'
		SET 'file_id' = ?
		WHERE 'covers'.'card_id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(c.FileID, c.CardID).
		WillReturnError(fmt.Errorf("Error 1062: Duplicate entry '%d' for key", c.CardID))

	mock.ExpectRollback()

	err := r.Update(c, newFileID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, fmt.Sprintf("%d has already been taken", c.CardID))
}

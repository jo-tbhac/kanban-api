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

func TestShouldSuccessfullyValidateUIDOnFileRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	userID := uint(1)
	cardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		Join lists ON boards.id = lists.board_id
		Join cards ON lists.id = cards.list_id
		WHERE 'boards'.'deleted_at' IS NULL AND ((cards.id = ?) AND (boards.user_id = ?))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

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

func TestShouldFailureValidateUIDOnFileRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	userID := uint(1)
	cardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		Join lists ON boards.id = lists.board_id
		Join cards ON lists.id = cards.list_id
		WHERE 'boards'.'deleted_at' IS NULL AND ((cards.id = ?) AND (boards.user_id = ?))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

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

	assert.Equal(t, err[0].Text, ErrorInvalidSession)
}

func TestShouldSuccessfullyCreateFile(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	f := entity.File{
		DisplayName: "new file",
		URL:         "http://dfmjsfo",
		ContentType: "image/png",
		CardID:      uint(1),
	}

	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'files' ('created_at','updated_at','display_name','key','url','content_type','card_id')
		VALUES (?,?,?,?,?,?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(createdAt, updatedAt, f.DisplayName, f.Key, f.URL, f.ContentType, f.CardID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Create(&f); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, f.ID, uint(1))
}

func TestShouldSuccessfullyFindFile(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	userID := uint(1)
	fileID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'files'.*
		FROM 'files'
		Join cards ON files.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.user_id = ?) AND ('files'.'id' = %d) ORDER BY 'files'.'id' ASC LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, fileID))).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(fileID))

	f, err := r.Find(fileID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, f.ID, fileID)
}

func TestShouldNotFindFile(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	userID := uint(1)
	fileID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'files'.*
		FROM 'files'
		Join cards ON files.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.user_id = ?) AND ('files'.'id' = %d) ORDER BY 'files'.'id' ASC LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, fileID))).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(fileID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorRecordNotFound)
}

func TestShouldSuccessfullyDeleteFile(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	f := &entity.File{
		ID: uint(1),
	}

	query := "DELETE FROM `files`  WHERE `files`.`id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(f.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Delete(f); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldNotDeleteFile(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	f := &entity.File{
		ID: uint(1),
	}

	query := "DELETE FROM `files` WHERE `files`.`id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(f.ID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err := r.Delete(f)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidRequest)
}

func TestShouldSuccessfullyGetAllFiles(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	userID := uint(1)
	boardID := uint(2)

	mockFile := entity.File{
		ID:          uint(3),
		DisplayName: "file1",
		URL:         "http://test",
		ContentType: "image/png",
		CardID:      uint(4),
	}

	query := utils.ReplaceQuotationForQuery(`
		SELECT files.id, files.display_name, files.url, files.content_type, files.card_id
		FROM 'files'
		Join cards ON files.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.id = ?) AND (boards.user_id = ?)`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(boardID, userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "display_name", "url", "content_type", "card_id"}).
				AddRow(mockFile.ID, mockFile.DisplayName, mockFile.URL, mockFile.ContentType, mockFile.CardID))

	fs := r.GetAll(boardID, userID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, (*fs)[0].ID, mockFile.ID)
	assert.Equal(t, (*fs)[0].DisplayName, mockFile.DisplayName)
	assert.Equal(t, (*fs)[0].URL, mockFile.URL)
	assert.Equal(t, (*fs)[0].ContentType, mockFile.ContentType)
	assert.Equal(t, (*fs)[0].CardID, mockFile.CardID)
}

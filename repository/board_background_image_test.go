package repository

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"local.packages/utils"
	"local.packages/validator"
)

func TestShouldSuccessfullyValidateUIDOnBoardBackgroundImageRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardBackgroundImageRepository(db)

	userID := uint(1)
	boardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, boardID))).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(boardID))

	if err := r.ValidateUID(boardID, userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureValidateUIDOnBoardBackgroundImageRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardBackgroundImageRepository(db)

	userID := uint(1)
	boardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, boardID))).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	err := r.ValidateUID(boardID, userID)

	if err == nil {
		t.Errorf("was expected an error but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidSession)
}

func TestShouldSuccessfullyCreateBoardBackgroundImage(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardBackgroundImageRepository(db)

	boardID := uint(1)
	backgroundImageID := uint(2)

	query := "INSERT INTO `board_background_images` (`board_id`,`background_image_id`) VALUES (?,?)"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(boardID, backgroundImageID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	b, err := r.Create(boardID, backgroundImageID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, b.BoardID, boardID)
	assert.Equal(t, b.BackgroundImageID, backgroundImageID)
}

func TestShouldFailureCreateBoardBackgroundImageWhenDuplicatePrimaryKey(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardBackgroundImageRepository(db)

	boardID := uint(1)
	backgroundImageID := uint(2)

	query := "INSERT INTO `board_background_images` (`board_id`,`background_image_id`) VALUES (?,?)"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(boardID, backgroundImageID).
		WillReturnError(errors.New("Error 1062: Duplicate entry"))

	mock.ExpectRollback()

	_, err := r.Create(boardID, backgroundImageID)

	if err == nil {
		t.Errorf("was expected an error but did not recieved it")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, validator.ErrorAlreadyBeenTaken)
}

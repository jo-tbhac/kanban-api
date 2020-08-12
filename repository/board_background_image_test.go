package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"local.packages/entity"
	"local.packages/utils"
	"local.packages/validator"
)

func TestShouldSuccessfullyFindBoardBackgroundImageRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardBackgroundImageRepository(db)

	userID := uint(1)
	boardID := uint(2)
	backgroundImageID := uint(3)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'board_background_images'.*
		FROM 'board_background_images'
		Join boards ON board_background_images.board_id = boards.id
		WHERE (boards.user_id = ?) AND (boards.id = ?)
		ORDER BY 'board_background_images'.'board_id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, boardID).
		WillReturnRows(
			sqlmock.NewRows([]string{"board_id", "background_image_id"}).
				AddRow(boardID, backgroundImageID))

	b, err := r.Find(boardID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, b.BoardID, boardID)
	assert.Equal(t, b.BackgroundImageID, backgroundImageID)
}

func TestShouldNotFindBoardBackgroundImageRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardBackgroundImageRepository(db)

	userID := uint(1)
	boardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'board_background_images'.*
		FROM 'board_background_images'
		Join boards ON board_background_images.board_id = boards.id
		WHERE (boards.user_id = ?) AND (boards.id = ?)
		ORDER BY 'board_background_images'.'board_id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, boardID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(boardID, userID)

	if err == nil {
		t.Errorf("was expected an error but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidSession)
}

func TestShouldSuccessfullyUpdateBoardBackgroundImage(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardBackgroundImageRepository(db)

	newBackgroundImageID := uint(1)

	b := &entity.BoardBackgroundImage{
		BoardID:           uint(2),
		BackgroundImageID: uint(3),
	}

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'board_background_images'
		SET 'background_image_id' = ?
		WHERE 'board_background_images'.'board_id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newBackgroundImageID, b.BoardID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Update(b, newBackgroundImageID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, b.BackgroundImageID, newBackgroundImageID)
}

func TestShouldFailureUpdateBoardBackgroundImageWhenDuplicatePrimaryKey(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardBackgroundImageRepository(db)

	newBackgroundImageID := uint(1)

	b := &entity.BoardBackgroundImage{
		BoardID:           uint(2),
		BackgroundImageID: uint(3),
	}

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'board_background_images'
		SET 'background_image_id' = ?
		WHERE 'board_background_images'.'board_id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newBackgroundImageID, b.BoardID).
		WillReturnError(errors.New("Error 1062: Duplicate entry"))

	mock.ExpectRollback()

	err := r.Update(b, newBackgroundImageID)

	if err == nil {
		t.Errorf("was expected an error but did not recieved it")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, validator.ErrorAlreadyBeenTaken)
}

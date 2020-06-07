package repository

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"local.packages/entity"
	"local.packages/utils"
)

func TestShouldSuccessfullyValidateUIDOnListRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewListRepository(db)

	userID := uint(1)
	boardID := uint(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(boardID, userID))

	if err := r.ValidateUID(boardID, userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureValidateUIDOnListRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewListRepository(db)

	userID := uint(1)
	boardID := uint(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`")).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	err := r.ValidateUID(boardID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyFindList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewListRepository(db)

	userID := uint(1)
	listID := uint(2)
	boardID := uint(3)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `lists`.* FROM `lists` Join boards on boards.id = lists.board_id")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "board_id"}).AddRow(listID, boardID))

	l, err := r.Find(listID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, l.ID, listID)
	assert.Equal(t, l.BoardID, boardID)
}

func TestShouldNotFindList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewListRepository(db)

	userID := uint(1)
	listID := uint(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `lists`.* FROM `lists` Join boards on boards.id = lists.board_id")).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(listID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyCreateList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewListRepository(db)

	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}
	name := strings.Repeat("a", 50)
	boardID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists` (`created_at`,`updated_at`,`deleted_at`,`name`,`board_id`)")).
		WithArgs(createdAt, updatedAt, nil, name, boardID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	l, err := r.Create(name, boardID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, l.ID, uint(1))
	assert.Equal(t, l.Name, name)
	assert.Equal(t, l.BoardID, boardID)
	assert.IsType(t, l.CreatedAt, time.Time{})
	assert.IsType(t, l.UpdatedAt, time.Time{})
	assert.Nil(t, l.DeletedAt)
}

func TestShouldNotCreateList(t *testing.T) {
	type testCase struct {
		testName      string
		listName      string
		boardID       uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a name",
			listName:      "",
			boardID:       uint(1),
			expectedError: "Name must exist",
		}, {
			testName:      "when name size more than 50 characters",
			listName:      strings.Repeat("a", 51),
			boardID:       uint(1),
			expectedError: "Name is too long (maximum is 50 characters)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewListRepository(db)

			mock.ExpectBegin()

			_, err := r.Create(tc.listName, tc.boardID)

			if err == nil {
				t.Error("was expected an error, but did not recieved it.")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, err[0].Text, tc.expectedError)
		})
	}
}

func TestShouldSuccessfullyUpdateList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewListRepository(db)

	l := &entity.List{
		ID:   uint(1),
		Name: "sample list",
	}

	updatedAt := utils.AnyTime{}
	name := strings.Repeat("a", 50)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET")).
		WithArgs(name, updatedAt, l.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Update(l, name); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, l.Name, name)
}

func TestShouldNotUpdateList(t *testing.T) {
	type testCase struct {
		testName      string
		listName      string
		boardID       uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a name",
			listName:      "",
			expectedError: "Name must exist",
		}, {
			testName:      "when name size more than 50 characters",
			listName:      strings.Repeat("a", 51),
			expectedError: "Name is too long (maximum is 50 characters)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewListRepository(db)

			l := &entity.List{
				ID:   uint(1),
				Name: "sample list",
			}

			mock.ExpectBegin()

			err := r.Update(l, tc.listName)

			if err == nil {
				t.Error("was expected an error, but did not recieved it.")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, err[0].Text, tc.expectedError)
		})
	}
}

func TestShouldSuccessfullyDeleteList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewListRepository(db)

	l := &entity.List{
		ID:        uint(1),
		DeletedAt: nil,
	}

	deletedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET")).
		WithArgs(deletedAt, l.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Delete(l); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.IsType(t, l.DeletedAt, &time.Time{})
}

func TestShouldNotDeleteList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewListRepository(db)

	l := &entity.List{
		ID:        uint(1),
		DeletedAt: nil,
	}

	deletedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET")).
		WithArgs(deletedAt, l.ID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err := r.Delete(l)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Nil(t, l.DeletedAt)
	assert.Equal(t, err[0].Text, "invalid request")
}

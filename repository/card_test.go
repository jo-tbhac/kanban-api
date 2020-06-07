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

func TestShouldSuccessfullyValidateUIDOnCardRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardRepository(db)

	listID := uint(1)
	userID := uint(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards` Join lists ON boards.id = lists.board_id")).
		WithArgs(listID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))

	if err := r.ValidateUID(listID, userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureValidateUIDOnCardRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardRepository(db)

	listID := uint(1)
	userID := uint(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards` Join lists ON boards.id = lists.board_id")).
		WithArgs(listID, userID).
		WillReturnError(gorm.ErrRecordNotFound)

	err := r.ValidateUID(listID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyFindCard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardRepository(db)

	userID := uint(1)
	cardID := uint(2)
	title := "sample card"
	description := "sample description"
	listID := uint(3)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `cards`.* FROM `cards` Join lists ON lists.id = cards.list_id Join boards ON boards.id = lists.board_id")).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "list_id"}).
			AddRow(cardID, title, description, listID))

	c, err := r.Find(cardID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, c.ID, cardID)
	assert.Equal(t, c.Title, title)
	assert.Equal(t, c.Description, description)
	assert.Equal(t, c.ListID, listID)
}

func TestShouldNotFindCard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardRepository(db)

	userID := uint(1)
	cardID := uint(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `cards`.* FROM `cards` Join lists ON lists.id = cards.list_id Join boards ON boards.id = lists.board_id")).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(cardID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyCreateCard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardRepository(db)

	title := "sample card"
	listID := uint(1)
	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}
	description := ""

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `cards` (`created_at`,`updated_at`,`deleted_at`,`title`,`description`,`list_id`)")).
		WithArgs(createdAt, updatedAt, nil, title, description, listID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	c, err := r.Create(title, listID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, c.ID, uint(1))
	assert.Equal(t, c.Title, title)
	assert.Equal(t, c.Description, description)
	assert.Equal(t, c.ListID, listID)
}

func TestShouldNotCreateCard(t *testing.T) {
	type testCase struct {
		testName      string
		title         string
		listID        uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a title",
			title:         "",
			listID:        uint(1),
			expectedError: "Title must exist",
		}, {
			testName:      "when name size more than 50 characters",
			title:         strings.Repeat("w", 51),
			listID:        uint(1),
			expectedError: "Title is too long (maximum is 50 characters)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewCardRepository(db)

			mock.ExpectBegin()

			_, err := r.Create(tc.title, tc.listID)

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

func TestShouldSuccessfullyUpdateCardTitle(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardRepository(db)

	title := strings.Repeat("c", 50)
	updatedAt := utils.AnyTime{}
	description := "sample description"
	c := &entity.Card{
		ID:          uint(1),
		Title:       "sample card",
		Description: description,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `cards` SET")).
		WithArgs(title, updatedAt, c.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.UpdateTitle(c, title); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, c.Title, title)
	assert.Equal(t, c.Description, description)
}

func TestShouldNotUpdateCardTitle(t *testing.T) {
	type testCase struct {
		testName      string
		title         string
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a title",
			title:         "",
			expectedError: "Title must exist",
		}, {
			testName:      "when name size more than 50 characters",
			title:         strings.Repeat("w", 51),
			expectedError: "Title is too long (maximum is 50 characters)",
		},
	}

	for _, tc := range testCases {
		db, mock := utils.NewDBMock(t)
		defer db.Close()

		r := NewCardRepository(db)

		c := &entity.Card{
			ID:          uint(1),
			Title:       "sample card",
			Description: "description",
		}

		mock.ExpectBegin()

		err := r.UpdateTitle(c, tc.title)

		if err == nil {
			t.Error("was expected an error, but did not recieved it.")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("there were unfulfilled expectations: %v", err)
		}

		assert.Equal(t, err[0].Text, tc.expectedError)
	}
}

func TestShouldSuccessfullyUpdateCardDescription(t *testing.T) {
	type testCase struct {
		testName    string
		description string
	}

	testCases := []testCase{
		{
			testName:    "when without a description",
			description: "",
		}, {
			testName:    "when with description",
			description: "updated description",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewCardRepository(db)

			title := "sample card"
			updatedAt := utils.AnyTime{}
			c := &entity.Card{
				ID:          uint(1),
				Title:       title,
				Description: "previous description",
			}

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta("UPDATE `cards` SET")).
				WithArgs(tc.description, updatedAt, c.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			if err := r.UpdateDescription(c, tc.description); err != nil {
				t.Errorf("was not expected an error. %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, c.Title, title)
			assert.Equal(t, c.Description, tc.description)
		})
	}
}

func TestShouldSuccessfullyDeleteCard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardRepository(db)

	c := &entity.Card{
		ID:        uint(1),
		DeletedAt: nil,
	}

	deletedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `cards` SET")).
		WithArgs(deletedAt, c.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Delete(c); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.IsType(t, c.DeletedAt, &time.Time{})
}

func TestShouldNotDeleteCard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCardRepository(db)

	c := &entity.Card{
		ID:        uint(1),
		DeletedAt: nil,
	}

	deletedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `cards` SET")).
		WithArgs(deletedAt, c.ID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err := r.Delete(c)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Nil(t, c.DeletedAt)
	assert.Equal(t, err[0].Text, "invalid request")
}

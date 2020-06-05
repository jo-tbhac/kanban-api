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

func TestShouldSuccessfullyFindBoard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	userID := uint(5)

	mockBoard := entity.Board{
		ID:        uint(1),
		UpdatedAt: time.Now(),
		Name:      "mockBoard",
		UserID:    userID,
	}

	mockList := entity.List{
		ID:      uint(2),
		Name:    "mockList",
		BoardID: mockBoard.ID,
	}

	mockCard := entity.Card{
		ID:          uint(3),
		Title:       "mockCard",
		Description: "mockDescription",
		ListID:      mockList.ID,
	}

	mockLabel := entity.Label{
		ID:    uint(4),
		Name:  "mockLabel",
		Color: "#ffffff",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "updated_at", "name", "user_id"}).
				AddRow(mockBoard.ID, mockBoard.UpdatedAt, mockBoard.Name, userID))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT lists.id, lists.name, lists.board_id FROM `lists`")).
		WithArgs(mockBoard.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "board_id"}).
				AddRow(mockList.ID, mockList.Name, mockList.BoardID))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT cards.id, cards.title, cards.description, cards.list_id FROM `cards`")).
		WithArgs(mockList.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "list_id"}).
				AddRow(mockCard.ID, mockCard.Title, mockCard.Description, mockCard.ListID))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT labels.id, labels.name, labels.color, labels.board_id, card_labels.card_id FROM `labels`")).
		WithArgs(mockCard.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "color", "card_id"}).
				AddRow(mockLabel.ID, mockLabel.Name, mockLabel.Color, mockCard.ID))

	b, err := r.Find(mockBoard.ID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, b.ID, mockBoard.ID)
	assert.Equal(t, b.Name, mockBoard.Name)
	assert.Equal(t, b.UpdatedAt, mockBoard.UpdatedAt)
	assert.Equal(t, b.UserID, mockBoard.UserID)

	assert.Equal(t, b.Lists[0].ID, mockList.ID)
	assert.Equal(t, b.Lists[0].Name, mockList.Name)
	assert.Equal(t, b.Lists[0].BoardID, mockList.BoardID)

	assert.Equal(t, b.Lists[0].Cards[0].ID, mockCard.ID)
	assert.Equal(t, b.Lists[0].Cards[0].Title, mockCard.Title)
	assert.Equal(t, b.Lists[0].Cards[0].Description, mockCard.Description)
	assert.Equal(t, b.Lists[0].Cards[0].ListID, mockCard.ListID)

	assert.Equal(t, b.Lists[0].Cards[0].Labels[0].ID, mockLabel.ID)
	assert.Equal(t, b.Lists[0].Cards[0].Labels[0].Name, mockLabel.Name)
	assert.Equal(t, b.Lists[0].Cards[0].Labels[0].Color, mockLabel.Color)
}

func TestShouldNotFindBoardWhenUserIdIsInvalid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	userID := uint(1)
	boardID := uint(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`")).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(boardID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieve it. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyFindBoardWithoutRelatedModel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	userID := uint(1)
	boardID := uint(2)
	updatedAt := time.Now()
	name := "sampleBoard"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`")).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "updated_at", "name", "user_id"}).
				AddRow(boardID, updatedAt, name, userID))

	b, err := r.FindWithoutPreload(boardID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, b.ID, boardID)
	assert.Equal(t, b.Name, name)
	assert.Equal(t, b.UpdatedAt, updatedAt)
	assert.Equal(t, b.UserID, userID)
	assert.Nil(t, b.Lists)
}

func TestShouldNotFindBoardWithoutRelatedModel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	userID := uint(1)
	boardID := uint(2)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`")).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.FindWithoutPreload(boardID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieve it. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldCreateBoard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}
	name := strings.Repeat("a", 50)
	userID := uint(1)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `boards` (`created_at`,`updated_at`,`deleted_at`,`name`,`user_id`)")).
		WithArgs(createdAt, updatedAt, nil, name, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	b, err := r.Create(name, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, b.ID, uint(1))
	assert.Equal(t, b.Name, name)
	assert.Equal(t, b.UserID, userID)
}

func TestShouldNotCreateBoard(t *testing.T) {
	type testCase struct {
		testName      string
		boardName     string
		userID        uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a name",
			boardName:     "",
			userID:        uint(1),
			expectedError: "Name must exist",
		}, {
			testName:      "when name size more than 50 characters",
			boardName:     strings.Repeat("a", 51),
			userID:        uint(1),
			expectedError: "Name is too long (maximum is 50 characters)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewBoardRepository(db)

			mock.ExpectBegin()

			_, err := r.Create(tc.boardName, tc.userID)

			if err == nil {
				t.Errorf("was expected an error, but did not recieve it. %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, err[0].Text, tc.expectedError)
		})
	}
}

func TestShouldUpdateBoard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	b := &entity.Board{
		ID:   uint(1),
		Name: "sample_board",
	}

	name := strings.Repeat("b", 50)
	updatedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `boards` SET")).
		WithArgs(name, updatedAt, b.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Update(b, name); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, b.Name, name)
}

func TestShouldNotUpdateBoard(t *testing.T) {
	type testCase struct {
		testName      string
		boardName     string
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a name",
			boardName:     "",
			expectedError: "Name must exist",
		}, {
			testName:      "when name size more than 50 characters",
			boardName:     strings.Repeat("a", 51),
			expectedError: "Name is too long (maximum is 50 characters)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewBoardRepository(db)
			b := &entity.Board{
				ID:   uint(1),
				Name: "sample_board",
			}

			mock.ExpectBegin()

			err := r.Update(b, tc.boardName)

			if err == nil {
				t.Errorf("was expected an error, but did not recieve it. %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, err[0].Text, tc.expectedError)
		})
	}
}

func TestShouldSuccessfullyDeleteBoard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	boardID := uint(1)
	userID := uint(2)
	deletedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `boards` SET")).
		WithArgs(deletedAt, boardID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Delete(boardID, userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldNotDeleteBoard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	boardID := uint(1)
	userID := uint(2)
	deletedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `boards` SET")).
		WithArgs(deletedAt, boardID, userID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err := r.Delete(boardID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid request")
}
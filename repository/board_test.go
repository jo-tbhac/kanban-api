package repository

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"local.packages/entity"
	"local.packages/utils"
	"local.packages/validator"
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

	mockBackgroundImage := entity.BoardBackgroundImage{
		BoardID:           mockBoard.ID,
		BackgroundImageID: uint(2),
	}

	mockList := entity.List{
		ID:      uint(2),
		Name:    "mockList",
		BoardID: mockBoard.ID,
		Index:   1,
	}

	mockCard := entity.Card{
		ID:          uint(3),
		Title:       "mockCard",
		Description: "mockDescription",
		ListID:      mockList.ID,
		Index:       1,
	}

	mockLabel := entity.Label{
		ID:    uint(4),
		Name:  "mockLabel",
		Color: "#ffffff",
	}

	mockCover := &entity.Cover{
		CardID: mockCard.ID,
		FileID: uint(5),
	}

	boardQuery := utils.ReplaceQuotationForQuery(`
		SELECT id, updated_at, name, user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(boardQuery, mockBoard.ID))).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "updated_at", "name", "user_id"}).
				AddRow(mockBoard.ID, mockBoard.UpdatedAt, mockBoard.Name, userID))

	backgroundImageQuery := utils.ReplaceQuotationForQuery(`
		SELECT *
		FROM 'board_background_images'
		WHERE ('board_id' IN (?))
		ORDER BY 'board_background_images'.'board_id' ASC`)

	mock.ExpectQuery(regexp.QuoteMeta(backgroundImageQuery)).
		WithArgs(mockBackgroundImage.BoardID).
		WillReturnRows(
			sqlmock.NewRows([]string{"board_id", "background_image_id"}).
				AddRow(mockBackgroundImage.BoardID, mockBackgroundImage.BackgroundImageID))

	listQuery := utils.ReplaceQuotationForQuery(`
		SELECT lists.id, lists.name, lists.board_id, lists.index
		FROM 'lists'
		WHERE 'lists'.'deleted_at' IS NULL AND (('board_id' IN (?)))
		ORDER BY lists.index asc,'lists'.'id' ASC`)

	mock.ExpectQuery(regexp.QuoteMeta(listQuery)).
		WithArgs(mockBoard.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "board_id", "index"}).
				AddRow(mockList.ID, mockList.Name, mockList.BoardID, mockList.Index))

	cardQuery := utils.ReplaceQuotationForQuery(`
		SELECT cards.id, cards.title, cards.description, cards.list_id, cards.index
		FROM 'cards'
		WHERE 'cards'.'deleted_at' IS NULL AND (('list_id' IN (?)))
		ORDER BY cards.index asc,'cards'.'id' ASC`)

	mock.ExpectQuery(regexp.QuoteMeta(cardQuery)).
		WithArgs(mockList.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "description", "list_id", "index"}).
				AddRow(mockCard.ID, mockCard.Title, mockCard.Description, mockCard.ListID, mockCard.Index))

	labelQuery := utils.ReplaceQuotationForQuery(`
		SELECT labels.id, labels.name, labels.color, labels.board_id, card_labels.card_id
		FROM 'labels'
		INNER JOIN 'card_labels' ON 'card_labels'.'label_id' = 'labels'.'id'
		WHERE 'labels'.'deleted_at' IS NULL AND (('card_labels'.'card_id' IN (?)))`)

	mock.ExpectQuery(regexp.QuoteMeta(labelQuery)).
		WithArgs(mockCard.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "color", "card_id"}).
				AddRow(mockLabel.ID, mockLabel.Name, mockLabel.Color, mockCard.ID))
	coverQuery := utils.ReplaceQuotationForQuery(`
		SELECT *
		FROM 'covers'
		WHERE ('card_id' IN (?))
		ORDER BY 'covers'.'card_id' ASC`)

	mock.ExpectQuery(regexp.QuoteMeta(coverQuery)).
		WithArgs(mockCard.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"card_id", "file_id"}).
				AddRow(mockCover.CardID, mockCover.FileID))

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

	assert.Equal(t, b.BackgroundImage.BoardID, mockBackgroundImage.BoardID)
	assert.Equal(t, b.BackgroundImage.BackgroundImageID, mockBackgroundImage.BackgroundImageID)

	assert.Equal(t, b.Lists[0].ID, mockList.ID)
	assert.Equal(t, b.Lists[0].Name, mockList.Name)
	assert.Equal(t, b.Lists[0].BoardID, mockList.BoardID)
	assert.Equal(t, b.Lists[0].Index, mockList.Index)

	assert.Equal(t, b.Lists[0].Cards[0].ID, mockCard.ID)
	assert.Equal(t, b.Lists[0].Cards[0].Title, mockCard.Title)
	assert.Equal(t, b.Lists[0].Cards[0].Description, mockCard.Description)
	assert.Equal(t, b.Lists[0].Cards[0].ListID, mockCard.ListID)

	assert.Equal(t, b.Lists[0].Cards[0].Labels[0].ID, mockLabel.ID)
	assert.Equal(t, b.Lists[0].Cards[0].Labels[0].Name, mockLabel.Name)
	assert.Equal(t, b.Lists[0].Cards[0].Labels[0].Color, mockLabel.Color)

	assert.Equal(t, b.Lists[0].Cards[0].Cover.CardID, mockCover.CardID)
	assert.Equal(t, b.Lists[0].Cards[0].Cover.FileID, mockCover.FileID)
}

func TestShouldNotFindBoardWhenUserIdIsInvalid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	userID := uint(1)
	boardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT id, updated_at, name, user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, boardID))).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(boardID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieve it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorRecordNotFound)
}

func TestShouldSuccessfullyFindBoardWithoutRelatedModel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	userID := uint(1)
	boardID := uint(2)
	updatedAt := time.Now()
	name := "sampleBoard"

	query := utils.ReplaceQuotationForQuery(`
		SELECT id, updated_at, name, user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, boardID))).
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

	query := utils.ReplaceQuotationForQuery(`
		SELECT id, updated_at, name, user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, boardID))).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.FindWithoutPreload(boardID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieve it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorRecordNotFound)
}

func TestShouldCreateBoard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}
	name := strings.Repeat("a", 50)
	userID := uint(1)
	backgroundImageID := uint(2)

	insertBoardQuery := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'boards' ('created_at','updated_at','deleted_at','name','user_id')
		VALUES (?,?,?,?,?)`)

	insertBackgroundImageQuery := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'board_background_images' ('board_id','background_image_id')
		VALUES (?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(insertBoardQuery)).
		WithArgs(createdAt, updatedAt, nil, name, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta(insertBackgroundImageQuery)).
		WithArgs(1, backgroundImageID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	b, err := r.Create(name, backgroundImageID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, b.ID, uint(1))
	assert.Equal(t, b.Name, name)
	assert.Equal(t, b.UserID, userID)
	assert.Equal(t, b.BackgroundImage.BoardID, uint(1))
	assert.Equal(t, b.BackgroundImage.BackgroundImageID, backgroundImageID)
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
			expectedError: validator.ErrorRequired("ボード名"),
		}, {
			testName:      "when name size more than 50 characters",
			boardName:     strings.Repeat("a", 51),
			userID:        uint(1),
			expectedError: validator.ErrorTooLong("ボード名", "50"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewBoardRepository(db)

			mock.ExpectBegin()
			mock.ExpectRollback()

			backgroundImageID := uint(2)

			_, err := r.Create(tc.boardName, backgroundImageID, tc.userID)

			if err == nil {
				t.Error("was expected an error, but did not recieve it.")
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

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'boards'
		SET 'name' = ?, 'updated_at' = ?
		WHERE 'boards'.'deleted_at' IS NULL AND 'boards'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
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
			expectedError: validator.ErrorRequired("ボード名"),
		}, {
			testName:      "when name size more than 50 characters",
			boardName:     strings.Repeat("a", 51),
			expectedError: validator.ErrorTooLong("ボード名", "50"),
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
				t.Error("was expected an error, but did not recieve it.")
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

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'boards'
		SET 'deleted_at'=?
		WHERE 'boards'.'deleted_at' IS NULL AND ((id = ? AND user_id = ?))`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
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

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'boards'
		SET 'deleted_at'=?
		WHERE 'boards'.'deleted_at' IS NULL AND ((id = ? AND user_id = ?))`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(deletedAt, boardID, userID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err := r.Delete(boardID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidRequest)
}

func TestShouldSuccessfullySearchBoard(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBoardRepository(db)

	userID := uint(1)
	boardID := uint(2)
	name := "sample"

	query := utils.ReplaceQuotationForQuery(`
		SELECT id FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND (name LIKE ?))`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID, "%"+name+"%").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(boardID))

	ids := r.Search(name, userID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	for _, id := range ids {
		assert.Equal(t, id, boardID)
	}
}

func TestShouldReturnsAllBoardInstances(t *testing.T) {
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

	mockBackgroundImage := entity.BoardBackgroundImage{
		BoardID:           mockBoard.ID,
		BackgroundImageID: uint(2),
	}

	boardQuery := utils.ReplaceQuotationForQuery(`
		SELECT id, updated_at, name, user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?))`)

	backgroundImageQuery := utils.ReplaceQuotationForQuery(`
		SELECT *
		FROM 'board_background_images'
		WHERE ('board_id' IN (?))`)

	mock.ExpectQuery(regexp.QuoteMeta(boardQuery)).
		WithArgs(userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "updated_at", "name", "user_id"}).
				AddRow(mockBoard.ID, mockBoard.UpdatedAt, mockBoard.Name, mockBoard.UserID))

	mock.ExpectQuery(regexp.QuoteMeta(backgroundImageQuery)).
		WithArgs(mockBackgroundImage.BoardID).
		WillReturnRows(
			sqlmock.NewRows([]string{"board_id", "background_image_id"}).
				AddRow(mockBackgroundImage.BoardID, mockBackgroundImage.BackgroundImageID))

	bs := r.GetAll(userID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, (*bs)[0].ID, mockBoard.ID)
	assert.Equal(t, (*bs)[0].UpdatedAt, mockBoard.UpdatedAt)
	assert.Equal(t, (*bs)[0].Name, mockBoard.Name)
	assert.Equal(t, (*bs)[0].UserID, mockBoard.UserID)

	assert.Equal(t, (*bs)[0].BackgroundImage.BoardID, mockBackgroundImage.BoardID)
	assert.Equal(t, (*bs)[0].BackgroundImage.BackgroundImageID, mockBackgroundImage.BackgroundImageID)
}

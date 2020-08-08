package repository

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"local.packages/entity"
	"local.packages/utils"
)

func TestShouldSuccessfullyValidateUIDOnCheckListRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

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

func TestShouldFailureValidateUIDOnCheckListRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

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

func TestShouldSuccessfullyFindCheckList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

	userID := uint(1)
	checkListID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'check_lists'.*
		FROM 'check_lists'
		Join cards ON check_lists.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.user_id = ?) AND ('check_lists'.'id' = %d)
		ORDER BY 'check_lists'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, checkListID))).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(checkListID))

	cl, err := r.Find(checkListID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, cl.ID, checkListID)
}

func TestShouldNotFindCheckList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

	userID := uint(1)
	checkListID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'check_lists'.*
		FROM 'check_lists'
		Join cards ON check_lists.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.user_id = ?) AND ('check_lists'.'id' = %d)
		ORDER BY 'check_lists'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, checkListID))).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(checkListID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorRecordNotFound)
}

func TestShouldSuccessfullyCreateCheckList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}
	title := "new check list"
	cardID := uint(1)

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'check_lists' ('created_at','updated_at','title','card_id')
		VALUES (?,?,?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(createdAt, updatedAt, title, cardID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	cl, err := r.Create(title, cardID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, cl.ID, uint(1))
	assert.Equal(t, cl.Title, title)
	assert.Equal(t, cl.CardID, cardID)
}

func TestShouldNotCreateCheckList(t *testing.T) {
	type testCase struct {
		testName      string
		title         string
		cardID        uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a title",
			title:         "",
			cardID:        uint(1),
			expectedError: "Title must exist",
		}, {
			testName:      "when name size more than 50 characters",
			title:         strings.Repeat("w", 51),
			cardID:        uint(1),
			expectedError: "Title is too long (maximum is 50 characters)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewCheckListRepository(db)

			mock.ExpectBegin()

			cl, err := r.Create(tc.title, tc.cardID)

			if err == nil {
				t.Errorf("was not expected an error. %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, cl.ID, uint(0))
			assert.Equal(t, err[0].Text, tc.expectedError)
		})
	}
}

func TestShouldSuccessfullyUpdateCheckList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

	cl := &entity.CheckList{
		ID:    uint(1),
		Title: "previous title",
	}

	updatedAt := utils.AnyTime{}
	title := "update title"

	query := "UPDATE `check_lists` SET `title` = ?, `updated_at` = ? WHERE `check_lists`.`id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(title, updatedAt, cl.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Update(cl, title); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, cl.Title, title)
}

func TestShouldNotUpdateCheckList(t *testing.T) {
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
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewCheckListRepository(db)

			cl := &entity.CheckList{
				ID:    uint(1),
				Title: "previous title",
			}

			mock.ExpectBegin()
			mock.ExpectRollback()

			err := r.Update(cl, tc.title)
			if err == nil {
				t.Errorf("was expected an error, but did not recieved it.")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, err[0].Text, tc.expectedError)
		})
	}
}

func TestShouldSuccessfullyDeleteCheckList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

	cl := &entity.CheckList{
		ID: uint(1),
	}

	query := "DELETE FROM `check_lists` WHERE `check_lists`.`id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(cl.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Delete(cl); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureDeleteCheckList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

	cl := &entity.CheckList{
		ID: uint(1),
	}

	query := "DELETE FROM `check_lists` WHERE `check_lists`.`id` = ?"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(cl.ID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err := r.Delete(cl)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidRequest)
}

func TestShouldSuccessfullyGetAllCheckList(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListRepository(db)

	boardID := uint(1)
	userID := uint(2)

	mockCheckList := entity.CheckList{
		ID:     uint(3),
		Title:  "mockCheckList",
		CardID: uint(4),
	}

	mockCheckListItem := entity.CheckListItem{
		ID:          uint(5),
		Name:        "mockCheckListItem",
		Check:       false,
		CheckListID: mockCheckList.ID,
	}

	checkListQuery := utils.ReplaceQuotationForQuery(`
		SELECT check_lists.id, check_lists.title, check_lists.card_id
		FROM 'check_lists'
		Join cards ON check_lists.card_id = cards.id
		Join lists ON cards.list_id = lists.id Join boards ON lists.board_id = boards.id
		WHERE (boards.id = ?) AND (boards.user_id = ?)`)

	mock.ExpectQuery(regexp.QuoteMeta(checkListQuery)).
		WithArgs(boardID, userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "card_id"}).
				AddRow(mockCheckList.ID, mockCheckList.Title, mockCheckList.CardID))

	checkListItemQuery := utils.ReplaceQuotationForQuery(`
		SELECT check_list_items.id, check_list_items.name, check_list_items.check_list_id, check_list_items.check
		FROM 'check_list_items'
		WHERE ('check_list_id' IN (?))`)

	mock.ExpectQuery(regexp.QuoteMeta(checkListItemQuery)).
		WithArgs(mockCheckListItem.CheckListID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "check", "check_list_id"}).
				AddRow(mockCheckListItem.ID, mockCheckListItem.Name, mockCheckListItem.Check, mockCheckListItem.CheckListID))

	cs := r.GetAll(boardID, userID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, (*cs)[0].ID, mockCheckList.ID)
	assert.Equal(t, (*cs)[0].Title, mockCheckList.Title)
	assert.Equal(t, (*cs)[0].CardID, mockCheckList.CardID)
	assert.Equal(t, (*cs)[0].Items[0].ID, mockCheckListItem.ID)
	assert.Equal(t, (*cs)[0].Items[0].Name, mockCheckListItem.Name)
	assert.Equal(t, (*cs)[0].Items[0].Check, mockCheckListItem.Check)
	assert.Equal(t, (*cs)[0].Items[0].CheckListID, mockCheckListItem.CheckListID)
}

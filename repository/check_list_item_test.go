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
	"local.packages/validator"
)

func TestShouldSuccessfullyValidateUIDOnCheckListItemRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	userID := uint(1)
	checkListID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		Join lists ON boards.id = lists.board_id
		Join cards ON lists.id = cards.list_id
		Join check_lists ON cards.id = check_lists.card_id
		WHERE 'boards'.'deleted_at' IS NULL AND ((check_lists.id = ?) AND (boards.user_id = ?))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(checkListID, userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userID))

	if err := r.ValidateUID(checkListID, userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureValidateUIDOnCheckListItemRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	userID := uint(1)
	checkListID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		Join lists ON boards.id = lists.board_id
		Join cards ON lists.id = cards.list_id
		Join check_lists ON cards.id = check_lists.card_id
		WHERE 'boards'.'deleted_at' IS NULL AND ((check_lists.id = ?) AND (boards.user_id = ?))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(checkListID, userID).
		WillReturnError(gorm.ErrRecordNotFound)

	err := r.ValidateUID(checkListID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidSession)
}

func TestShouldSuccessfullyFindCheckListItem(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	userID := uint(1)
	checkListItemID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'check_list_items'.*
		FROM 'check_list_items'
		Join check_lists ON check_list_items.check_list_id = check_lists.id
		Join cards ON check_lists.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.user_id = ?) AND ('check_list_items'.'id' = %d)
		ORDER BY 'check_list_items'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, checkListItemID))).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(checkListItemID))

	item, err := r.Find(checkListItemID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, item.ID, checkListItemID)
}

func TestShouldNotFindCheckListItem(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	userID := uint(1)
	checkListItemID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'check_list_items'.*
		FROM 'check_list_items'
		Join check_lists ON check_list_items.check_list_id = check_lists.id
		Join cards ON check_lists.card_id = cards.id
		Join lists ON cards.list_id = lists.id
		Join boards ON lists.board_id = boards.id
		WHERE (boards.user_id = ?) AND ('check_list_items'.'id' = %d)
		ORDER BY 'check_list_items'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, checkListItemID))).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(checkListItemID, userID)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorRecordNotFound)
}

func TestShouldSuccessfullyCreateCheckListItem(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}
	name := "new item"
	checkListID := uint(1)

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'check_list_items' ('created_at','updated_at','name','check_list_id')
		VALUES (?,?,?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(createdAt, updatedAt, name, checkListID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	item, err := r.Create(name, checkListID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, item.ID, uint(1))
	assert.Equal(t, item.Name, name)
	assert.Equal(t, item.CheckListID, checkListID)
}

func TestShouldNotCreateCheckListItem(t *testing.T) {
	type testCase struct {
		testName      string
		name          string
		checkListID   uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a title",
			name:          "",
			checkListID:   uint(1),
			expectedError: validator.ErrorRequired("アイテム名"),
		}, {
			testName:      "when name size more than 50 characters",
			name:          strings.Repeat("w", 51),
			checkListID:   uint(1),
			expectedError: validator.ErrorTooLong("アイテム名", "50"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewCheckListItemRepository(db)

			mock.ExpectBegin()

			item, err := r.Create(tc.name, tc.checkListID)

			if err == nil {
				t.Errorf("was not expected an error. %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, item.ID, uint(0))
			assert.Equal(t, err[0].Text, tc.expectedError)
		})
	}
}

func TestShouldSuccessfullyUpdateCheckListItemName(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	item := &entity.CheckListItem{
		ID:   uint(1),
		Name: "previous name",
	}

	updatedAt := utils.AnyTime{}
	name := "update item"

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'check_list_items'
		SET 'name' = ?, 'updated_at' = ?
		WHERE 'check_list_items'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(name, updatedAt, item.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Update(item, name); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, item.Name, name)
}

func TestShouldNotUpdateCheckListItemName(t *testing.T) {
	type testCase struct {
		testName      string
		name          string
		checkListID   uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a title",
			name:          "",
			checkListID:   uint(1),
			expectedError: validator.ErrorRequired("アイテム名"),
		}, {
			testName:      "when name size more than 50 characters",
			name:          strings.Repeat("w", 51),
			checkListID:   uint(1),
			expectedError: validator.ErrorTooLong("アイテム名", "50"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewCheckListItemRepository(db)

			item := &entity.CheckListItem{
				ID:   uint(1),
				Name: "previous name",
			}

			mock.ExpectBegin()

			err := r.Update(item, tc.name)

			if err == nil {
				t.Errorf("was not expected an error. %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, err[0].Text, tc.expectedError)
		})
	}
}

func TestShouldSuccessfullyUpdateCheckListItemsCheck(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	item := &entity.CheckListItem{
		ID:    uint(1),
		Name:  "previous name",
		Check: false,
	}

	updatedAt := utils.AnyTime{}
	check := true

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'check_list_items'
		SET 'check' = ?, 'updated_at' = ?
		WHERE 'check_list_items'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(check, updatedAt, item.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Check(item, check); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, item.Check, check)
}

func TestShouldSuccessfullyDeleteCheckListItem(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	item := &entity.CheckListItem{
		ID: uint(1),
	}

	query := utils.ReplaceQuotationForQuery(`
		DELETE FROM 'check_list_items'
		WHERE 'check_list_items'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(item.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.Delete(item); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldNotDeleteCheckListItem(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewCheckListItemRepository(db)

	item := &entity.CheckListItem{
		ID: uint(1),
	}

	query := utils.ReplaceQuotationForQuery(`
		DELETE FROM 'check_list_items'
		WHERE 'check_list_items'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(item.ID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	err := r.Delete(item)

	if err == nil {
		t.Errorf("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidRequest)
}

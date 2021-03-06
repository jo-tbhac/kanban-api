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

func TestShouldSuccessfullyValidateUIDOnLabelRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewLabelRepository(db)

	userID := uint(1)
	boardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, boardID))).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(boardID, userID))

	if err := r.ValidateUID(boardID, userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureValidateUIDOnLabelRepository(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewLabelRepository(db)

	userID := uint(1)
	boardID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, boardID))).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	err := r.ValidateUID(boardID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidSession)
}

func TestShouldSuccessfullyFindLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewLabelRepository(db)

	userID := uint(1)
	labelID := uint(2)
	boardID := uint(3)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'labels'.*
		FROM 'labels'
		Join boards on boards.id = labels.board_id
		WHERE 'labels'.'deleted_at' IS NULL AND ((boards.user_id = ?) AND ('labels'.'id' = %d))
		ORDER BY 'labels'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, labelID))).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "board_id"}).AddRow(labelID, boardID))

	l, err := r.Find(labelID, userID)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, l.ID, labelID)
	assert.Equal(t, l.BoardID, boardID)
}

func TestShouldNotFindLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewLabelRepository(db)

	userID := uint(1)
	labelID := uint(2)

	query := utils.ReplaceQuotationForQuery(`
		SELECT 'labels'.*
		FROM 'labels'
		Join boards on boards.id = labels.board_id
		WHERE 'labels'.'deleted_at' IS NULL AND ((boards.user_id = ?) AND ('labels'.'id' = %d))
		ORDER BY 'labels'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(query, labelID))).
		WithArgs(userID).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.Find(labelID, userID)

	if err == nil {
		t.Error("was expected an error, but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorRecordNotFound)
}

func TestShouldSuccessfullyCreateLabel(t *testing.T) {
	type testCase struct {
		testName  string
		labelName string
		color     string
		boardID   uint
	}

	testCases := []testCase{
		{
			testName:  "when it valid with all field(color code contains 3 digit).",
			labelName: "sample label",
			color:     "#fff",
			boardID:   uint(1),
		}, {
			testName:  "when it valid with all field(color code contains 6 digit).",
			labelName: "sample label",
			color:     "#c3c3c3",
			boardID:   uint(1),
		}, {
			testName:  "when it valid with all field(name size is 50).",
			labelName: strings.Repeat("a", 50),
			color:     "#eeeeee",
			boardID:   uint(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewLabelRepository(db)

			createdAt := utils.AnyTime{}
			updatedAt := utils.AnyTime{}

			query := utils.ReplaceQuotationForQuery(`
				INSERT INTO 'labels' ('created_at','updated_at','deleted_at','name','color','board_id')
				VALUES (?,?,?,?,?,?)`)

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(query)).
				WithArgs(createdAt, updatedAt, nil, tc.labelName, tc.color, tc.boardID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			l, err := r.Create(tc.labelName, tc.color, tc.boardID)

			if err != nil {
				t.Errorf("was not expected an error. %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, l.ID, uint(1))
			assert.Equal(t, l.Name, tc.labelName)
			assert.Equal(t, l.Color, tc.color)
			assert.Equal(t, l.BoardID, tc.boardID)
		})
	}
}

func TestShouldNotCreateLabel(t *testing.T) {
	type testCase struct {
		testName      string
		labelName     string
		color         string
		boardID       uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a name",
			labelName:     "",
			color:         "#fff",
			boardID:       uint(1),
			expectedError: validator.ErrorRequired("ラベル名"),
		}, {
			testName:      "when name size more than 50 characters",
			labelName:     strings.Repeat("a", 51),
			color:         "#fff",
			boardID:       uint(1),
			expectedError: validator.ErrorTooLong("ラベル名", "50"),
		}, {
			testName:      "when without a color",
			labelName:     "sample label",
			color:         "",
			boardID:       uint(1),
			expectedError: validator.ErrorRequired("ラベルカラー"),
		}, {
			testName:      "when color does not contains hashtag(#)",
			labelName:     "sample label",
			color:         "ffffff",
			boardID:       uint(1),
			expectedError: validator.ErrorHexcolor("ラベルカラー"),
		}, {
			testName:      "when color code more than 6 digit",
			labelName:     "sample label",
			color:         fmt.Sprintf("#%s", strings.Repeat("f", 7)),
			boardID:       uint(1),
			expectedError: validator.ErrorHexcolor("ラベルカラー"),
		}, {
			testName:      "when color contains invalid character",
			labelName:     "sample label",
			color:         "#fgfgfg",
			boardID:       uint(1),
			expectedError: validator.ErrorHexcolor("ラベルカラー"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewLabelRepository(db)

			mock.ExpectBegin()

			_, err := r.Create(tc.labelName, tc.color, tc.boardID)

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

func TestShouldSuccessfullyUpdateLabel(t *testing.T) {
	type testCase struct {
		testName  string
		labelName string
		color     string
	}

	testCases := []testCase{
		{
			testName:  "when it valid with all field(color code contains 3 digit).",
			labelName: "sample label",
			color:     "#fff",
		}, {
			testName:  "when it valid with all field(color code contains 6 digit).",
			labelName: "sample label",
			color:     "#c3c3c3",
		}, {
			testName:  "when it valid with all field(name size is 50).",
			labelName: strings.Repeat("a", 50),
			color:     "#eeeeee",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewLabelRepository(db)

			updatedAt := utils.AnyTime{}

			l := &entity.Label{
				ID:    uint(1),
				Name:  "previous name",
				Color: "#000",
			}

			query := utils.ReplaceQuotationForQuery(`
				UPDATE 'labels'
				SET 'color' = ?, 'name' = ?, 'updated_at' = ?
				WHERE 'labels'.'deleted_at' IS NULL AND 'labels'.'id' = ?`)

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(query)).
				WithArgs(tc.color, tc.labelName, updatedAt, l.ID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			if err := r.Update(l, tc.labelName, tc.color); err != nil {
				t.Errorf("was not expected an error. %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, l.Name, tc.labelName)
			assert.Equal(t, l.Color, tc.color)
		})
	}
}

func TestShouldNotUpdateLabel(t *testing.T) {
	type testCase struct {
		testName      string
		labelName     string
		color         string
		boardID       uint
		expectedError string
	}

	testCases := []testCase{
		{
			testName:      "when without a name",
			labelName:     "",
			color:         "#fff",
			expectedError: validator.ErrorRequired("ラベル名"),
		}, {
			testName:      "when name size more than 50 characters",
			labelName:     strings.Repeat("a", 51),
			color:         "#fff",
			expectedError: validator.ErrorTooLong("ラベル名", "50"),
		}, {
			testName:      "when without a color",
			labelName:     "sample label",
			color:         "",
			expectedError: validator.ErrorRequired("ラベルカラー"),
		}, {
			testName:      "when color does not contains hashtag(#)",
			labelName:     "sample label",
			color:         "ffffff",
			expectedError: validator.ErrorHexcolor("ラベルカラー"),
		}, {
			testName:      "when color code more than 6 digit",
			labelName:     "sample label",
			color:         fmt.Sprintf("#%s", strings.Repeat("f", 7)),
			expectedError: validator.ErrorHexcolor("ラベルカラー"),
		}, {
			testName:      "when color contains invalid character",
			labelName:     "sample label",
			color:         "#fgfgfg",
			expectedError: validator.ErrorHexcolor("ラベルカラー"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			r := NewLabelRepository(db)

			mock.ExpectBegin()

			l := &entity.Label{
				ID:    uint(1),
				Name:  "previous name",
				Color: "#000",
			}

			err := r.Update(l, tc.labelName, tc.color)

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

func TestShouldSuccessfullyDeleteLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewLabelRepository(db)

	l := &entity.Label{
		ID:        uint(1),
		DeletedAt: nil,
	}

	deletedAt := utils.AnyTime{}

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'labels'
		SET 'deleted_at'=?
		WHERE 'labels'.'deleted_at' IS NULL AND 'labels'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
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

func TestShouldNotDeleteLabel(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewLabelRepository(db)

	l := &entity.Label{
		ID:        uint(1),
		DeletedAt: nil,
	}

	deletedAt := utils.AnyTime{}

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'labels'
		SET 'deleted_at'=?
		WHERE 'labels'.'deleted_at' IS NULL AND 'labels'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
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
	assert.Equal(t, err[0].Text, ErrorInvalidRequest)
}

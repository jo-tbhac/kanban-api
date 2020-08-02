package repository

import (
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

	assert.Equal(t, err[0].Text, "invalid parameters")
}

func TestShouldSuccessfullyCreateFile(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewFileRepository(db)

	f := entity.File{
		Name:        "new file",
		URL:         "http://dfmjsfo",
		ContentType: "image/png",
		CardID:      uint(1),
	}

	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'files' ('created_at','updated_at','name','url','content_type','card_id')
		VALUES (?,?,?,?,?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(createdAt, updatedAt, f.Name, f.URL, f.ContentType, f.CardID).
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
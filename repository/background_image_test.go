package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"local.packages/utils"
)

func TestShouldReturnBackgroundImageInstances(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewBackgroundImageRepository(db)

	backgroundImageID := uint(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `background_images`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(backgroundImageID))

	bs := r.GetAll()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, (*bs)[0].ID, backgroundImageID)
}

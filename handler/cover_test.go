package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"local.packages/entity"
	"local.packages/repository"
	"local.packages/utils"
	"local.packages/validator"
)

func TestCreateCoverHandlerShouldReturnsStatusCreatedWithCoverData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCoverHandler(repository.NewCoverRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	cardID := uint(1)
	fileID := uint(2)

	url := fmt.Sprintf("/card/%d/cover/%d", cardID, fileID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, url, nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `covers` (`card_id`,`file_id`)")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.POST("/card/:cardID/cover/:fileID", ch.CreateCover)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.Cover{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["cover"].CardID, cardID)
	assert.Equal(t, res["cover"].FileID, fileID)
}

func TestShouldFailureCreateCoverHandlerWhenDuplicatePrimaryKey(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCoverHandler(repository.NewCoverRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	cardID := uint(1)
	fileID := uint(2)

	url := fmt.Sprintf("/card/%d/cover/%d", cardID, fileID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, url, nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `covers` (`card_id`,`file_id`)")).
		WillReturnError(fmt.Errorf("Error 1062: Duplicate entry '%d-%d' for key", cardID, fileID))

	mock.ExpectRollback()

	r.POST("/card/:cardID/cover/:fileID", ch.CreateCover)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, fmt.Sprintf("%d-%d has already been taken", cardID, fileID))
}

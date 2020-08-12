package handler

import (
	"encoding/json"
	"errors"
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

func TestShouldReturnStatusCreatedAsHTTPRepsponse(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardBackgroundImageHandler(repository.NewBoardBackgroundImageRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	boardID := uint(1)
	backgroundImageID := uint(2)

	url := fmt.Sprintf("/board/%d/background_image/%d", boardID, backgroundImageID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, url, nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	findQuery := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC LIMIT 1`)

	insertQuery := "INSERT INTO `board_background_images` (`board_id`,`background_image_id`) VALUES (?,?)"

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(findQuery, boardID))).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(boardID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).
		WithArgs(boardID, backgroundImageID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.POST("/board/:boardID/background_image/:backgroundImageID", bh.CreateBoardBackgroundImage)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.BoardBackgroundImage{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["board_background_image"].BoardID, boardID)
	assert.Equal(t, res["board_background_image"].BackgroundImageID, backgroundImageID)
}

func TestShouldReturnStatusBadRequestAsHTTPRepsponse(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardBackgroundImageHandler(repository.NewBoardBackgroundImageRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	boardID := uint(1)
	backgroundImageID := uint(2)

	url := fmt.Sprintf("/board/%d/background_image/%d", boardID, backgroundImageID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, url, nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	findQuery := utils.ReplaceQuotationForQuery(`
		SELECT user_id
		FROM 'boards'
		WHERE 'boards'.'deleted_at' IS NULL AND ((user_id = ?) AND ('boards'.'id' = %d))
		ORDER BY 'boards'.'id' ASC LIMIT 1`)

	insertQuery := "INSERT INTO `board_background_images` (`board_id`,`background_image_id`) VALUES (?,?)"

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(findQuery, boardID))).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(boardID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).
		WithArgs(boardID, backgroundImageID).
		WillReturnError(errors.New("Error 1062: Duplicate entry"))

	mock.ExpectRollback()

	r.POST("/board/:boardID/background_image/:backgroundImageID", bh.CreateBoardBackgroundImage)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, validator.ErrorAlreadyBeenTaken)
}

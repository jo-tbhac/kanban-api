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
	"local.packages/repository"
	"local.packages/utils"
	"local.packages/validator"
)

func TestShouldReturnStatusOKWhenSucceedUpdateBoardBackgroundImage(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardBackgroundImageHandler(repository.NewBoardBackgroundImageRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	boardID := uint(1)
	previousBackgroundImageID := uint(2)
	newBackgroundImageID := uint(3)

	url := fmt.Sprintf("/board/%d/background_image/%d", boardID, newBackgroundImageID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, url, nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	findQuery := utils.ReplaceQuotationForQuery(`
		SELECT 'board_background_images'.*
		FROM 'board_background_images'
		Join boards ON board_background_images.board_id = boards.id
		WHERE (boards.user_id = ?) AND (boards.id = ?)
		ORDER BY 'board_background_images'.'board_id' ASC
		LIMIT 1`)

	updateQuery := utils.ReplaceQuotationForQuery(`
		UPDATE 'board_background_images'
		SET 'background_image_id' = ?
		WHERE 'board_background_images'.'board_id' = ?`)

	mock.ExpectQuery(regexp.QuoteMeta(findQuery)).
		WillReturnRows(
			sqlmock.NewRows([]string{"board_id", "background_image_id"}).
				AddRow(boardID, previousBackgroundImageID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
		WithArgs(newBackgroundImageID, boardID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.PATCH("/board/:boardID/background_image/:backgroundImageID", bh.UpdateBoardBackgroundImage)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestShouldReturnStatusBadRequestWhenFailedUpdateBoardBackgroundImage(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardBackgroundImageHandler(repository.NewBoardBackgroundImageRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	boardID := uint(1)
	previousBackgroundImageID := uint(2)
	newBackgroundImageID := uint(3)

	url := fmt.Sprintf("/board/%d/background_image/%d", boardID, newBackgroundImageID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, url, nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	findQuery := utils.ReplaceQuotationForQuery(`
		SELECT 'board_background_images'.*
		FROM 'board_background_images'
		Join boards ON board_background_images.board_id = boards.id
		WHERE (boards.user_id = ?) AND (boards.id = ?)
		ORDER BY 'board_background_images'.'board_id' ASC
		LIMIT 1`)

	updateQuery := utils.ReplaceQuotationForQuery(`
		UPDATE 'board_background_images'
		SET 'background_image_id' = ?
		WHERE 'board_background_images'.'board_id' = ?`)

	mock.ExpectQuery(regexp.QuoteMeta(findQuery)).
		WillReturnRows(
			sqlmock.NewRows([]string{"board_id", "background_image_id"}).
				AddRow(boardID, previousBackgroundImageID))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
		WithArgs(newBackgroundImageID, boardID).
		WillReturnError(errors.New("Error 1062: Duplicate entry"))

	mock.ExpectRollback()

	r.PATCH("/board/:boardID/background_image/:backgroundImageID", bh.UpdateBoardBackgroundImage)
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

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"local.packages/entity"
	"local.packages/repository"
	"local.packages/utils"
	"local.packages/validator"
)

type boardRequestBody struct {
	Name string `json:"name"`
}

func TestCreateBoardHandlerShouldReturnsStatusCreatedWithBoardData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	name := "sample board"

	b, err := json.Marshal(boardRequestBody{
		Name: name,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/board", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `boards` (`created_at`,`updated_at`,`deleted_at`,`name`,`user_id`)")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.POST("/board", bh.CreateBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.Board{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["board"].Name, name)
	assert.Equal(t, res["board"].ID, uint(1))
}

func TestShouldFailureCreateBoardHandlerWhenWithoutBoardName(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	b, err := json.Marshal(boardRequestBody{
		Name: "",
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/board", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectBegin()

	r.POST("/board", bh.CreateBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, "Name must exist")
}

func TestUpdateBoardHandlerShouldReturnsStatusOKWithBoardData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	name := "update board"

	b, err := json.Marshal(boardRequestBody{
		Name: name,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/board/1", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `boards` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.PATCH("/board/:boardID", bh.UpdateBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.Board{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["board"].Name, name)
}

func TestShouldFailureUpdateBoardHandler(t *testing.T) {
	type testCase struct {
		testName         string
		expectedStatus   int
		expectedError    string
		queryParameter   string
		boardRequestBody boardRequestBody
	}

	testCases := []testCase{
		{
			testName:       "when with invalid query parameter",
			expectedStatus: 400,
			expectedError:  "boardID must be an integer",
			queryParameter: "eee",
			boardRequestBody: boardRequestBody{
				Name: "sample board",
			},
		}, {
			testName:       "when without name",
			expectedStatus: 400,
			expectedError:  "Name must exist",
			queryParameter: "1",
			boardRequestBody: boardRequestBody{
				Name: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			bh := NewBoardHandler(repository.NewBoardRepository(db))
			uh := NewUserHandler(repository.NewUserRepository(db))

			r := utils.SetUpRouter()

			b, err := json.Marshal(tc.boardRequestBody)

			if err != nil {
				t.Fatalf("fail to marshal json: %v", err)
			}

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/board/%s", tc.queryParameter)
			req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(b))

			utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

			if tc.testName == "when without name" {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`"))
				mock.ExpectBegin()
			}

			r.PATCH("/board/:boardID", bh.UpdateBoard)
			r.ServeHTTP(w, req)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			res := map[string][]validator.ValidationError{}

			if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
				t.Fatalf("fail to unmarshal response body. %v", err)
			}

			assert.Equal(t, w.Code, tc.expectedStatus)
			assert.Equal(t, res["errors"][0].Text, tc.expectedError)
		})
	}
}

func TestIndexBoardHandlerShouldReturnsStatusOKWithBoardData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/boards", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))

	r.GET("/boards", bh.IndexBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]entity.Board{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Len(t, res["boards"], 1)
}

func TestShowBoardHandlerShouldReturnsStatusOKWithBoardData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/board/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))

	r.GET("/board/:boardID", bh.ShowBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.Board{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["board"].ID, uint(1))
}

func TestShowBoardHandlerShouldReturnsStatusBadRequestWhenRecordNotFound(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/board/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, updated_at, name, user_id FROM `boards`")).
		WillReturnError(gorm.ErrRecordNotFound)

	r.GET("/board/:boardID", bh.ShowBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, ErrorInvalidParameter)
}

func TestDeleteBoardHandlerShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/board/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `boards` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.DELETE("/board/:boardID", bh.DeleteBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestDeleteBoardHandlerShouldReturnsStatusBadRequest(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	bh := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/board/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `boards` SET")).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	r.DELETE("/board/:boardID", bh.DeleteBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, repository.ErrorInvalidRequest)
}

func TestSearchBoardHandlerShouldReturnsStatusOKWithBoardIDs(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewBoardHandler(repository.NewBoardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()

	name := "test"
	boardID := uint(2)

	req, _ := http.NewRequest(http.MethodGet, "/boards/search", nil)
	params := req.URL.Query()
	params.Set("name", name)
	req.URL.RawQuery = params.Encode()

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	query := "SELECT id FROM `boards` WHERE `boards`.`deleted_at` IS NULL AND ((user_id = ?) AND (name LIKE ?))"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(boardID))

	r.GET("boards/search", ch.SearchBoard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]uint{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["board_ids"][0], boardID)
}

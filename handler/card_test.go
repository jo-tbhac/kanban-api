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
	"github.com/stretchr/testify/assert"

	"local.packages/entity"
	"local.packages/repository"
	"local.packages/utils"
	"local.packages/validator"
)

type cardRequestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func TestCreateCardHandlerShouldReturnsStatusCreatedWithCardData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCardHandler(repository.NewCardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	title := "sample card"

	b, err := json.Marshal(cardRequestBody{
		Title: title,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/list/1/card", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards` Join lists ON boards.id = lists.board_id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `cards` (`created_at`,`updated_at`,`deleted_at`,`title`,`description`,`list_id`)")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.POST("/list/:listID/card", ch.CreateCard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.Card{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["card"].Title, title)
	assert.Equal(t, res["card"].ID, uint(1))
}

func TestShouldFailureCreateCardHandlerWhenWithoutTitle(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCardHandler(repository.NewCardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	b, err := json.Marshal(cardRequestBody{
		Title: "",
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/list/1/card", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards` Join lists ON boards.id = lists.board_id"))
	mock.ExpectBegin()

	r.POST("/list/:listID/card", ch.CreateCard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, "Title must exist")
}

func TestUpdateCardHandlerShouldReturnsStatusOK(t *testing.T) {
	type testCase struct {
		testName        string
		attribute       string
		cardRequestBody cardRequestBody
	}

	testCases := []testCase{
		{
			testName:  "when with valid card title",
			attribute: "title",
			cardRequestBody: cardRequestBody{
				Title: "sample card",
			},
		}, {
			testName:  "when with valid card description",
			attribute: "description",
			cardRequestBody: cardRequestBody{
				Description: "sample card description",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			ch := NewCardHandler(repository.NewCardRepository(db))
			uh := NewUserHandler(repository.NewUserRepository(db))

			r := utils.SetUpRouter()

			b, err := json.Marshal(tc.cardRequestBody)

			if err != nil {
				t.Fatalf("fail to marshal json: %v", err)
			}

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/card/1/%s", tc.attribute)
			req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(b))

			utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

			mock.ExpectQuery(regexp.QuoteMeta("SELECT `cards`.* FROM `cards` Join lists ON lists.id = cards.list_id Join boards ON boards.id = lists.board_id")).
				WillReturnRows(sqlmock.NewRows([]string{"title"}).AddRow("sample title"))

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta("UPDATE `cards` SET")).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			r.PATCH("/card/:cardID/:attribute", ch.UpdateCard)
			r.ServeHTTP(w, req)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			res := map[string]entity.Card{}

			if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
				t.Fatalf("fail to unmarshal response body. %v", err)
			}

			assert.Equal(t, w.Code, 200)

			if tc.attribute == "title" {
				assert.Equal(t, res["card"].Title, tc.cardRequestBody.Title)
			} else if tc.attribute == "description" {
				assert.Equal(t, res["card"].Description, tc.cardRequestBody.Description)
			}
		})
	}
}

func TestShouldFailureUpdateCardHandler(t *testing.T) {
	type testCase struct {
		testName        string
		attribute       string
		expectedStatus  int
		expectedError   string
		cardRequestBody cardRequestBody
	}

	testCases := []testCase{
		{
			testName:       "when without a title",
			attribute:      "title",
			expectedStatus: 400,
			expectedError:  "Title must exist",
			cardRequestBody: cardRequestBody{
				Title: "",
			},
		}, {
			testName:        "when with invalid query parameter",
			attribute:       "dfsdfhsksg",
			expectedStatus:  400,
			expectedError:   "invalid parameters",
			cardRequestBody: cardRequestBody{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			ch := NewCardHandler(repository.NewCardRepository(db))
			uh := NewUserHandler(repository.NewUserRepository(db))

			r := utils.SetUpRouter()

			b, err := json.Marshal(tc.cardRequestBody)

			if err != nil {
				t.Fatalf("fail to marshal json: %v", err)
			}

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/card/1/%s", tc.attribute)
			req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(b))

			utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

			mock.ExpectQuery(regexp.QuoteMeta("SELECT `cards`.* FROM `cards` Join lists ON lists.id = cards.list_id Join boards ON boards.id = lists.board_id"))

			if tc.attribute == "title" {
				mock.ExpectBegin()
			}

			r.PATCH("/card/:cardID/:attribute", ch.UpdateCard)
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

func TestDeleteCardHandlerShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCardHandler(repository.NewCardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/card/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `cards`.* FROM `cards` Join lists ON lists.id = cards.list_id Join boards ON boards.id = lists.board_id")).
		WillReturnRows(sqlmock.NewRows([]string{"title"}).AddRow("sample title"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `cards` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.DELETE("/card/:cardID", ch.DeleteCard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestDeleteCardHandlerShouldReturnsStatusBadRequest(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCardHandler(repository.NewCardRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/card/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `cards`.* FROM `cards` Join lists ON lists.id = cards.list_id Join boards ON boards.id = lists.board_id")).
		WillReturnRows(sqlmock.NewRows([]string{"title"}).AddRow("sample title"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `cards` SET")).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	r.DELETE("/card/:cardID", ch.DeleteCard)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, "invalid request")
}
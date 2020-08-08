package handler

import (
	"bytes"
	"encoding/json"
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

type labelRequestBody struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func TestCreateLabelHandlerShouldReturnsStatusCreatedWithLabelData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewLabelHandler(repository.NewLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	name := "sample label"
	color := "#fff"

	b, err := json.Marshal(labelRequestBody{
		Name:  name,
		Color: color,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/board/1/label", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `labels` (`created_at`,`updated_at`,`deleted_at`,`name`,`color`,`board_id`)")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.POST("/board/:boardID/label", lh.CreateLabel)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.Label{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["label"].Color, color)
	assert.Equal(t, res["label"].Name, name)
	assert.Equal(t, res["label"].ID, uint(1))
}

func TestShouldFailureCreateLabelHandler(t *testing.T) {
	type testCase struct {
		testName         string
		expectedStatus   int
		expectedError    string
		labelRequestBody labelRequestBody
	}

	testCases := []testCase{
		{
			testName:       "when without a name",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("ラベル名"),
			labelRequestBody: labelRequestBody{
				Name:  "",
				Color: "#fff",
			},
		}, {
			testName:       "when without a color",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("ラベルカラー"),
			labelRequestBody: labelRequestBody{
				Name:  "sample label",
				Color: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			lh := NewLabelHandler(repository.NewLabelRepository(db))
			uh := NewUserHandler(repository.NewUserRepository(db))

			r := utils.SetUpRouter()

			b, err := json.Marshal(tc.labelRequestBody)

			if err != nil {
				t.Fatalf("fail to marshal json: %v", err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/board/1/label", bytes.NewReader(b))

			utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

			mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))
			mock.ExpectBegin()

			r.POST("/board/:boardID/label", lh.CreateLabel)
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

func TestUpdateLabelHandlerShouldReturnsStatusOKWithLabelData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewLabelHandler(repository.NewLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	name := "sample label"
	color := "#fff"

	b, err := json.Marshal(labelRequestBody{
		Name:  name,
		Color: color,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/label/1", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `labels`.* FROM `labels` Join boards on boards.id = labels.board_id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `labels` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.PATCH("/label/:labelID", lh.UpdateLabel)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.Label{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["label"].Color, color)
	assert.Equal(t, res["label"].Name, name)
}

func TestShouldFailureUpdateLabelHandler(t *testing.T) {
	type testCase struct {
		testName         string
		expectedStatus   int
		expectedError    string
		labelRequestBody labelRequestBody
	}

	testCases := []testCase{
		{
			testName:       "when without a name",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("ラベル名"),
			labelRequestBody: labelRequestBody{
				Name:  "",
				Color: "#fff",
			},
		}, {
			testName:       "when without a color",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("ラベルカラー"),
			labelRequestBody: labelRequestBody{
				Name:  "sample label",
				Color: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			lh := NewLabelHandler(repository.NewLabelRepository(db))
			uh := NewUserHandler(repository.NewUserRepository(db))

			r := utils.SetUpRouter()

			b, err := json.Marshal(tc.labelRequestBody)

			if err != nil {
				t.Fatalf("fail to marshal json: %v", err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, "/label/1", bytes.NewReader(b))

			utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

			mock.ExpectQuery(regexp.QuoteMeta("SELECT `labels`.* FROM `labels` Join boards on boards.id = labels.board_id"))
			mock.ExpectBegin()

			r.PATCH("/label/:labelID", lh.UpdateLabel)
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

func TestIndexLabelHandlerShouldReturnsStatusOKWithLabelData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewLabelHandler(repository.NewLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/board/1/labels", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT labels.id, labels.name, labels.color, labels.board_id FROM `labels` Join boards on boards.id = labels.board_id")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))

	r.GET("/board/:boardID/labels", lh.IndexLabel)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]entity.Label{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Len(t, res["labels"], 1)
}

func TestDeleteLabelHandlerShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewLabelHandler(repository.NewLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/label/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `labels`.* FROM `labels` Join boards on boards.id = labels.board_id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `labels` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.DELETE("/label/:labelID", lh.DeleteLabel)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestDeleteLabelHandlerShouldReturnsStatusBadRequest(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewLabelHandler(repository.NewLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/label/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `labels`.* FROM `labels` Join boards on boards.id = labels.board_id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `labels` SET")).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	r.DELETE("/label/:labelID", lh.DeleteLabel)
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

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

type listRequestBody struct {
	Name string `json:"name"`
}

func TestCreateListHandlerShouldReturnsStatusCreatedWithListData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewListHandler(repository.NewListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	name := "sample list"

	b, err := json.Marshal(listRequestBody{
		Name: name,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/board/1/list", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `lists`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.POST("/board/:boardID/list", lh.CreateList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.List{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["list"].Name, name)
	assert.Equal(t, res["list"].ID, uint(1))
}

func TestShouldFailureCreateListHandlerWhenWithoutName(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewListHandler(repository.NewListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	b, err := json.Marshal(listRequestBody{
		Name: "",
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/board/1/list", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))
	mock.ExpectBegin()

	r.POST("/board/:boardID/list", lh.CreateList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, validator.ErrorRequired("リスト名"))
}

func TestUpdateListShouldReturnsStatusOKWithListData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewListHandler(repository.NewListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	name := "sample list"

	b, err := json.Marshal(listRequestBody{
		Name: name,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/list/1", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `lists`.* FROM `lists` Join boards on boards.id = lists.board_id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.PATCH("/list/:listID", lh.UpdateList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.List{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["list"].Name, name)
}

func TestShouldFailureUpdateListHandlerWhenWithoutName(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewListHandler(repository.NewListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	name := ""

	b, err := json.Marshal(listRequestBody{
		Name: name,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/list/1", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `lists`.* FROM `lists` Join boards on boards.id = lists.board_id"))
	mock.ExpectBegin()

	r.PATCH("/list/:listID", lh.UpdateList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, validator.ErrorRequired("リスト名"))
}

func TestUpdateListIndexShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewListHandler(repository.NewListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	b, err := json.Marshal([]struct {
		ID    uint
		Index int
	}{
		{ID: 1, Index: 1},
		{ID: 2, Index: 3},
		{ID: 3, Index: 2},
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/lists/index", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET `index` = ELT(FIELD(id,1,2,3),1,3,2) WHERE id IN (1,2,3)")).
		WillReturnResult(sqlmock.NewResult(1, 3))

	r.PATCH("/lists/index", lh.UpdateListIndex)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestDeleteListShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewListHandler(repository.NewListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/list/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `lists`.* FROM `lists` Join boards on boards.id = lists.board_id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.DELETE("/list/:listID", lh.DeleteList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestDeleteListHandlerShouldReturnsStatusBadRequest(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	lh := NewListHandler(repository.NewListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/list/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `lists`.* FROM `lists` Join boards on boards.id = lists.board_id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `lists` SET")).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	r.DELETE("/list/:listID", lh.DeleteList)
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

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

type checkListRequestBody struct {
	Title string `json:"title"`
}

func TestCreateCheckListHandlerShouldReturnsStatusCreatedWithCheckListData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListHandler(repository.NewCheckListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	title := "new check list"
	cardID := uint(1)

	b, err := json.Marshal(checkListRequestBody{
		Title: title,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/card/%d/check_list", cardID), bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `check_lists`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.POST("/card/:cardID/check_list", ch.CreateCheckList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.CheckList{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["check_list"].Title, title)
	assert.Equal(t, res["check_list"].CardID, cardID)
	assert.Equal(t, res["check_list"].ID, uint(1))
}

func TestShouldFailureCreateCheckListHandlerWhenRecievedInvalidParameter(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListHandler(repository.NewCheckListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	cardID := uint(1)

	b, err := json.Marshal(checkListRequestBody{
		Title: "",
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/card/%d/check_list", cardID), bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))
	mock.ExpectBegin()

	r.POST("/card/:cardID/check_list", ch.CreateCheckList)
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

func TestUpdateCheckListHandlerShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListHandler(repository.NewCheckListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	title := "new check list"
	checkListID := uint(1)

	b, err := json.Marshal(checkListRequestBody{
		Title: title,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/check_list/%d", checkListID), bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `check_lists`.* FROM `check_lists`"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `check_lists` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.PATCH("/check_list/:checkListID", ch.UpdateCheckList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestShouldFailureUpdateCheckListHandlerWhenRecievedInvalidParameter(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListHandler(repository.NewCheckListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	checkListID := uint(1)

	b, err := json.Marshal(checkListRequestBody{
		Title: "",
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/check_list/%d", checkListID), bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `check_lists`.* FROM `check_lists`"))
	mock.ExpectBegin()

	r.PATCH("/check_list/:checkListID", ch.UpdateCheckList)
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

func TestDeleteCheckListHandlerShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListHandler(repository.NewCheckListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	checkListID := uint(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/check_list/%d", checkListID), nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `check_lists`.* FROM `check_lists`"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `check_lists`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.DELETE("/check_list/:checkListID", ch.DeleteCheckList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestDeleteCheckListHandlerShouldReturnsStatusBadRequest(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListHandler(repository.NewCheckListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	checkListID := uint(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/check_list/%d", checkListID), nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `check_lists`.* FROM `check_lists`"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `check_lists`")).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	r.DELETE("/check_list/:checkListID", ch.DeleteCheckList)
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

func TestIndexCheckListHandlerShouldReturnsStatusOKWithCheckListData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListHandler(repository.NewCheckListRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	boardID := uint(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/board/%d/check_list", boardID), nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mockCheckList := entity.CheckList{
		ID:     uint(3),
		Title:  "mockCheckList",
		CardID: uint(4),
	}

	mockCheckListItem := entity.CheckListItem{
		ID:          uint(5),
		Name:        "mockCheckListItem",
		Check:       false,
		CheckListID: mockCheckList.ID,
	}

	checkListQuery := utils.ReplaceQuotationForQuery(`
		SELECT check_lists.id, check_lists.title, check_lists.card_id
		FROM 'check_lists'`)

	mock.ExpectQuery(regexp.QuoteMeta(checkListQuery)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "card_id"}).
				AddRow(mockCheckList.ID, mockCheckList.Title, mockCheckList.CardID))

	checkListItemQuery := utils.ReplaceQuotationForQuery(`
		SELECT check_list_items.id, check_list_items.name, check_list_items.check_list_id, check_list_items.check
		FROM 'check_list_items'
		WHERE ('check_list_id' IN (?))`)

	mock.ExpectQuery(regexp.QuoteMeta(checkListItemQuery)).
		WithArgs(mockCheckListItem.CheckListID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "check", "check_list_id"}).
				AddRow(mockCheckListItem.ID, mockCheckListItem.Name, mockCheckListItem.Check, mockCheckListItem.CheckListID))

	r.GET("/board/:boardID/check_list", ch.IndexCheckList)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]entity.CheckList{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["check_lists"][0].ID, mockCheckList.ID)
	assert.Equal(t, res["check_lists"][0].ID, mockCheckList.ID)
	assert.Equal(t, res["check_lists"][0].Title, mockCheckList.Title)
	assert.Equal(t, res["check_lists"][0].CardID, mockCheckList.CardID)
	assert.Equal(t, res["check_lists"][0].Items[0].ID, mockCheckListItem.ID)
	assert.Equal(t, res["check_lists"][0].Items[0].Name, mockCheckListItem.Name)
	assert.Equal(t, res["check_lists"][0].Items[0].Check, mockCheckListItem.Check)
	assert.Equal(t, res["check_lists"][0].Items[0].CheckListID, mockCheckListItem.CheckListID)
}

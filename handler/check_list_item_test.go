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

type checkListItemRequestBody struct {
	Name  string `json:"name"`
	Check bool   `json:"check"`
}

func TestCreateCheckListItemHandlerShouldReturnsStatusCreatedWithCheckListItemData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListItemHandler(repository.NewCheckListItemRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	name := "new check list item"
	checkListID := uint(1)

	b, err := json.Marshal(checkListItemRequestBody{
		Name: name,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/check_list/%d/item", checkListID), bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `check_list_items`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.POST("/check_list/:checkListID/item", ch.CreateCheckListItem)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.CheckListItem{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["check_list_item"].Name, name)
	assert.Equal(t, res["check_list_item"].CheckListID, checkListID)
	assert.Equal(t, res["check_list_item"].ID, uint(1))
}

func TestShouldFailureCreateCheckListItemHandlerWhenRecievedInvalidParameter(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListItemHandler(repository.NewCheckListItemRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	checkListID := uint(1)

	b, err := json.Marshal(checkListItemRequestBody{
		Name: "",
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/check_list/%d/item", checkListID), bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards`"))
	mock.ExpectBegin()

	r.POST("/check_list/:checkListID/item", ch.CreateCheckListItem)
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

func TestUpdateCheckListItemHandlerShouldReturnsStatusOK(t *testing.T) {
	type testCase struct {
		testName    string
		requestBody checkListItemRequestBody
		url         string
	}

	const checkListItemID = uint(1)

	testCases := []testCase{
		{
			testName:    "when update a name",
			requestBody: checkListItemRequestBody{Name: "update name"},
			url:         fmt.Sprintf("/check_list_item/%d/name", checkListItemID),
		}, {
			testName:    "when update a check",
			requestBody: checkListItemRequestBody{Check: true},
			url:         fmt.Sprintf("/check_list_item/%d/check", checkListItemID),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock := utils.NewDBMock(t)
			defer db.Close()

			ch := NewCheckListItemHandler(repository.NewCheckListItemRepository(db))
			uh := NewUserHandler(repository.NewUserRepository(db))

			r := utils.SetUpRouter()

			b, err := json.Marshal(tc.requestBody)

			if err != nil {
				t.Fatalf("fail to marshal json: %v", err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, tc.url, bytes.NewReader(b))

			utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

			mock.ExpectQuery(regexp.QuoteMeta("SELECT `check_list_items`.* FROM `check_list_items`")).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(checkListItemID, "sample name"))

			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta("UPDATE `check_list_items` SET")).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()

			r.PATCH("/check_list_item/:checkListItemID/:attribute", ch.UpdateCheckListItem)
			r.ServeHTTP(w, req)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %v", err)
			}

			assert.Equal(t, w.Code, 200)
		})
	}
}

func TestShouldFailureUpdateCheckListItemHandlerWhenRecievedInvalidParameter(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListItemHandler(repository.NewCheckListItemRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	checkListItemID := uint(1)

	b, err := json.Marshal(checkListItemRequestBody{
		Name: "efwe",
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/check_list_item/%d/test", checkListItemID), bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `check_list_items`.* FROM `check_list_items`"))

	r.PATCH("/check_list_item/:checkListItemID/:attribute", ch.UpdateCheckListItem)
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

func TestDeleteCheckListItemHandlerShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListItemHandler(repository.NewCheckListItemRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	checkListItemID := uint(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/check_list_item/%d", checkListItemID), nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `check_list_items`.* FROM `check_list_items`"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `check_list_items`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.DELETE("/check_list_item/:checkListItemID", ch.DeleteCheckListItem)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestDeleteCheckListItemHandlerShouldReturnsStatusBadRequest(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCheckListItemHandler(repository.NewCheckListItemRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	checkListItemID := uint(1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/check_list_item/%d", checkListItemID), nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `check_list_items`.* FROM `check_list_items`"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `check_list_items`")).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	r.DELETE("/check_list_item/:checkListItemID", ch.DeleteCheckListItem)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 400)
}

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

type cardLabelRequestBody struct {
	LabelID uint `json:"label_id"`
}

func TestCreateCardLabelHandlerShouldReturnsStatusCreatedWithLabelData(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCardLabelHandler(repository.NewCardLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	labelID := uint(1)

	b, err := json.Marshal(cardLabelRequestBody{
		LabelID: labelID,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/card/1/card_label", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id FROM `boards` Join lists ON boards.id = lists.board_id Join labels ON boards.id = labels.board_id Join cards ON lists.id = cards.list_id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `card_labels` (`card_id`,`label_id`)")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `labels`")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(labelID))

	r.POST("/card/:cardID/card_label", ch.CreateCardLabel)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]entity.Label{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["label"].ID, labelID)
}

func TestShouldFailureCreateCardLabelHandlerWhenWithoutLabeID(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCardLabelHandler(repository.NewCardLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	b, err := json.Marshal(cardLabelRequestBody{})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/card/1/card_label", bytes.NewReader(b))

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	r.POST("/card/:cardID/card_label", ch.CreateCardLabel)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, "invalid parameters")
}

func TestDeleteCardLabelHandlerShouldReturnsStatusOK(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCardLabelHandler(repository.NewCardLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/card/1/card_label/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `card_labels`.* FROM `card_labels` Join labels ON card_labels.label_id = labels.id Join boards ON labels.board_id = boards.id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `card_labels`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.DELETE("/card/:cardID/card_label/:labelID", ch.DeleteCardLabel)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestDeleteCardLabelHandlerShouldReturnsStatusBadRequest(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	ch := NewCardLabelHandler(repository.NewCardLabelRepository(db))
	uh := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/card/1/card_label/1", nil)

	utils.SetUpAuthentication(r, req, mock, uh.Authenticate(), MapIDParamsToContext())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `card_labels`.* FROM `card_labels` Join labels ON card_labels.label_id = labels.id Join boards ON labels.board_id = boards.id"))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `card_labels`")).
		WillReturnResult(sqlmock.NewResult(1, 0))

	mock.ExpectCommit()

	r.DELETE("/card/:cardID/card_label/:labelID", ch.DeleteCardLabel)
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

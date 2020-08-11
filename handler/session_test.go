package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"local.packages/repository"
	"local.packages/utils"
	"local.packages/validator"
)

type sessionRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenRequestBody struct {
	RefreshToken string `json:"refresh_token"`
}

func TestShouldReturnsStatusOKWithSessionTokenUponUserSignIn(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	h := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()
	r.POST("/session", h.CreateSession)

	email := "gopher@sample.com"
	password := "12345678"

	b, err := json.Marshal(sessionRequestBody{
		Email:    email,
		Password: password,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	passwordDigest, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password_digest"}).AddRow(uint(1), email, passwordDigest))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/session", bytes.NewReader(b))

	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]interface{}{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["email"], email)
	assert.NotNil(t, res["access_token"])
	assert.NotNil(t, res["refresh_token"])
	assert.NotNil(t, res["expires_in"])
}

func TestShouldFailureCreateSessionHandler(t *testing.T) {
	type testCase struct {
		testName           string
		expectedStatus     int
		expectedError      string
		sessionRequestBody sessionRequestBody
	}

	testCases := []testCase{
		{
			testName:       "when without an email",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("メールアドレス"),
			sessionRequestBody: sessionRequestBody{
				Email:    "",
				Password: "12345678",
			},
		}, {
			testName:       "when without a password",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("パスワード"),
			sessionRequestBody: sessionRequestBody{
				Email:    "gopher@sample.com",
				Password: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, _ := utils.NewDBMock(t)
			defer db.Close()

			h := NewUserHandler(repository.NewUserRepository(db))

			r := utils.SetUpRouter()
			r.POST("/session", h.CreateSession)

			b, err := json.Marshal(tc.sessionRequestBody)

			if err != nil {
				t.Fatalf("fail to marshal json: %v", err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/session", bytes.NewReader(b))

			r.ServeHTTP(w, req)

			res := map[string][]validator.ValidationError{}

			if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
				t.Fatalf("fail to unmarshal response body. %v", err)
			}

			assert.Equal(t, w.Code, tc.expectedStatus)
			assert.Equal(t, res["errors"][0].Text, tc.expectedError)
		})
	}
}

func TestShouldReturnsStatusOKWithNewSessionPayload(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	h := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()
	r.PATCH("/session", h.UpdateSession)

	email := "gopher@sample.com"
	accessToken := "ercmewaorijno"
	refreshToken := "oirjnoinoiaec"

	findQuery := utils.ReplaceQuotationForQuery(`
		SELECT * FROM 'users'
		WHERE (remember_token = ?) AND (refresh_token = ?)
		ORDER BY 'users'.'id' ASC
		LIMIT 1`)

	updateQuery := utils.ReplaceQuotationForQuery(`
		UPDATE 'users'
		SET 'expires_at' = ?, 'refresh_token' = ?, 'remember_token' = ?, 'updated_at' = ?
		WHERE 'users'.'id' = ?`)

	mock.ExpectQuery(regexp.QuoteMeta(findQuery)).
		WithArgs(accessToken, refreshToken).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(uint(1), email))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	b, err := json.Marshal(tokenRequestBody{
		RefreshToken: refreshToken,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/session", bytes.NewReader(b))

	req.Header.Add("X-Auth-Token", accessToken)

	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]interface{}{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["ok"], true)
	assert.Equal(t, res["email"], email)
	assert.NotNil(t, res["access_token"])
	assert.NotNil(t, res["refresh_token"])
	assert.NotNil(t, res["expires_in"])
}

func TestShouldReturnsStatusOKWithoutNewSessionPayloadWhenTokenIsInvalid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	h := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()
	r.PATCH("/session", h.UpdateSession)

	accessToken := "ercmewaorijno"
	refreshToken := "oirjnoinoiaec"

	findQuery := utils.ReplaceQuotationForQuery(`
		SELECT * FROM 'users'
		WHERE (remember_token = ?) AND (refresh_token = ?)
		ORDER BY 'users'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(findQuery)).
		WithArgs(accessToken, refreshToken).
		WillReturnError(gorm.ErrRecordNotFound)

	b, err := json.Marshal(tokenRequestBody{
		RefreshToken: refreshToken,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/session", bytes.NewReader(b))

	req.Header.Add("X-Auth-Token", accessToken)

	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]interface{}{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, res["ok"], false)
	assert.Nil(t, res["access_token"])
	assert.Nil(t, res["refresh_token"])
	assert.Nil(t, res["expires_in"])
}

func TestShouldReturnsStatusBadRequestWhenFailedUpdateUser(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	h := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()
	r.PATCH("/session", h.UpdateSession)

	accessToken := "ercmewaorijno"
	refreshToken := "oirjnoinoiaec"

	findQuery := utils.ReplaceQuotationForQuery(`
		SELECT * FROM 'users'
		WHERE (remember_token = ?) AND (refresh_token = ?)
		ORDER BY 'users'.'id' ASC
		LIMIT 1`)

	updateQuery := utils.ReplaceQuotationForQuery(`
		UPDATE 'users'
		SET 'expires_at' = ?, 'refresh_token' = ?, 'remember_token' = ?, 'updated_at' = ?
		WHERE 'users'.'id' = ?`)

	mock.ExpectQuery(regexp.QuoteMeta(findQuery)).
		WithArgs(accessToken, refreshToken).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
		WillReturnError(errors.New("some error"))

	mock.ExpectRollback()

	b, err := json.Marshal(tokenRequestBody{
		RefreshToken: refreshToken,
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPatch, "/session", bytes.NewReader(b))

	req.Header.Add("X-Auth-Token", accessToken)

	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, repository.ErrorAuthenticationFailed)
}

func TestShouldReturnStatusOKWhenSucceedDeleteSession(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	h := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/session", nil)

	utils.SetUpAuthentication(r, req, mock, h.Authenticate(), MapIDParamsToContext())

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'users'
		SET 'refresh_token' = ?, 'remember_token' = ?
		WHERE (id = ?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(nil, nil, uint(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r.DELETE("/session", h.DeleteSession)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, w.Code, 200)
}

func TestShouldReturnStatusBadRequestWhenFailedDeleteSession(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	h := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/session", nil)

	utils.SetUpAuthentication(r, req, mock, h.Authenticate(), MapIDParamsToContext())

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'users'
		SET 'refresh_token' = ?, 'remember_token' = ?
		WHERE (id = ?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(nil, nil, uint(1)).
		WillReturnError(errors.New("some error"))

	mock.ExpectRollback()

	r.DELETE("/session", h.DeleteSession)
	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string][]validator.ValidationError{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 400)
	assert.Equal(t, res["errors"][0].Text, repository.ErrorAuthenticationFailed)
}

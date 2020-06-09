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
	"golang.org/x/crypto/bcrypt"

	"local.packages/repository"
	"local.packages/utils"
	"local.packages/validator"
)

type sessionRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
			sqlmock.NewRows([]string{"id", "password_digest"}).AddRow(uint(1), passwordDigest))

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

	res := map[string]string{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 200)
	assert.NotNil(t, res["token"])
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
			expectedError:  "Email must exist",
			sessionRequestBody: sessionRequestBody{
				Email:    "",
				Password: "12345678",
			},
		}, {
			testName:       "when without a password",
			expectedStatus: 400,
			expectedError:  "Password must exist",
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

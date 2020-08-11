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

type userRequestBody struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func TestShouldReturnsStatusCreatedWithSessionTokenUponUserSignUp(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	h := NewUserHandler(repository.NewUserRepository(db))

	r := utils.SetUpRouter()
	r.POST("/user", h.CreateUser)

	email := "gopher@sample.com"

	b, err := json.Marshal(userRequestBody{
		Name:                 "gopher",
		Email:                email,
		Password:             "12345678",
		PasswordConfirmation: "12345678",
	})

	if err != nil {
		t.Fatalf("fail to marshal json: %v", err)
	}

	passwordDigest, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs("gopher@sample.com").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password_digest"}).AddRow(uint(1), "gopher@sample.com", passwordDigest))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(b))

	r.ServeHTTP(w, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}

	res := map[string]interface{}{}

	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("fail to unmarshal response body. %v", err)
	}

	assert.Equal(t, w.Code, 201)
	assert.Equal(t, res["email"], email)
	assert.NotNil(t, res["access_token"])
	assert.NotNil(t, res["refresh_token"])
	assert.NotNil(t, res["expires_in"])
}

func TestShouldFailureCreateUserHandler(t *testing.T) {
	type testCase struct {
		testName        string
		expectedStatus  int
		expectedError   string
		userRequestBody userRequestBody
	}

	testCases := []testCase{
		{
			testName:       "when without a name",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("ユーザー名"),
			userRequestBody: userRequestBody{
				Name:                 "",
				Email:                "gopher@sample.com",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
		}, {
			testName:       "when without an email",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("メールアドレス"),
			userRequestBody: userRequestBody{
				Name:                 "gopher",
				Email:                "",
				Password:             "12345678",
				PasswordConfirmation: "12345678",
			},
		}, {
			testName:       "when without a password",
			expectedStatus: 400,
			expectedError:  validator.ErrorRequired("パスワード"),
			userRequestBody: userRequestBody{
				Name:                 "gopher",
				Email:                "gopher@sample.com",
				Password:             "",
				PasswordConfirmation: "12345678",
			},
		}, {
			testName:       "when does not match password and password confirmation",
			expectedStatus: 400,
			expectedError:  validator.ErrorEqualField("パスワード", "パスワード（確認用）"),
			userRequestBody: userRequestBody{
				Name:                 "gopher",
				Email:                "gopher@sample.com",
				Password:             "12345678",
				PasswordConfirmation: "123456789",
			},
		}, {
			testName:       "when password less than 8 characters",
			expectedStatus: 400,
			expectedError:  validator.ErrorTooShort("パスワード", "8"),
			userRequestBody: userRequestBody{
				Name:                 "gopher",
				Email:                "gopher@sample.com",
				Password:             "1234567",
				PasswordConfirmation: "1234567",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, _ := utils.NewDBMock(t)
			defer db.Close()

			h := NewUserHandler(repository.NewUserRepository(db))

			r := utils.SetUpRouter()
			r.POST("/user", h.CreateUser)

			b, err := json.Marshal(tc.userRequestBody)

			if err != nil {
				t.Fatalf("fail to marshal json: %v", err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(b))

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

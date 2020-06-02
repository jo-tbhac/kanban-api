package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"local.packages/utils"
)

func TestShouldSuccessfullyCreateUser(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	name := "gopher"
	email := "gopher@sample.com"
	password := "password"
	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`name`,`email`,`password_digest`,`remember_token`)")).
		WithArgs(createdAt, updatedAt, name, email, password, "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	u, _ := r.Create(name, email, password)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, u.ID, uint(1))
	assert.Equal(t, u.Name, name)
	assert.Equal(t, u.Email, email)
	assert.Equal(t, u.PasswordDigest, password)
}

func TestShouldNotCreateUserWhenDuplicateEmail(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	name := "gopher"
	email := "gopher@sample.com"
	password := "password"
	createdAt := utils.AnyTime{}
	updatedAt := utils.AnyTime{}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `users` (`created_at`,`updated_at`,`name`,`email`,`password_digest`,`remember_token`)")).
		WithArgs(createdAt, updatedAt, name, email, password, "").
		WillReturnError(fmt.Errorf("Error 1062: Duplicate entry '%s' for key 'email'", email))

	mock.ExpectRollback()

	_, err := r.Create(name, email, password)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	if err == nil {
		t.Error("there were unfulfilled expectations: it did not recieve error")
	}

	assert.Equal(t, err[0].Text, fmt.Sprintf("%s has already been taken", email))
}

func TestShouldSuccessfullySignIn(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	id := uint(1)
	email := "gopher@sample.com"
	password := "password"
	passwordDigest, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password_digest"}).AddRow(id, email, passwordDigest))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `users`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	u, err := r.SignIn(email, password)

	if err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, u.ID, id)
	assert.Equal(t, u.Email, email)
	assert.Equal(t, u.PasswordDigest, string(passwordDigest))
}

func TestShouldNotSignInWhenEmailDoesNotExist(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	email := "gopher@sample.com"
	password := "password"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.SignIn(email, password)

	if err == nil {
		t.Errorf("was expected an error, but did not recieve it. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "user does not exist")
}

func TestShouldNotSignInWhenPasswordIsInvalid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	email := "gopher@sample.com"
	registeredPassword := "password"
	password := "invalid_password"
	passwordDigest, _ := bcrypt.GenerateFromPassword([]byte(registeredPassword), bcrypt.DefaultCost)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"password_digest"}).AddRow(passwordDigest))

	_, err := r.SignIn(email, password)

	if err == nil {
		t.Errorf("was expected an error, but did not recieve it. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, "invalid password")
}

func TestIsSignedInShouldReturnTrueWhenTokenIsValid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	id := uint(1)
	token := "sample_token"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

	u, ok := r.IsSignedIn(token)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, u.ID, id)
	assert.True(t, ok)
}

func TestIsSignedInShouldReturnFalseWhenTokenIsInvalid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	token := "sample_token"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
		WithArgs(token).
		WillReturnError(gorm.ErrRecordNotFound)

	_, ok := r.IsSignedIn(token)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.False(t, ok)
}

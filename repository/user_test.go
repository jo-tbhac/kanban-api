package repository

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

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

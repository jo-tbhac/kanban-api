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
	"local.packages/validator"
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

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'users' ('created_at','updated_at','name','email','password_digest','remember_token')
		VALUES (?,?,?,?,?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
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

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'users' ('created_at','updated_at','name','email','password_digest','remember_token')
		VALUES (?,?,?,?,?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(createdAt, updatedAt, name, email, password, "").
		WillReturnError(fmt.Errorf("Error 1062: Duplicate entry '%s' for key 'email'", email))

	mock.ExpectRollback()

	_, err := r.Create(name, email, password)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	if err == nil {
		t.Error("was not expected an error, but did not recieve it")
	}

	assert.Equal(t, err[0].Text, validator.ErrorAlreadyBeenTaken)
}

func TestShouldSuccessfullySignIn(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	id := uint(1)
	email := "gopher@sample.com"
	password := "password"
	passwordDigest, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	findQuery := "SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1"
	updateQuery := "UPDATE `users` SET `remember_token` = ?, `updated_at` = ? WHERE `users`.`id` = ?"

	mock.ExpectQuery(regexp.QuoteMeta(findQuery)).
		WithArgs(email).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "email", "password_digest"}).AddRow(id, email, passwordDigest))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(updateQuery)).
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

	query := "SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := r.SignIn(email, password)

	if err == nil {
		t.Error("was not expected an error, but did not recieve it")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorUserDoesNotExist)
}

func TestShouldNotSignInWhenPasswordIsInvalid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	email := "gopher@sample.com"
	registeredPassword := "password"
	password := "invalid_password"
	passwordDigest, _ := bcrypt.GenerateFromPassword([]byte(registeredPassword), bcrypt.DefaultCost)

	query := "SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"password_digest"}).AddRow(passwordDigest))

	_, err := r.SignIn(email, password)

	if err == nil {
		t.Error("was expected an error, but did not recieve it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorInvalidPassword)
}

func TestIsSignedInShouldReturnTrueWhenTokenIsValid(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	id := uint(1)
	token := "sample_token"

	query := "SELECT * FROM `users`  WHERE (remember_token = ?) ORDER BY `users`.`id` ASC LIMIT 1"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
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

	query := "SELECT * FROM `users`  WHERE (remember_token = ?) ORDER BY `users`.`id` ASC LIMIT 1"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(token).
		WillReturnError(gorm.ErrRecordNotFound)

	_, ok := r.IsSignedIn(token)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.False(t, ok)
}

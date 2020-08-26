package repository

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"local.packages/entity"
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
	expiresAt := utils.AnyTime{}

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'users' ('created_at','updated_at','name','email','password_digest','remember_token','refresh_token','expires_at')
		VALUES (?,?,?,?,?,?,?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(createdAt, updatedAt, name, email, password, "", "", expiresAt).
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
	expiresAt := utils.AnyTime{}

	query := utils.ReplaceQuotationForQuery(`
		INSERT INTO 'users' ('created_at','updated_at','name','email','password_digest','remember_token','refresh_token','expires_at')
		VALUES (?,?,?,?,?,?,?,?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(createdAt, updatedAt, name, email, password, "", "", expiresAt).
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
	updateQuery := "UPDATE `users` SET `expires_at` = ?, `refresh_token` = ?, `remember_token` = ?, `updated_at` = ? WHERE `users`.`id` = ?"

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
		t.Error("was expected an error, but did not recieve it")
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

func TestShouldSuccessfullyValidateToken(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	accessToken := "dsfjsefsefljfsf"
	refreshToken := "eskljfnaejfauh"
	userID := uint(1)

	query := utils.ReplaceQuotationForQuery(`
		SELECT * FROM 'users'
		WHERE (remember_token = ?) AND (refresh_token = ?)
		ORDER BY 'users'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(accessToken, refreshToken).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "remember_token", "refresh_token"}).
				AddRow(userID, accessToken, refreshToken))

	u, ok := r.ValidateToken(accessToken, refreshToken)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.True(t, ok)
	assert.Equal(t, u.ID, userID)
	assert.Equal(t, u.RememberToken, accessToken)
	assert.Equal(t, u.RefreshToken, refreshToken)
}

func TestShouldFailureValidateToken(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	accessToken := "dsfjsefsefljfsf"
	refreshToken := "eskljfnaejfauh"

	query := utils.ReplaceQuotationForQuery(`
		SELECT * FROM 'users'
		WHERE (remember_token = ?) AND (refresh_token = ?)
		ORDER BY 'users'.'id' ASC
		LIMIT 1`)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(accessToken, refreshToken).
		WillReturnError(gorm.ErrRecordNotFound)

	_, ok := r.ValidateToken(accessToken, refreshToken)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.False(t, ok)
}

func TestShouldSuccessfullyUpdateUserSession(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	u := &entity.User{ID: uint(1)}

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'users'
		SET 'expires_at' = ?, 'refresh_token' = ?, 'remember_token' = ?, 'updated_at' = ?
		WHERE 'users'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.UpdateUserSession(u); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.NotNil(t, u.RememberToken)
	assert.NotNil(t, u.RefreshToken)
}

func TestShouldFailureUpdateUserSession(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	u := &entity.User{ID: uint(1)}

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'users'
		SET 'expires_at' = ?, 'refresh_token' = ?, 'remember_token' = ?, 'updated_at' = ?
		WHERE 'users'.'id' = ?`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WillReturnError(errors.New("some error"))

	err := r.UpdateUserSession(u)

	if err == nil {
		t.Errorf("was expected an error but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorAuthenticationFailed)
}

func TestShouldSuccessfullySignOut(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	userID := uint(1)

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'users'
		SET 'expires_at' = ?, 'refresh_token' = ?, 'remember_token' = ?
		WHERE (id = ?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(nil, nil, nil, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	if err := r.SignOut(userID); err != nil {
		t.Errorf("was not expected an error. %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestShouldFailureSignOut(t *testing.T) {
	db, mock := utils.NewDBMock(t)
	defer db.Close()

	r := NewUserRepository(db)

	userID := uint(1)

	query := utils.ReplaceQuotationForQuery(`
		UPDATE 'users'
		SET 'expires_at' = ?, 'refresh_token' = ?, 'remember_token' = ?
		WHERE (id = ?)`)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(nil, nil, nil, userID).
		WillReturnError(errors.New("some error"))

	mock.ExpectRollback()

	err := r.SignOut(userID)

	if err == nil {
		t.Errorf("was expected an error but did not recieved it.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}

	assert.Equal(t, err[0].Text, ErrorAuthenticationFailed)
}

package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	Email          string    `json:"email" gorm:"unique;not null"`
	PasswordDigest string    `json:"password_digest" gorm:"not null"`
	RememberToken  string    `json:"remember_token" gorm:"not null"`
	Boards         []Board   `json:"boards" gorm:"foreignkey:UserID"`
}

func (u *User) Create(name, email, pw string) []validator.ValidationError {
	db := db.Get()
	passwordDigest, err := encryptPassword(pw)

	if err != nil {
		log.Printf("fail to encrypted password: %v", err)
		return validator.NewValidationErrors("internal server error")
	}

	u.Name = name
	u.Email = email
	u.PasswordDigest = passwordDigest

	if err := db.Create(u).Error; err != nil {
		return validator.FormattedMySQLError(err)
	}

	return nil
}

func (u *User) SignIn(email, password string) error {
	db := db.Get()

	if db.Where("email = ?", email).First(u); u.ID == 0 {
		return errors.New("user does not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password)); err != nil {
		return errors.New("invalid password")
	}

	t, err := newSessionToken()

	if err != nil {
		log.Printf("failed create token. %v", err)
		return errors.New("internal server error")
	}

	db.Model(u).Select("remember_token").Updates(map[string]interface{}{"remember_token": t})

	return nil
}

func (u *User) IsSignedIn(token string) bool {
	db := db.Get()
	db.Where("remember_token = ?", token).First(u)

	return u.ID != 0
}

func newSessionToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func encryptPassword(password string) (digest string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	digest = string(h)

	return
}

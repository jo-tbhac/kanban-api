package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/jo-tbhac/kanban-api/db"
	"golang.org/x/crypto/bcrypt"
)

const UserDoesNotExist = 0

type User struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordDigest string    `json:"password_digest"`
	RememberToken  string    `json:"remember_token"`
	Boards         []Board   `json:"boards" gorm:"foreignkey:UserID"`
}

type UserParams struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

func init() {
	db := db.Get()
	db.AutoMigrate(&User{})
	db.Model(&User{}).AddUniqueIndex("idx_users_email", "email")
}

func (u *User) Create(p UserParams) error {
	db := db.Get()

	passwordDigest, err := encryptPassword(p.Password)
	if err != nil {
		return err
	}

	u.PasswordDigest = passwordDigest
	u.Name = p.Name
	u.Email = p.Email

	if err := db.Create(u).Error; err != nil {
		return err
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

	return u.ID != UserDoesNotExist
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

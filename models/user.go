package models

import (
	"time"

	"github.com/jo-tbhac/kanban-api/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordDigest string    `json:"password_digest"`
	RememberToken  string    `json:"remember_token"`
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

	if err := db.Create(&u).Error; err != nil {
		return err
	}

	return nil
}

func encryptPassword(password string) (digest string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	digest = string(h)

	return
}

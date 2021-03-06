package entity

import "time"

// User is model of users table.
type User struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	Email          string    `json:"email" gorm:"unique;not null"`
	PasswordDigest string    `json:"password_digest" gorm:"not null"`
	RememberToken  string    `json:"remember_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresAt      time.Time `json:"-" gorm:"not null"`
	Boards         []Board   `json:"boards" gorm:"foreignkey:UserID"`
}

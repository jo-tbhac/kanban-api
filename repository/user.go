package repository

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/jo-tbhac/kanban-api/entity"
	"github.com/jo-tbhac/kanban-api/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(name, email, password string) (*entity.User, []validator.ValidationError) {
	u := &entity.User{
		Name:           name,
		Email:          email,
		PasswordDigest: password,
	}

	if err := r.db.Create(u).Error; err != nil {
		return u, validator.FormattedMySQLError(err)
	}

	return u, nil
}

func (r *UserRepository) SignIn(email, password string) (*entity.User, []validator.ValidationError) {
	u := &entity.User{}

	if r.db.Where("email = ?", email).First(u).RecordNotFound() {
		return u, validator.NewValidationErrors("user does not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password)); err != nil {
		return u, validator.NewValidationErrors("invalid password")
	}

	t, err := newSessionToken()

	if err != nil {
		log.Printf("fail to create token: %v", err)
		return u, validator.NewValidationErrors("internal server error")
	}

	r.db.Model(u).Select("remember_token").Updates(map[string]interface{}{"remember_token": t})

	return u, nil
}

func (r *UserRepository) IsSignedIn(token string) (*entity.User, bool) {
	u := &entity.User{}

	if r.db.Where("remember_token = ?", token).First(u).RecordNotFound() {
		return u, false
	}

	return u, true
}

func (r *UserRepository) EncryptPassword(password string) (digest string, err error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	digest = string(h)

	return
}

func newSessionToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

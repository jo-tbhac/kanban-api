package repository

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"local.packages/entity"
	"local.packages/validator"
)

// UserRepository ...
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository is constructor for UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create insert a new record to users table.
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

// SignIn returns instance of User that contain a new session token.
// returns errors with message if an email or password is invalid.
func (r *UserRepository) SignIn(email, password string) (*entity.User, []validator.ValidationError) {
	u := &entity.User{}

	if r.db.Where("email = ?", email).First(u).RecordNotFound() {
		return u, validator.NewValidationErrors(ErrorUserDoesNotExist)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password)); err != nil {
		return u, validator.NewValidationErrors(ErrorInvalidPassword)
	}

	t, err := newSessionToken()

	if err != nil {
		log.Printf("fail to create token: %v", err)
		return u, validator.NewValidationErrors(ErrorAuthenticationFailed)
	}

	expire := time.Now().Add(time.Hour * 2)

	r.db.Model(u).Updates(map[string]interface{}{"remember_token": t, "expires_at": expire})

	return u, nil
}

// IsSignedIn returns an instance of User that found by session token.
// returns `false` if the record not found.
func (r *UserRepository) IsSignedIn(token string) (*entity.User, bool) {
	u := &entity.User{}

	if r.db.Where("remember_token = ?", token).First(u).RecordNotFound() {
		return u, false
	}

	return u, true
}

// EncryptPassword returns the bcrypt hash of the password.
func (r *UserRepository) EncryptPassword(password string) (string, []validator.ValidationError) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", validator.NewValidationErrors(ErrorAuthenticationFailed)
	}

	digest := string(h)

	return digest, nil
}

func newSessionToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

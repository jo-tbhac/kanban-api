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

// TestUsers returns user instances for testing.
func (r *UserRepository) TestUsers() *[]entity.User {
	var us []entity.User

	r.db.Select("id, email, expires_at").Where("id IN (?)", []uint{2, 3, 4, 5}).Find(&us)

	return &us
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

// UpdateUserSession updates access token, refresh token and expires.
func (r *UserRepository) UpdateUserSession(u *entity.User) []validator.ValidationError {
	at, err := newSessionToken()

	if err != nil {
		log.Printf("fail to create access token: %v", err)
		return validator.NewValidationErrors(ErrorAuthenticationFailed)
	}

	rt, err := newSessionToken()

	if err != nil {
		log.Printf("fail to create refresh token: %v", err)
		return validator.NewValidationErrors(ErrorAuthenticationFailed)
	}

	expire := time.Now().Add(time.Hour * 1)

	if err := r.db.Model(u).Updates(map[string]interface{}{"remember_token": at, "refresh_token": rt, "expires_at": expire}).Error; err != nil {
		log.Printf("fail to update user session: %v", err)
		return validator.NewValidationErrors(ErrorAuthenticationFailed)
	}

	return nil
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

	if err := r.UpdateUserSession(u); err != nil {
		return u, err
	}

	return u, nil
}

// SignOut delete remember token and refresh token.
func (r *UserRepository) SignOut(uid uint) []validator.ValidationError {
	if err := r.db.Table("users").Where("id = ?", uid).Updates(map[string]interface{}{"remember_token": nil, "refresh_token": nil}).Error; err != nil {
		log.Printf("fail to delete session: %v", err)
		return validator.NewValidationErrors(ErrorAuthenticationFailed)
	}

	return nil
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

// ValidateToken returns an instance of User that found by access token and refresh token.
// returns `false` if the record not found.
func (r *UserRepository) ValidateToken(at, rt string) (*entity.User, bool) {
	u := &entity.User{}

	if r.db.Where("remember_token = ?", at).Where("refresh_token = ?", rt).First(u).RecordNotFound() {
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

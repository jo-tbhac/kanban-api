package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"local.packages/repository"
	"local.packages/validator"
)

type userParams struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

type UserHandler struct {
	repository *repository.UserRepository
}

func NewUserHandler(r *repository.UserRepository) *UserHandler {
	return &UserHandler{repository: r}
}

func (h UserHandler) CreateUser(c *gin.Context) {
	var p userParams

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.FormattedValidationError(err)})
		return
	}

	passwordDigest, err := h.repository.EncryptPassword(p.Password)

	if err != nil {
		log.Printf("fail to encrypted password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err})
		return
	}

	if _, err := h.repository.Create(p.Name, p.Email, passwordDigest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	u, err := h.repository.SignIn(p.Email, p.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": u.RememberToken})
}
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type UserParams struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var p UserParams

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.FormattedValidationError(err)})
		return
	}

	var u models.User

	if err := u.Create(p.Name, p.Email, p.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := u.SignIn(p.Email, p.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.NewValidationErrors(err.Error())})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": u.RememberToken})
}

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type sessionParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func createSession(c *gin.Context) {
	var p sessionParams

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.NewValidationErrors("invalid parameters")})
		return
	}

	var u models.User

	if err := u.SignIn(p.Email, p.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.NewValidationErrors(err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": u.RememberToken})
}

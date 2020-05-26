package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type SessionParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateSession(c *gin.Context) {
	var p SessionParams

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.MakeErrors("invalid parameters")})
		return
	}

	var u models.User

	if err := u.SignIn(p.Email, p.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.MakeErrors(err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": u.RememberToken})
}

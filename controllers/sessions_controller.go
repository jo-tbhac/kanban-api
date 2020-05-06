package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jo-tbhac/kanban-api/models"
)

func CreateSession(c *gin.Context) {
	var p models.SessionParams

	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var u models.User

	if err := u.SignIn(p.Email, p.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": u.RememberToken})
}

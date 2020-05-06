package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
)

func CreateUser(c *gin.Context) {
	var p models.UserParams

	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, []string{err.Error()})
		return
	}

	var u models.User

	if err := u.Create(p); err != nil {
		c.JSON(http.StatusBadRequest, []string{err.Error()})
		return
	}

	if err := u.SignIn(p.Email, p.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": u.RememberToken})
}

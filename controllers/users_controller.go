package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
)

func CreateUser(c *gin.Context) {
	var p models.UserParams

	if err := c.BindJSON(&p); err != nil {
		models.CustomValidateMessages(err)
		c.JSON(http.StatusBadRequest, []string{err.Error()})
		return
	}

	var u models.User

	if err := u.Create(p); err != nil {
		c.JSON(http.StatusBadRequest, []string{err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

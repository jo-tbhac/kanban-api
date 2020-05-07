package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
)

func CreateBoard(c *gin.Context) {
	var b models.Board

	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b.UserID = c.Keys["user"].(models.User).ID

	if err := b.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"board": b})
}

func IndexBoard(c *gin.Context) {
	var b []models.Board
	u := c.Keys["user"].(models.User)

	models.IndexBoard(&b, &u)
	c.JSON(http.StatusOK, gin.H{"boards": b})
}

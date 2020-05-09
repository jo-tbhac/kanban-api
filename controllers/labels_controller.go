package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
)

func CreateLabel(c *gin.Context) {
	var l models.Label

	if err := c.BindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := l.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"label": l})
}

func IndexLabel(c *gin.Context) {
	var l []models.Label

	bid, err := strconv.Atoi(c.Query("board_id"))

	if err != nil {
		log.Println("invalid query parameter `board_id`")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter"})
		return
	}

	uid := CurrentUser(c).ID

	if err := models.IndexLabel(&l, uint(bid), uid); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"labels": l})
}

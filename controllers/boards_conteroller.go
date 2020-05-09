package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
)

func CreateBoard(c *gin.Context) {
	var b models.Board

	if err := c.BindJSON(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	b.UserID = CurrentUser(c).ID

	if err := b.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"board": b})
}

func IndexBoard(c *gin.Context) {
	var b []models.Board
	u := CurrentUser(c)

	models.GetAllBoard(&b, &u)
	c.JSON(http.StatusOK, gin.H{"boards": b})
}

func ShowBoard(c *gin.Context) {
	bid, err := strconv.Atoi(c.Query("board_id"))

	if err != nil {
		log.Printf("failed cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	b := models.Board{ID: uint(bid)}

	uid := CurrentUser(c).ID

	if err := b.Get(uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": b})
}

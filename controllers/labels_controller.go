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

	if !models.RelatedBoardOwnerIsValid(l.BoardID, CurrentUser(c).ID) {
		log.Println("does not match uid and board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	if err := l.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"label": l})
}

func UpdateLabel(c *gin.Context) {
	var l models.Label

	if err := c.BindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	l.GetBoardID()

	if !models.RelatedBoardOwnerIsValid(l.BoardID, CurrentUser(c).ID) {
		log.Println("does not match uid and board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	if err := l.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": l})
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

	if err := models.GetAllLabel(&l, uint(bid), uid); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"labels": l})
}

func DeleteLabel(c *gin.Context) {
	lid, err := strconv.Atoi(c.Query("label_id"))

	if err != nil {
		log.Printf("failed cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	l := models.Label{ID: uint(lid)}

	l.GetBoardID()

	if !models.RelatedBoardOwnerIsValid(l.BoardID, CurrentUser(c).ID) {
		log.Println("does not match uid and board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	if err := l.Delete(); err != nil {
		log.Printf("failed delete a label: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed delete a label"})
		return
	}

	c.Status(http.StatusOK)
}

package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
)

func CreateList(c *gin.Context) {
	var l models.List

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

	c.JSON(http.StatusCreated, gin.H{"list": l})
}

func UpdateList(c *gin.Context) {
	var l models.List

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

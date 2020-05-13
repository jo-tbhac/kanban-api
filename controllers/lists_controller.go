package controllers

import (
	"log"
	"net/http"
	"strconv"

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

func DeleteList(c *gin.Context) {
	lid, err := strconv.Atoi(c.Query("list_id"))

	if err != nil {
		log.Printf("failed cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	l := models.List{ID: uint(lid)}

	l.GetBoardID()

	if !models.RelatedBoardOwnerIsValid(l.BoardID, CurrentUser(c).ID) {
		log.Println("does not match uid and board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	if err := l.Delete(); err != nil {
		log.Printf("failed delete a list: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed delete a list"})
		return
	}

	c.Status(http.StatusOK)
}

package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
)

func CreateCard(c *gin.Context) {
	var ca models.Card

	if err := c.BindJSON(&ca); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bid := ca.GetBoardID()

	if !models.RelatedBoardOwnerIsValid(bid, CurrentUser(c).ID) {
		log.Println("does not match uid and board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters"})
		return
	}

	if err := ca.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"card": ca})
}

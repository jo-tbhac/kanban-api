package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
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

func UpdateCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("failed cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	var ca models.Card

	if ca.Find(uint(id), CurrentUser(c).ID); ca.ID == 0 {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.MakeErrors("id is invalid")})
		return
	}

	if err := c.BindJSON(&ca); err != nil {
		log.Printf("failed bind JSON: %v", err)
		c.AbortWithStatus(500)
		return
	}

	if err := ca.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"card": ca})
}

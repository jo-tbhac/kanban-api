package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type CardParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createCard(c *gin.Context) {
	var p CardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	ca := models.Card{
		Title:  p.Title,
		ListID: getIDParam(c, "listID"),
	}

	if !ca.ValidateUID(currentUser(c).ID) {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid request")})
		return
	}

	if err := ca.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"card": ca})
}

func updateCard(c *gin.Context) {
	id := getIDParam(c, "cardID")
	var ca models.Card

	if ca.Find(id, currentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.NewValidationErrors("id is invalid")})
		return
	}

	var p CardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	ca.Title = p.Title

	if err := ca.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"card": ca})
}

func deleteCard(c *gin.Context) {
	id := getIDParam(c, "cardID")
	var ca models.Card

	if ca.Find(id, currentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.NewValidationErrors("id is invalid")})
		return
	}

	if err := ca.Delete(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

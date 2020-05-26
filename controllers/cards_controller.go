package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type CardParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateCard(c *gin.Context) {
	lid, err := strconv.Atoi(c.Query("list_id"))

	if err != nil {
		log.Printf("failed cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("list_id must be an integer")})
		return
	}

	var p CardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	ca := models.Card{
		Title:  p.Title,
		ListID: uint(lid),
	}

	if !ca.ValidateUID(CurrentUser(c).ID) {
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

func UpdateCard(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("failed cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("id must be an integer")})
		return
	}

	var ca models.Card

	if ca.Find(uint(id), CurrentUser(c).ID).RecordNotFound() {
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

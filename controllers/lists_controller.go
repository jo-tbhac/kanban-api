package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type ListParams struct {
	Name string `json:"name"`
}

func createList(c *gin.Context) {
	var p ListParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	bid := getIDParam(c, "boardID")

	if !models.ValidateUID(uint(bid), currentUser(c).ID) {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("board_id is invalid")})
		return
	}

	l := models.List{
		Name:    p.Name,
		BoardID: bid,
	}

	if err := l.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"list": l})
}

func updateList(c *gin.Context) {
	id := getIDParam(c, "labelID")
	var l models.List

	if l.Find(id, currentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id associated with the list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("id is invalid")})
		return
	}

	var p ListParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	l.Name = p.Name

	if err := l.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": l})
}

func deleteList(c *gin.Context) {
	id := getIDParam(c, "labelID")
	var l models.List

	if l.Find(id, currentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id associated with the list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("id is invalid")})
		return
	}

	if err := l.Delete(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}

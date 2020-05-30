package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type labelParams struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func createLabel(c *gin.Context) {
	var p labelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	bid := getIDParam(c, "boardID")

	if !models.ValidateUID(bid, currentUser(c).ID) {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("board_id is invalid")})
		return
	}

	l := models.Label{
		Name:    p.Name,
		Color:   p.Color,
		BoardID: bid,
	}

	if err := l.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"label": l})
}

func updateLabel(c *gin.Context) {
	id := getIDParam(c, "labelID")
	var l models.Label

	if l.Find(id, currentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("id is invalid")})
		return
	}

	var p labelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	l.Name = p.Name
	l.Color = p.Color

	if err := l.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": l})
}

func indexLabel(c *gin.Context) {
	bid := getIDParam(c, "boardID")
	var ls models.Labels

	ls.GetAll(bid, currentUser(c).ID)

	c.JSON(http.StatusOK, gin.H{"labels": ls})
}

func deleteLabel(c *gin.Context) {
	id := getIDParam(c, "labelID")
	var l models.Label

	if l.Find(id, currentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("id is invalid")})
		return
	}

	if err := l.Delete(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

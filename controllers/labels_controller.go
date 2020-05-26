package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type LabelParams struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func CreateLabel(c *gin.Context) {
	bid, err := strconv.Atoi(c.Query("board_id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("board_id must be an integer")})
		return
	}

	var p LabelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.AbortWithStatus(500)
		return
	}

	if !models.ValidateUID(uint(bid), CurrentUser(c).ID) {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("board_id is invalid")})
		return
	}

	l := models.Label{
		Name:    p.Name,
		Color:   p.Color,
		BoardID: uint(bid),
	}

	if err := l.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"label": l})
}

func UpdateLabel(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	var l models.Label

	if l.Find(uint(id), CurrentUser(c).ID); l.ID == 0 {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id is invalid")})
		return
	}

	var p LabelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.AbortWithStatus(500)
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

func IndexLabel(c *gin.Context) {
	bid, err := strconv.Atoi(c.Query("board_id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	var l []models.Label

	models.GetAllLabel(&l, uint(bid), CurrentUser(c).ID)

	c.JSON(http.StatusOK, gin.H{"labels": l})
}

func DeleteLabel(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	var l models.Label

	if l.Find(uint(id), CurrentUser(c).ID); l.ID == 0 {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id is invalid")})
		return
	}

	if err := l.Delete(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

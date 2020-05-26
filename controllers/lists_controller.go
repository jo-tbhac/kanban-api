package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type ListParams struct {
	Name string `json:"name"`
}

func CreateList(c *gin.Context) {
	bid, err := strconv.Atoi(c.Query("board_id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("board_id must be an integer")})
		return
	}

	var p ListParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("invalid parameters")})
		return
	}

	if !models.ValidateUID(uint(bid), CurrentUser(c).ID) {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("board_id is invalid")})
		return
	}

	l := models.List{
		Name:    p.Name,
		BoardID: uint(bid),
	}

	if err := l.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"list": l})
}

func UpdateList(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	var l models.List

	if l.Find(uint(id), CurrentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id associated with the list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id is invalid")})
		return
	}

	var p ListParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("invalid parameters")})
		return
	}

	l.Name = p.Name

	if err := l.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": l})
}

func DeleteList(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	var l models.List

	if l.Find(uint(id), CurrentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id associated with the list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id is invalid")})
		return
	}

	if err := l.Delete(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}

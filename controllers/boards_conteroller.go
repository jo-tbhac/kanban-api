package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type BoardParams struct {
	Name string `json:"name"`
}

func CreateBoard(c *gin.Context) {
	var p BoardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.AbortWithStatus(500)
		return
	}

	b := models.Board{
		Name:   p.Name,
		UserID: CurrentUser(c).ID,
	}

	if err := b.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"board": b})
}

func UpdateBoard(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	var b models.Board

	if b.Find(uint(id), CurrentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id is invalid")})
		return
	}

	var p BoardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.AbortWithStatus(500)
		return
	}

	b.Name = p.Name

	if err := b.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": b})
}

func IndexBoard(c *gin.Context) {
	var b []models.Board
	u := CurrentUser(c)

	models.GetAllBoard(&b, &u)
	c.JSON(http.StatusOK, gin.H{"boards": b})
}

func ShowBoard(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	var b models.Board

	uid := CurrentUser(c).ID

	if b.Find(uint(id), uid).RecordNotFound() {
		log.Println("uid does not match board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id is invalid")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": b})
}

func DeleteBoard(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		log.Printf("fail to cast string to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.MakeErrors("id must be an integer")})
		return
	}

	b := models.Board{
		ID:     uint(id),
		UserID: CurrentUser(c).ID,
	}

	if err := b.Delete(); err != nil {
		log.Printf("fail to delete a board: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

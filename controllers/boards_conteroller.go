package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type boardParams struct {
	Name string `json:"name"`
}

func createBoard(c *gin.Context) {
	var p boardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	b := models.Board{
		Name:   p.Name,
		UserID: currentUser(c).ID,
	}

	if err := b.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"board": b})
}

func updateBoard(c *gin.Context) {
	id := getIDParam(c, "boardID")
	var b models.Board

	if b.Find(id, currentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("id is invalid")})
		return
	}

	var p boardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	b.Name = p.Name

	if err := b.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": b})
}

func indexBoard(c *gin.Context) {
	var bs models.Boards
	u := currentUser(c)

	bs.GetAll(&u)
	c.JSON(http.StatusOK, gin.H{"boards": bs})
}

func showBoard(c *gin.Context) {
	id := getIDParam(c, "boardID")
	var b models.Board

	if b.Find(id, currentUser(c).ID).RecordNotFound() {
		log.Println("uid does not match board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("id is invalid")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": b})
}

func deleteBoard(c *gin.Context) {
	b := models.Board{
		ID:     getIDParam(c, "boardID"),
		UserID: currentUser(c).ID,
	}

	if err := b.Delete(); err != nil {
		log.Printf("fail to delete a board: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

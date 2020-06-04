package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"local.packages/repository"
	"local.packages/validator"
)

type boardParams struct {
	Name string `json:"name"`
}

type BoardHandler struct {
	repository *repository.BoardRepository
}

func NewBoardHandler(r *repository.BoardRepository) *BoardHandler {
	return &BoardHandler{repository: r}
}

func (h BoardHandler) CreateBoard(c *gin.Context) {
	var p boardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	b, err := h.repository.Create(p.Name, currentUserID(c))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"board": b})
}

func (h BoardHandler) UpdateBoard(c *gin.Context) {
	id := getIDParam(c, "boardID")
	b, err := h.repository.FindWithoutPreload(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	var p boardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	if err := h.repository.Update(b, p.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": b})
}

func (h BoardHandler) IndexBoard(c *gin.Context) {
	bs := h.repository.GetAll(currentUserID(c))
	c.JSON(http.StatusOK, gin.H{"boards": bs})
}

func (h BoardHandler) ShowBoard(c *gin.Context) {
	id := getIDParam(c, "boardID")
	b, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"board": b})
}

func (h BoardHandler) DeleteBoard(c *gin.Context) {
	id := getIDParam(c, "boardID")

	if err := h.repository.Delete(id, currentUserID(c)); err != nil {
		log.Printf("fail to delete a board: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

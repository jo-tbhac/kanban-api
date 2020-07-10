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

// BoardHandler ...
type BoardHandler struct {
	repository *repository.BoardRepository
}

// NewBoardHandler is constructor for BoardHandler.
func NewBoardHandler(r *repository.BoardRepository) *BoardHandler {
	return &BoardHandler{repository: r}
}

// CreateBoard call a function that create a new record to boards table.
// if creation was successful, returns status 201 and instance of Board as http response.
// if creation was failure, returns status 400 and error with messages.
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

// UpdateBoard call a function that update a record in boards table.
// if update was successful, returns status 200 and updated instance of Board as http response.
// if update was failure, returns status 400 and error with messages.
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

// IndexBoard returns status 200 and slice of Board instance as http response.
func (h BoardHandler) IndexBoard(c *gin.Context) {
	bs := h.repository.GetAll(currentUserID(c))
	c.JSON(http.StatusOK, gin.H{"boards": bs})
}

// ShowBoard returns status 200 and an instance of Board as response.
// if recieved invalid request, returns status 400 and errors with message.
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

// DeleteBoard call a function that delete a record from boards table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
func (h BoardHandler) DeleteBoard(c *gin.Context) {
	id := getIDParam(c, "boardID")

	if err := h.repository.Delete(id, currentUserID(c)); err != nil {
		log.Printf("fail to delete a board: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

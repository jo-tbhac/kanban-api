package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"local.packages/repository"
)

// CoverHandler ...
type CoverHandler struct {
	repository *repository.CoverRepository
}

// NewCoverHandler is constructor for CoverHandler.
func NewCoverHandler(r *repository.CoverRepository) *CoverHandler {
	return &CoverHandler{repository: r}
}

// CreateCover call a function that create a new record to covers table.
// if creation was successful, returns status 201 and instance of Cover as http response.
// if creation was failure, returns status 400 and error with messages.
func (h CoverHandler) CreateCover(c *gin.Context) {
	cid := getIDParam(c, "cardID")
	fid := getIDParam(c, "fileID")

	if err := h.repository.ValidateUID(cid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	co, err := h.repository.Create(cid, fid)

	if err != nil {
		log.Printf("failed insert record to covers table: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"cover": co})
}

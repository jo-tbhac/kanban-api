package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"local.packages/repository"
)

// BoardBackgroundImageHandler ...
type BoardBackgroundImageHandler struct {
	repository *repository.BoardBackgroundImageRepository
}

// NewBoardBackgroundImageHandler is constructor for BoardBackgroundImageHandler.
func NewBoardBackgroundImageHandler(r *repository.BoardBackgroundImageRepository) *BoardBackgroundImageHandler {
	return &BoardBackgroundImageHandler{
		repository: r,
	}
}

// CreateBoardBackgroundImage call a function that create a new record to board_background_images table.
// if creation was successful, returns status 201 and instance of BoardBackgroundImage as http response.
// if creation was failure, returns status 400 and error with messages.
func (h BoardBackgroundImageHandler) CreateBoardBackgroundImage(c *gin.Context) {
	bid := getIDParam(c, "boardID")
	iid := getIDParam(c, "backgroundImageID")

	if err := h.repository.ValidateUID(bid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id and current user")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	b, err := h.repository.Create(bid, iid)

	if err != nil {
		log.Printf("failed insert record to board_background_image table: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"board_background_image": b})
}

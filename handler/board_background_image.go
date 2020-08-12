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

// UpdateBoardBackgroundImage call a function that update a record for board_background_images table.
// if creation was successful, returns status 200 as http response.
// if creation was failure, returns status 400 and error with messages.
func (h BoardBackgroundImageHandler) UpdateBoardBackgroundImage(c *gin.Context) {
	bid := getIDParam(c, "boardID")
	iid := getIDParam(c, "backgroundImageID")

	b, err := h.repository.Find(bid, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id and current user")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Update(b, iid); err != nil {
		log.Printf("failed update for board_background_image tables record: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

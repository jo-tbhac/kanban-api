package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"local.packages/repository"
)

// FileHandler ...
type FileHandler struct {
	repository *repository.FileRepository
}

// NewFileHandler is constructor for FileHandler.
func NewFileHandler(r *repository.FileRepository) *FileHandler {
	return &FileHandler{repository: r}
}

func (h FileHandler) UploadFile(c *gin.Context) {
	cid := getIDParam(c, "cardID")

	if err := h.repository.ValidateUID(cid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	fh, _ := c.FormFile("file")

	f := h.repository.Upload(fh, cid)

	c.JSON(http.StatusOK, gin.H{"url": f.URL, "content-type": f.ContentType, "name": f.Name, "card_id": f.CardID})
}

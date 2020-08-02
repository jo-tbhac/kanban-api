package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"local.packages/repository"
	"local.packages/validator"
)

const maxFileSize = 8388608

// FileHandler ...
type FileHandler struct {
	repository *repository.FileRepository
}

// NewFileHandler is constructor for FileHandler.
func NewFileHandler(r *repository.FileRepository) *FileHandler {
	return &FileHandler{repository: r}
}

// UploadFile call a function that upload a file to storage and create a new record to files table.
// if creation was successful, returns status 201 and instance of File as http response.
// if creation was failure, returns status 400 and error with messages.
func (h FileHandler) UploadFile(c *gin.Context) {
	fh, _ := c.FormFile("file")

	if fh.Size > maxFileSize {
		log.Println("file size over 8MiB")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("file size over 8MiB")})
		return
	}

	cid := getIDParam(c, "cardID")

	if err := h.repository.ValidateUID(cid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	f := h.repository.Upload(fh, cid)

	if f == nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid request")})
		return
	}

	if err := h.repository.Create(f); err != nil {
		log.Printf("failed insert record to files table: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"file": f})
}

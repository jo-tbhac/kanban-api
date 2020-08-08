package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"local.packages/repository"
	"local.packages/validator"
)

type coverParams struct {
	FileID uint `json:"new_file_id"`
	CardID uint `json:"card_id"`
}

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

// UpdateCover call a function that update a record in covers table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
func (h CoverHandler) UpdateCover(c *gin.Context) {
	var p coverParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors(ErrorInvalidParameter)})
		return
	}

	co, err := h.repository.Find(p.CardID, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the cover")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Update(co, p.FileID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteCover call a function that delete a record from covers table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
func (h CoverHandler) DeleteCover(c *gin.Context) {
	cid := getIDParam(c, "cardID")

	co, err := h.repository.Find(cid, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the cover")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Delete(co); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

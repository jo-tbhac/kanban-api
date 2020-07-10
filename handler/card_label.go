package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"local.packages/repository"
	"local.packages/validator"
)

type cardLabelParams struct {
	LabelID uint `json:"label_id" binding:"required"`
}

// CardLabelHandler ...
type CardLabelHandler struct {
	repository *repository.CardLabelRepository
}

// NewCardLabelHandler is constructor for CardLabelHandler.
func NewCardLabelHandler(r *repository.CardLabelRepository) *CardLabelHandler {
	return &CardLabelHandler{repository: r}
}

// CreateCardLabel call a function that create a new record to card_labels table.
// if creation was successful, returns status 201 and instance of Label as http response.
// if creation was failure, returns status 400 and error with messages.
func (h *CardLabelHandler) CreateCardLabel(c *gin.Context) {
	var p cardLabelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	cid := getIDParam(c, "cardID")

	if err := h.repository.ValidateUID(p.LabelID, cid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the card or label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	l, err := h.repository.Create(p.LabelID, cid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"label": l})
}

// DeleteCardLabel call a function that delete a record from card_labels table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
func (h *CardLabelHandler) DeleteCardLabel(c *gin.Context) {
	cid := getIDParam(c, "cardID")
	lid := getIDParam(c, "labelID")

	cl, err := h.repository.Find(lid, cid, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the card or label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Delete(cl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

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

type CardLabelHandler struct {
	repository *repository.CardLabelRepository
}

func NewCardLabelHandler(r *repository.CardLabelRepository) *CardLabelHandler {
	return &CardLabelHandler{repository: r}
}

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

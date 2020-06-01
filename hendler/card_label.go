package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/repository"
	"github.com/jo-tbhac/kanban-api/validator"
)

type cardLabelParams struct {
	LabelID uint `json:"label_id" form:"label_id" binding:"required"`
}

type CardLabelHandler struct {
	repository repository.CardLabelRepository
}

func NewCardLabelHandler(r *repository.CardLabelRepository) *CardLabelHandler {
	return &CardLabelHandler{}
}

func (h *CardLabelHandler) createCardLabel(c *gin.Context) {
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

func (h *CardLabelHandler) deleteCardLabel(c *gin.Context) {
	var p cardLabelParams

	if err := c.ShouldBindQuery(&p); err != nil {
		log.Printf("fail to bind Query: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	cid := getIDParam(c, "cardID")

	cl, err := h.repository.Find(p.LabelID, cid, currentUserID(c))

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

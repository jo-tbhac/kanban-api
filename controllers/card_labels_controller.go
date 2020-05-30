package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"
)

type cardLabelParams struct {
	LabelID uint `json:"label_id" binding:"required"`
}

func createCardLabel(c *gin.Context) {
	var p cardLabelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	cl := models.CardLabel{
		LabelID: p.LabelID,
		CardID:  getIDParam(c, "cardID"),
	}

	if !cl.ValidateUID(currentUser(c).ID) {
		log.Println("uid does not match board.user_id associated with the card or label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid request")})
		return
	}

	l, err := cl.Create()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"label": l})
}

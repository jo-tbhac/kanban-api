package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"local.packages/repository"
	"local.packages/validator"
)

type cardParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CardHandler struct {
	repository *repository.CardRepository
}

func NewCardHandler(r *repository.CardRepository) *CardHandler {
	return &CardHandler{repository: r}
}

func (h CardHandler) CreateCard(c *gin.Context) {
	var p cardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	lid := getIDParam(c, "listID")

	if err := h.repository.ValidateUID(lid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	ca, err := h.repository.Create(p.Title, lid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"card": ca})
}

func (h CardHandler) UpdateCard(c *gin.Context) {
	id := getIDParam(c, "cardID")
	ca, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var p cardParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	switch c.Param("attribute") {
	case "title":
		if err := h.repository.UpdateTitle(ca, p.Title); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err})
			return
		}
	case "description":
		if err := h.repository.UpdateDescription(ca, p.Description); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"card": ca})
}

func (h CardHandler) DeleteCard(c *gin.Context) {
	id := getIDParam(c, "cardID")
	ca, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := h.repository.Delete(ca); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

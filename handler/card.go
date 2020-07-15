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

// CardHandler ...
type CardHandler struct {
	repository *repository.CardRepository
}

// NewCardHandler is constructor for CardHandler.
func NewCardHandler(r *repository.CardRepository) *CardHandler {
	return &CardHandler{repository: r}
}

// CreateCard call a function that create a new record to cards table.
// if creation was successful, returns status 201 and instance of Card as http response.
// if creation was failure, returns status 400 and error with messages.
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

// UpdateCard call a function that update a record in cards table.
// if update was successful, returns status 200 and updated instance of Card as http response.
// if update was failure, returns status 400 and error with messages.
func (h CardHandler) UpdateCard(c *gin.Context) {
	id := getIDParam(c, "cardID")
	ca, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
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

// UpdateCardIndex call a function that update cards order.
// if update was successful, returns status 200.
// if update was failure, returns status 400 and error with messages.
func (h CardHandler) UpdateCardIndex(c *gin.Context) {
	var ps []struct {
		ID     uint `json:"id"`
		Index  int  `json:"index"`
		ListID uint `json:"list_id"`
	}

	if err := c.ShouldBindJSON(&ps); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	if err := h.repository.UpdateIndex(ps); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteCard call a function that delete a record from cards table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
func (h CardHandler) DeleteCard(c *gin.Context) {
	id := getIDParam(c, "cardID")
	ca, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Delete(ca); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

// SearchCard returns status 200 and slice of Card instance as http response.
func (h CardHandler) SearchCard(c *gin.Context) {
	p := struct {
		Title   string `form:"title"`
		BoardID uint   `form:"board_id"`
	}{}

	if err := c.ShouldBindQuery(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	ids := h.repository.Search(p.BoardID, currentUserID(c), p.Title)

	c.JSON(http.StatusOK, gin.H{"card_ids": ids})
}

package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"local.packages/repository"
	"local.packages/validator"
)

type labelParams struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type LabelHandler struct {
	repository *repository.LabelRepository
}

func NewLabelHandler(r *repository.LabelRepository) *LabelHandler {
	return &LabelHandler{repository: r}
}

func (h LabelHandler) CreateLabel(c *gin.Context) {
	var p labelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	bid := getIDParam(c, "boardID")

	if err := h.repository.ValidateUID(bid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	l, err := h.repository.Create(p.Name, p.Color, bid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"label": l})
}

func (h LabelHandler) UpdateLabel(c *gin.Context) {
	id := getIDParam(c, "labelID")
	l, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	var p labelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	if err := h.repository.Update(l, p.Name, p.Color); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": l})
}

func (h LabelHandler) IndexLabel(c *gin.Context) {
	bid := getIDParam(c, "boardID")
	ls := h.repository.GetAll(bid, currentUserID(c))

	c.JSON(http.StatusOK, gin.H{"labels": ls})
}

func (h LabelHandler) DeleteLabel(c *gin.Context) {
	id := getIDParam(c, "labelID")
	l, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the label")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Delete(l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

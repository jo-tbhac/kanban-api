package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"local.packages/repository"
	"local.packages/validator"
)

type listParams struct {
	Name string `json:"name"`
}

type ListHandler struct {
	repository repository.ListRepository
}

func NewListHandler(r *repository.ListRepository) *ListHandler {
	return &ListHandler{}
}

func (h ListHandler) createList(c *gin.Context) {
	var p listParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	bid := getIDParam(c, "boardID")

	if err := h.repository.ValidateUID(bid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	l, err := h.repository.Create(p.Name, bid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"list": l})
}

func (h ListHandler) updateList(c *gin.Context) {
	id := getIDParam(c, "listID")
	l, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	var p listParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	if err := h.repository.Update(l, p.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": l})
}

func (h ListHandler) deleteList(c *gin.Context) {
	id := getIDParam(c, "listID")
	l, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Delete(l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}

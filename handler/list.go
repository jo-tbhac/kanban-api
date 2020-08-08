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

// ListHandler ...
type ListHandler struct {
	repository *repository.ListRepository
}

// NewListHandler is constructor for ListHandler.
func NewListHandler(r *repository.ListRepository) *ListHandler {
	return &ListHandler{repository: r}
}

// CreateList call a function that create a new record to lists table.
// if creation was successful, returns status 201 and instance of List as http response.
// if creation was failure, returns status 400 and error with messages.
func (h ListHandler) CreateList(c *gin.Context) {
	var p listParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors(ErrorInvalidParameter)})
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
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"list": l})
}

// UpdateList call a function that update a record in lists table.
// if update was successful, returns status 200 and updated instance of List as http response.
// if update was failure, returns status 400 and error with messages.
func (h ListHandler) UpdateList(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors(ErrorInvalidParameter)})
		return
	}

	if err := h.repository.Update(l, p.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"list": l})
}

// UpdateListIndex call a function that update lists order.
// if update was successful, returns status 200.
// if update was failure, returns status 400 and error with messages.
func (h ListHandler) UpdateListIndex(c *gin.Context) {
	var ps []struct {
		ID    uint
		Index int
	}
	if err := c.ShouldBindJSON(&ps); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors(ErrorInvalidParameter)})
		return
	}

	if err := h.repository.UpdateIndex(ps); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteList call a function that delete a record from lists table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
func (h ListHandler) DeleteList(c *gin.Context) {
	id := getIDParam(c, "listID")
	l, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Delete(l); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

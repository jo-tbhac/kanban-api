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

// LabelHandler ...
type LabelHandler struct {
	repository *repository.LabelRepository
}

// NewLabelHandler is constructor for LabelHandler.
func NewLabelHandler(r *repository.LabelRepository) *LabelHandler {
	return &LabelHandler{repository: r}
}

// CreateLabel call a function that create a new record to labels table.
// if creation was successful, returns status 201 and instance of Label as http response.
// if creation was failure, returns status 400 and error with messages.
func (h LabelHandler) CreateLabel(c *gin.Context) {
	var p labelParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors(ErrorInvalidParameter)})
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

// UpdateLabel call a function that update a record in labels table.
// if update was successful, returns status 200 and updated instance of Label as http response.
// if update was failure, returns status 400 and error with messages.
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
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors(ErrorInvalidParameter)})
		return
	}

	if err := h.repository.Update(l, p.Name, p.Color); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"label": l})
}

// IndexLabel returns status 200 and slice of Label instance as http response.
func (h LabelHandler) IndexLabel(c *gin.Context) {
	bid := getIDParam(c, "boardID")
	ls := h.repository.GetAll(bid, currentUserID(c))

	c.JSON(http.StatusOK, gin.H{"labels": ls})
}

// DeleteLabel call a function that delete a record from labels table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
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

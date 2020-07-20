package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"local.packages/repository"
	"local.packages/validator"
)

type checkListItemParams struct {
	Name  string `json:"name"`
	Check bool   `json:"check"`
}

// CheckListItemHandler ...
type CheckListItemHandler struct {
	repository *repository.CheckListItemRepository
}

// NewCheckListItemHandler is constructor for CardHandler.
func NewCheckListItemHandler(r *repository.CheckListItemRepository) *CheckListItemHandler {
	return &CheckListItemHandler{repository: r}
}

// CreateCheckListItem call a function that create a new record to check_list_items table.
// if creation was successful, returns status 201 and instance of CheckListItem as http response.
// if creation was failure, returns status 400 and error with messages.
func (h CheckListItemHandler) CreateCheckListItem(c *gin.Context) {
	var p checkListItemParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	cid := getIDParam(c, "checkListID")

	if err := h.repository.ValidateUID(cid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the check_list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	item, err := h.repository.Create(p.Name, cid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"check_list_item": item})
}

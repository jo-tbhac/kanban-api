package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"local.packages/repository"
	"local.packages/validator"
)

type checkListParams struct {
	Title string `json:"title"`
}

// CheckListHandler ...
type CheckListHandler struct {
	repository *repository.CheckListRepository
}

// NewCheckListHandler is constructor for CardHandler.
func NewCheckListHandler(r *repository.CheckListRepository) *CheckListHandler {
	return &CheckListHandler{repository: r}
}

// CreateCheckList call a function that create a new record to check_lists table.
// if creation was successful, returns status 201 and instance of CheckList as http response.
// if creation was failure, returns status 400 and error with messages.
func (h CheckListHandler) CreateCheckList(c *gin.Context) {
	var p checkListParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	cid := getIDParam(c, "cardID")

	if err := h.repository.ValidateUID(cid, currentUserID(c)); err != nil {
		log.Println("uid does not match board.user_id associated with the card")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	cl, err := h.repository.Create(p.Title, cid)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"check_list": cl})
}

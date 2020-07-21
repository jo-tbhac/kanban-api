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

// UpdateCheckList call a function that update a record in check_lists table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
func (h CheckListHandler) UpdateCheckList(c *gin.Context) {
	var p checkListParams

	if err := c.ShouldBindJSON(&p); err != nil {
		log.Printf("fail to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.NewValidationErrors("invalid parameters")})
		return
	}

	cid := getIDParam(c, "checkListID")

	cl, err := h.repository.Find(cid, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the check_list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Update(cl, p.Title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteCheckList call a function that delete a record from check_lists table.
// if deletion was successful, returns status 200.
// if deletion was failure, returns status 400 and errors with message.
func (h CheckListHandler) DeleteCheckList(c *gin.Context) {
	id := getIDParam(c, "checkListID")

	cl, err := h.repository.Find(id, currentUserID(c))

	if err != nil {
		log.Println("uid does not match board.user_id associated with the check_list")
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	if err := h.repository.Delete(cl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

// IndexCheckList returns status 200 and slice of CheckList instance as http response.
func (h CheckListHandler) IndexCheckList(c *gin.Context) {
	bid := getIDParam(c, "boardID")
	cs := h.repository.GetAll(bid, currentUserID(c))
	c.JSON(http.StatusOK, gin.H{"check_lists": cs})
}

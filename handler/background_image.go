package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"local.packages/repository"
)

// BackgroundImageHandler ...
type BackgroundImageHandler struct {
	repository *repository.BackgroundImageRepository
}

// NewBackgroundImageHandler is constructor for BackgroundImageHandler.
func NewBackgroundImageHandler(r *repository.BackgroundImageRepository) *BackgroundImageHandler {
	return &BackgroundImageHandler{
		repository: r,
	}
}

// IndexBackgroundImage returns status 200 and slice of BackgroundImage instance as http response.
func (h BackgroundImageHandler) IndexBackgroundImage(c *gin.Context) {
	bs := h.repository.GetAll()

	c.JSON(http.StatusOK, gin.H{"background_images": bs})
}

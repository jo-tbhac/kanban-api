package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"local.packages/validator"
)

type sessionParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateSession call a function that authenticate by request params.
// returns a session token if authentication was valid.
func (h UserHandler) CreateSession(c *gin.Context) {
	var p sessionParams

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.FormattedValidationError(err)})
		return
	}

	u, err := h.repository.SignIn(p.Email, p.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": u.RememberToken})
}

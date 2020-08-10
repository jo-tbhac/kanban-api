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
// returns access token and refresh token if authentication was valid.
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

	d := time.Until(u.ExpiresAt)

	c.JSON(
		http.StatusOK,
		gin.H{
			"access_token":  u.RememberToken,
			"refresh_token": u.RefreshToken,
			"expires_in":    d.Milliseconds,
		},
	)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"local.packages/utils"
	"local.packages/validator"
)

type sessionParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateSession call a function that authenticate by request params.
// returns access token, refresh token and expires if authentication was valid.
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

	c.JSON(
		http.StatusOK,
		gin.H{
			"email":         u.Email,
			"name":          u.Name,
			"access_token":  u.RememberToken,
			"refresh_token": u.RefreshToken,
			"expires_in":    utils.CalcExpiresIn(u.ExpiresAt),
		},
	)
}

// UpdateSession call a function that update access token and refresh token.
// returns access token, refresh token and expires if authentication was valid.
func (h UserHandler) UpdateSession(c *gin.Context) {
	at := c.Request.Header.Get("X-Auth-Token")

	p := struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validator.FormattedValidationError(err)})
		return
	}

	u, ok := h.repository.ValidateToken(at, p.RefreshToken)

	if !ok {
		c.JSON(http.StatusOK, gin.H{"ok": false})
		return
	}

	if err := h.repository.UpdateUserSession(u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"ok":            true,
			"name":          u.Name,
			"email":         u.Email,
			"access_token":  u.RememberToken,
			"refresh_token": u.RefreshToken,
			"expires_in":    utils.CalcExpiresIn(u.ExpiresAt),
		},
	)
}

// DeleteSession call a function that delete access token and refresh token.
func (h UserHandler) DeleteSession(c *gin.Context) {
	if err := h.repository.SignOut(currentUserID(c)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	c.Status(http.StatusOK)
}

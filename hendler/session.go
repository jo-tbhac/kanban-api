package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/validator"
)

type sessionParams struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h UserHandler) createSession(c *gin.Context) {
	var p sessionParams

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.NewValidationErrors("invalid parameters")})
		return
	}

	u, err := h.repository.SignIn(p.Email, p.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": u.RememberToken})
}

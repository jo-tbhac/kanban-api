package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/validator"
)

func (h UserHandler) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		if u, ok := h.repository.IsSignedIn(token); ok {
			c.Set("uid", u.ID)
			return
		}

		c.AbortWithStatus(401)
	}
}

func MapIDParamsToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, p := range c.Params {
			if ok, err := regexp.MatchString(`ID$`, p.Key); err != nil {
				log.Printf("fail to regexp.MatchString: %v", err)
				c.AbortWithStatus(500)
				return
			} else if !ok {
				continue
			}

			id, err := strconv.Atoi(c.Param(p.Key))

			if err != nil {
				log.Printf("fail to cast string to int: %v", err)
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"errors": validator.NewValidationErrors(fmt.Sprintf("%s must be an integer", p.Key))})
				return
			}

			c.Set(p.Key, uint(id))
		}
	}
}

func currentUserID(c *gin.Context) uint {
	return c.Keys["uid"].(uint)
}

func getIDParam(c *gin.Context, key string) uint {
	return c.Keys[key].(uint)
}

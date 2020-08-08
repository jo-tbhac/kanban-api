package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	"local.packages/validator"
)

// Authenticate call a function that validate a session token.
// map a login user id to context if authentication was valid.
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

// MapIDParamsToContext map URL params that has suffix `ID` to context.
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
					gin.H{"errors": validator.NewValidationErrors(fmt.Sprintf("%s"+ErrorMustBeAnInteger, p.Key))})
				return
			}

			c.Set(p.Key, uint(id))
		}
	}
}

// CORSMiddleware is a configure about CORS.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func currentUserID(c *gin.Context) uint {
	return c.Keys["uid"].(uint)
}

func getIDParam(c *gin.Context, key string) uint {
	return c.Keys[key].(uint)
}

package controllers

import (
	"fmt"

	"github.com/jo-tbhac/kanban-api/config"
	"github.com/jo-tbhac/kanban-api/models"

	"github.com/gin-gonic/gin"
)

func authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		var u models.User

		if u.IsSignedIn(token) {
			c.Set("user", u)
			return
		}

		c.AbortWithStatus(401)
	}
}

func StartServer() {
	r := gin.Default()

	authorized := r.Group("/", authenticate())

	r.POST("/users", CreateUser)
	r.POST("/sessions", CreateSession)

	authorized.POST("/boards", CreateBoard)
	authorized.GET("/boards", IndexBoard)

	r.Run(fmt.Sprintf(":%v", config.Config.Web.Port))
}

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

func CurrentUser(c *gin.Context) models.User {
	return c.Keys["user"].(models.User)
}

func StartServer() {
	r := gin.Default()

	authorized := r.Group("/", authenticate())

	r.POST("/users", CreateUser)
	r.POST("/sessions", CreateSession)

	authorized.POST("/boards", CreateBoard)
	authorized.PATCH("/boards", UpdateBoard)
	authorized.GET("/boards", IndexBoard)
	authorized.GET("/board", ShowBoard)
	authorized.DELETE("/board", DeleteBoard)

	authorized.POST("/labels", CreateLabel)
	authorized.PATCH("/labels", UpdateLabel)
	authorized.GET("/labels", IndexLabel)
	authorized.DELETE("/label", DeleteLabel)

	authorized.POST("/lists", CreateList)
	authorized.PATCH("/lists", UpdateList)
	authorized.DELETE("/list", DeleteList)

	r.Run(fmt.Sprintf(":%v", config.Config.Web.Port))
}

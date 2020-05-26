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

	r.POST("/users", createUser)
	r.POST("/sessions", createSession)

	authorized.POST("/boards", createBoard)
	authorized.PATCH("/boards", updateBoard)
	authorized.GET("/boards", indexBoard)
	authorized.GET("/board", showBoard)
	authorized.DELETE("/board", deleteBoard)

	authorized.POST("/labels", createLabel)
	authorized.PATCH("/labels", updateLabel)
	authorized.GET("/labels", indexLabel)
	authorized.DELETE("/label", deleteLabel)

	authorized.POST("/lists", createList)
	authorized.PATCH("/lists", updateList)
	authorized.DELETE("/list", deleteList)

	authorized.POST("/cards", createCard)
	authorized.PATCH("/cards", updateCard)

	r.Run(fmt.Sprintf(":%v", config.Config.Web.Port))
}

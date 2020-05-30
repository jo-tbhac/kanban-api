package controllers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/jo-tbhac/kanban-api/config"
	"github.com/jo-tbhac/kanban-api/models"
	"github.com/jo-tbhac/kanban-api/validator"

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

func mapIDParamsToContext() gin.HandlerFunc {
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

func currentUser(c *gin.Context) models.User {
	return c.Keys["user"].(models.User)
}

func getIDParam(c *gin.Context, key string) uint {
	return c.Keys[key].(uint)
}

func StartServer() {
	r := gin.Default()

	r.Use(mapIDParamsToContext())

	authorized := r.Group("/", authenticate())

	r.POST("/user", createUser)
	r.POST("/session", createSession)

	authorized.POST("/board", createBoard)
	authorized.GET("/boards", indexBoard)
	authorized.GET("/board/:boardID", showBoard)
	authorized.PATCH("/board/:boardID", updateBoard)
	authorized.DELETE("/board/:boardID", deleteBoard)

	authorized.POST("/board/:boardID/label", createLabel)
	authorized.GET("/board/:boardID/labels", indexLabel)
	authorized.PATCH("/label/:labelID", updateLabel)
	authorized.DELETE("/label/:labelID", deleteLabel)

	authorized.POST("/board/:boardID/list", createList)
	authorized.PATCH("/list/:listID", updateList)
	authorized.DELETE("/list/:listID", deleteList)

	authorized.POST("/list/:listID/card", createCard)
	authorized.PATCH("/card/:cardID", updateCard)
	authorized.DELETE("/card/:cardID", deleteCard)

	authorized.POST("/card/:cardID/card_label", createCardLabel)
	authorized.DELETE("/card/:cardID/card_label", deleteCardLabel)

	r.Run(fmt.Sprintf(":%v", config.Config.Web.Port))
}

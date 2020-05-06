package controllers

import (
	"fmt"

	"github.com/jo-tbhac/kanban-api/config"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	r.POST("/users", CreateUser)

	r.Run(fmt.Sprintf(":%v", config.Config.Web.Port))
}

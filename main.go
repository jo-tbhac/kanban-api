package main

import (
	"fmt"

	"github.com/jo-tbhac/kanban-api/config"
	"github.com/jo-tbhac/kanban-api/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db := db.Get()
	defer db.Close()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(fmt.Sprintf(":%v", config.Config.Web.Port))
}

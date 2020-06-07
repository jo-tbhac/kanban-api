package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"local.packages/config"
	"local.packages/db"
	"local.packages/handler"
	"local.packages/migration"
	"local.packages/repository"
)

var (
	userHandler      *handler.UserHandler
	boardHandler     *handler.BoardHandler
	labelHandler     *handler.LabelHandler
	listHandler      *handler.ListHandler
	cardHandler      *handler.CardHandler
	cardLabelHandler *handler.CardLabelHandler
)

func main() {
	db := db.Get()
	defer db.Close()

	if os.Getenv("GIN_MODE") == "debug" {
		db.LogMode(true)
	}

	userHandler = handler.NewUserHandler(repository.NewUserRepository(db))
	boardHandler = handler.NewBoardHandler(repository.NewBoardRepository(db))
	labelHandler = handler.NewLabelHandler(repository.NewLabelRepository(db))
	listHandler = handler.NewListHandler(repository.NewListRepository(db))
	cardHandler = handler.NewCardHandler(repository.NewCardRepository(db))
	cardLabelHandler = handler.NewCardLabelHandler(repository.NewCardLabelRepository(db))

	migration.Migrate()
	startServer()
}

func startServer() {
	r := gin.Default()

	r.Use(handler.MapIDParamsToContext())

	authorized := r.Group("/", userHandler.Authenticate())

	r.POST("/user", userHandler.CreateUser)
	r.POST("/session", userHandler.CreateSession)

	authorized.POST("/board", boardHandler.CreateBoard)
	authorized.GET("/boards", boardHandler.IndexBoard)
	authorized.GET("/board/:boardID", boardHandler.ShowBoard)
	authorized.PATCH("/board/:boardID", boardHandler.UpdateBoard)
	authorized.DELETE("/board/:boardID", boardHandler.DeleteBoard)

	authorized.POST("/board/:boardID/label", labelHandler.CreateLabel)
	authorized.GET("/board/:boardID/labels", labelHandler.IndexLabel)
	authorized.PATCH("/label/:labelID", labelHandler.UpdateLabel)
	authorized.DELETE("/label/:labelID", labelHandler.DeleteLabel)

	authorized.POST("/board/:boardID/list", listHandler.CreateList)
	authorized.PATCH("/list/:listID", listHandler.UpdateList)
	authorized.DELETE("/list/:listID", listHandler.DeleteList)

	authorized.POST("/list/:listID/card", cardHandler.CreateCard)
	authorized.PATCH("/card/:cardID/:attribute", cardHandler.UpdateCard)
	authorized.DELETE("/card/:cardID", cardHandler.DeleteCard)

	authorized.POST("/card/:cardID/card_label", cardLabelHandler.CreateCardLabel)
	authorized.DELETE("/card/:cardID/card_label/:labelID", cardLabelHandler.DeleteCardLabel)

	r.Run(fmt.Sprintf(":%v", config.Config.Web.Port))
}

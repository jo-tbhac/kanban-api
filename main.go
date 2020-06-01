package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jo-tbhac/kanban-api/config"
	"github.com/jo-tbhac/kanban-api/db"
	"github.com/jo-tbhac/kanban-api/handler"
	"github.com/jo-tbhac/kanban-api/migration"
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

	authorized := r.Group("/", handler.Authenticate())

	r.POST("/user", userHandler.createUser)
	r.POST("/session", userHandler.createSession)

	authorized.POST("/board", boardHandler.createBoard)
	authorized.GET("/boards", boardHandler.indexBoard)
	authorized.GET("/board/:boardID", boardHandler.showBoard)
	authorized.PATCH("/board/:boardID", boardHandler.updateBoard)
	authorized.DELETE("/board/:boardID", boardHandler.deleteBoard)

	authorized.POST("/board/:boardID/label", labelHandler.createLabel)
	authorized.GET("/board/:boardID/labels", labelHandler.indexLabel)
	authorized.PATCH("/label/:labelID", labelHandler.updateLabel)
	authorized.DELETE("/label/:labelID", labelHandler.deleteLabel)

	authorized.POST("/board/:boardID/list", listHandler.createList)
	authorized.PATCH("/list/:listID", listHandler.updateList)
	authorized.DELETE("/list/:listID", listHandler.deleteList)

	authorized.POST("/list/:listID/card", cardHandler.createCard)
	authorized.PATCH("/card/:cardID", cardHandler.updateCard)
	authorized.DELETE("/card/:cardID", cardHandler.deleteCard)

	authorized.POST("/card/:cardID/card_label", cardLabelHandler.createCardLabel)
	authorized.DELETE("/card/:cardID/card_label", cardLabelHandler.deleteCardLabel)

	r.Run(fmt.Sprintf(":%v", config.Config.Web.Port))
}

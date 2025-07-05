package api

import (
	"lesta-battleship/server-core/internal/game-core/event"
	"lesta-battleship/server-core/internal/game-core/ws"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, dispatcher *event.MatchEventDispatcher) {
	r.POST("/start-match", StartMatch)
	r.GET("/ws", func(c *gin.Context) {
		ws.WebSocketHandler(c, dispatcher)
	})
}

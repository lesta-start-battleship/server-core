package api

import (
	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/ws"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, dispatcher *event.MatchEventDispatcher) {
	r.POST("/api/v1/start-match", StartMatch)
	r.GET("/ws", func(c *gin.Context) {
		ws.WebSocketHandler(c, dispatcher)
	})
}

package websocket

import (
	"lesta-battleship/server-core/internal/app/multiplayer"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router gin.IRouter, websocketServer *WebsocketServer, engine *multiplayer.Engine) {
	randomHandler := HandleGetJoinRandom(websocketServer, engine)
	rankedHandler := HandleGetJoinRanked(websocketServer, engine)
	customHandler := HandleGetJoinCustom(websocketServer, engine)

	router.GET("/join/random", randomHandler)
	router.GET("/join/ranked", rankedHandler)
	router.GET("/join/custom", customHandler)
}

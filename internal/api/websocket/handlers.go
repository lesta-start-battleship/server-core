package websocket

import (
	"lesta-battleship/server-core/internal/app/multiplayer"
	"lesta-battleship/server-core/internal/app/multiplayer/actors/matchmakers"

	"github.com/gin-gonic/gin"
)

func HandleGetJoinRandom(websocketServer *WebsocketServer, engine *multiplayer.Engine) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			interfacer, _ := websocketServer.Connect(c.Writer, c.Request)
			go interfacer.ReadPump()
			go interfacer.WritePump()

			player := engine.CreatePlayer(interfacer)
			engine.SendToMatchmaking(player, matchmakers.RandomMatch)
		},
	)
}

func HandleGetJoinRanked(websocketServer *WebsocketServer, engine *multiplayer.Engine) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			interfacer, _ := websocketServer.Connect(c.Writer, c.Request)
			go interfacer.ReadPump()
			go interfacer.WritePump()

			player := engine.CreatePlayer(interfacer)
			engine.SendToMatchmaking(player, matchmakers.RankedMatch)
		},
	)
}

func HandleGetJoinCustom(websocketServer *WebsocketServer, engine *multiplayer.Engine) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			interfacer, _ := websocketServer.Connect(c.Writer, c.Request)
			go interfacer.ReadPump()
			go interfacer.WritePump()

			player := engine.CreatePlayer(interfacer)
			engine.SendToMatchmaking(player, matchmakers.CustomMatch)
		},
	)
}

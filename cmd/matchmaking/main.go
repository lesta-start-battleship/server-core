package main

import (
	"lesta-battleship/server-core/internal/matchmaking/api/websocket"
	"lesta-battleship/server-core/internal/matchmaking/app/multiplayer"
	"lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/matchmakers"
	"lesta-battleship/server-core/internal/matchmaking/infra/registries"

	"github.com/gin-gonic/gin"
)

func main() {
	matchmakerRegistry := registries.NewMatchmakerRegistry()
	roomRegistry := registries.NewRoomRegistry()
	playerRegistry := registries.NewPlayerRegistry()

	engine := multiplayer.NewEngine(matchmakerRegistry, roomRegistry, playerRegistry)

	engine.CreateHub()

	engine.CreateMatchmaker(matchmakers.RandomMatch)
	engine.CreateMatchmaker(matchmakers.RankedMatch)
	engine.CreateMatchmaker(matchmakers.CustomMatch)

	websocketServer := websocket.NewWebsocketServer()

	router := gin.Default()

	websocket.SetupRouter(router, websocketServer, engine)

	router.Run(":8080")
}

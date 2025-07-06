package api

import (
	"lesta-battleship/server-core/internal/game-core/game"
	"lesta-battleship/server-core/internal/game-core/items"
	"lesta-battleship/server-core/internal/game-core/match"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartMatch(c *gin.Context) {
	var payload struct {
		RoomID  string `json:"room_id"`
		Player1 string `json:"player1"`
		Player2 string `json:"player2"`
		Mode    string `json:"mode"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if items, err := items.GetItemsInfo(payload.Player1, payload.Player2); err != nil {
		log.Printf("пизда рулям")
	}

	player1State := game.NewGameState()
	player2State := game.NewGameState()

	room := &match.GameRoom{
		RoomID: payload.RoomID,
		Mode:   payload.Mode,
		Player1: &match.PlayerConn{
			ID: payload.Player1,
			States: &game.States{
				PlayerState: player1State,
				EnemyState:  player2State,
			},
			Items: items.ItemsPlayer1,
		},
		Player2: &match.PlayerConn{
			ID: payload.Player2,
			States: &game.States{
				PlayerState: player2State,
				EnemyState:  player1State,
			},
			Items: items.ItemsPlayer2,
		},
		Status:    "waiting",
		CreatedAt: time.Now(),
		Items:     items.Items,
	}

	match.Rooms.Store(payload.RoomID, room)
	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

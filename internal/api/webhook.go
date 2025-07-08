package api

import (
	"log"
	"net/http"
	"time"

	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/items"
	"github.com/lesta-battleship/server-core/internal/match"

	"github.com/gin-gonic/gin"
)

func StartMatch(c *gin.Context) {
	var payload struct {
		RoomID     string `json:"room_id"`
		Player1    string `json:"player1"`
		Player2    string `json:"player2"`
		Mode       string `json:"mode"`
		GuildWarID string `json:"guild_war_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// testoviy token (admina), ispolzuyem inventar odnogo usera na 2
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoiYWRtaW4ifQ.JzhxV6sJhyTWgr4F-_EeDHg3-urRQiZUWYU9EvMZNHU"

	itemsPlayer1, err := items.GetUserItems(token)
	if err != nil {
		log.Printf("error fetching number of items for player1 %v", err)
	}

	itemsPlayer2, err := items.GetUserItems(token)
	if err != nil {
		log.Printf("error fetching number of items for player2 %v", err)
	}

	allItems, err := items.GetAllItems()
	if err != nil {
		log.Printf("error fetching all items %v", err)
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
			Items:     itemsPlayer1,
			ItemUsage: make(map[items.ItemID]*items.ItemUsageData),
			MoveCount: 0,
			ChessFigureCount: 0, // sobitiye tolko dlya chess
		},
		Player2: &match.PlayerConn{
			ID: payload.Player2,
			States: &game.States{
				PlayerState: player2State,
				EnemyState:  player1State,
			},
			Items:     itemsPlayer2,
			ItemUsage: make(map[items.ItemID]*items.ItemUsageData),
			MoveCount: 0,
			ChessFigureCount: 0, // sobitiye tolko dlya chess
		},
		Status:     "waiting",
		CreatedAt:  time.Now(),
		Items:      allItems,
		GuildWarID: payload.GuildWarID,
	}

	match.Rooms.Store(payload.RoomID, room)
	c.JSON(http.StatusOK, gin.H{"status": "created"})
}

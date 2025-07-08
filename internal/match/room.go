package match

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/items"

	"github.com/gorilla/websocket"
)

type PlayerConn struct {
	ID                string
	Ready             bool
	States            *game.States // ВАЖНО чтобы при создании GameRoom было создано всего 2 GameState, а не 4
	Conn              *websocket.Conn
	Items             map[items.ItemID]int                  // хранит количество предметов у юзера
	ItemUsage         map[items.ItemID]*items.ItemUsageData // хранит сколько раз использовался предмет, тоже для cd и limita
	MoveCount         int                                   // сколько раз игрок ходил (используется для cd)
	ChessFigureCount  int                                   // нужно чтобы использование шахматных фигур не было > 2 (sobitiye tolko dlya chess)
	LastSubmarineTurn int
	Disconnected      bool
	ReconnectTimer    *time.Timer
}

type GameRoom struct {
	RoomID     string
	GuildWarID string
	Mode       string
	Player1    *PlayerConn
	Player2    *PlayerConn
	Status     string // waiting, ready, playing, ended
	Turn       string // player ID
	WinnerID   string
	Mutex      sync.Mutex
	CreatedAt  time.Time
	Items      map[items.ItemID]*items.Item // хранит артефакты доступные в игре
}

var Rooms sync.Map

func (p *PlayerConn) WriteMessage(msgType int, data []byte) error {
	if p.Conn != nil {
		return p.Conn.WriteMessage(msgType, data)
	}
	return nil
}

func (r *GameRoom) DeclareVictory(winnerID string, dispatcher *event.MatchEventDispatcher) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if r.Status == "ended" || r.WinnerID != "" {
		return
	}

	r.WinnerID = winnerID
	r.Status = "ended"

	var winner, loser *PlayerConn
	if r.Player1.ID == winnerID {
		winner = r.Player1
		loser = r.Player2
	} else {
		winner = r.Player2
		loser = r.Player1
	}

	log.Printf("[MATCH] Player %s declared as winner in room %s (by technical reason)", winnerID, r.RoomID)

	matchResult := event.MatchResult{
		WinnerID:  winner.ID,
		LoserID:   loser.ID,
		MatchID:   r.RoomID,
		MatchDate: r.CreatedAt,
		MatchType: r.Mode,
	}

	switch r.Mode {
	case "ranked":
		matchResult.Experience = &event.Experience{WinnerGain: 30, LoserGain: 15}
		matchResult.Rating = &event.Rating{WinnerGain: 30, LoserGain: -15}
	case "random":
		matchResult.Experience = &event.Experience{WinnerGain: 30, LoserGain: 15}
	case "guild_war_match":
		matchResult.WarID = r.GuildWarID
		matchResult.Experience = &event.Experience{WinnerGain: 30, LoserGain: 15}
	case "custom":
	}

	if dispatcher != nil {
		if err := dispatcher.DispatchMatchResult(matchResult); err != nil {
			log.Printf("[MATCH] Failed to dispatch result: %v", err)
		}
	}

	endPayload := map[string]any{
		"event": "game_end",
		"data": map[string]string{
			"winner": winnerID,
		},
	}
	data, _ := json.Marshal(endPayload)
	for _, p := range []*PlayerConn{r.Player1, r.Player2} {
		if p.Conn != nil {
			_ = p.Conn.WriteMessage(websocket.TextMessage, data)
			_ = p.Conn.Close()
		}
	}

	Rooms.Delete(r.RoomID)
}

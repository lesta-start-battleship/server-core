package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game-core/match"

	"github.com/gorilla/websocket"
)

func HandleReady(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, _ EventInput) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if player.States.PlayerState.NumShips < 10 {
		err := errors.New("you must place 10 ships before ready")
		Send(conn, "ready_error", err.Error())
		return err
	}

	player.Ready = true
	allReady := room.Player1.Ready && room.Player2.Ready
	shouldStart := false

	if allReady && room.Status == "waiting" {
		room.Status = "playing"
		room.Turn = room.Player1.ID
		shouldStart = true
	}

	Send(conn, "ready_confirmed", map[string]any{"all_ready": allReady})

	if shouldStart {
		Broadcast(room, "game_start", map[string]any{"first_turn": room.Turn})
	}

	return nil
}

package handlers

import (
	// "errors"
	"github.com/lesta-battleship/server-core/internal/match"

	"github.com/gorilla/websocket"
)

func HandleReady(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	// TODO: после тестов убрать комменты

	// if player.States.PlayerState.NumShips < 10 {
	// 	err := errors.New("you must place 10 ships before ready")
	// 	SendError(conn, err.Error())
	// 	return err
	// }

	player.Ready = true
	allReady := room.Player1.Ready && room.Player2.Ready
	shouldStart := false

	if allReady && room.Status == "waiting" {
		room.Status = "playing"
		room.Turn = room.Player1.ID
		shouldStart = true
	}

	SendSuccess(conn, EventReadyConfirmed, ReadyConfirmedResponse{
		AllReady: allReady,
	})

	if shouldStart {
		Broadcast(room, EventGameStart, GameStartResponse{
			FirstTurn: room.Turn,
		})
	}

	return nil
}

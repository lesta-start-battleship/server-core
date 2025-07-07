package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/lesta-battleship/server-core/internal/match"
	"github.com/lesta-battleship/server-core/internal/wsiface"
)

func SendSuccess(conn *websocket.Conn, event string, data any) error {
	resp := wsiface.WSResponse{
		Event: event,
		Data:  data,
	}
	if err := conn.WriteJSON(resp); err != nil {
		return fmt.Errorf("send success failed: %w", err)
	}
	return nil
}

func SendError(conn *websocket.Conn, message string) error {
	resp := wsiface.WSResponse{
		Event: wsiface.EventError,
		Data:  message,
	}
	if err := conn.WriteJSON(resp); err != nil {
		return fmt.Errorf("send error failed: %w", err)
	}
	return nil
}

func Broadcast(room *match.GameRoom, event string, data any) error {
	resp := wsiface.WSResponse{
		Event: event,
		Data:  data,
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("broadcast marshal error: %w", err)
	}

	var err1, err2 error

	if room.Player1.Conn != nil {
		err1 = room.Player1.Conn.WriteMessage(websocket.TextMessage, raw)
	}
	if room.Player2.Conn != nil {
		err2 = room.Player2.Conn.WriteMessage(websocket.TextMessage, raw)
	}

	if err1 != nil {
		return fmt.Errorf("player1 write error: %w", err1)
	}
	if err2 != nil {
		return fmt.Errorf("player2 write error: %w", err2)
	}
	return nil
}

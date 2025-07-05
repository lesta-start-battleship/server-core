package handlers

import (
	"encoding/json"
	"lesta-battleship/server-core/internal/game-core/match"
	"log"

	"github.com/gorilla/websocket"
)

func Send(conn *websocket.Conn, event string, data any) {
	err := conn.WriteJSON(map[string]any{
		"event": event,
		"data":  data,
	})
	if err != nil {
		log.Println("[WS] Send failed:", err)
	}
}

func Broadcast(room *match.GameRoom, event string, data any) {
	msg := map[string]any{
		"event": event,
		"data":  data,
	}
	raw, _ := json.Marshal(msg)

	if room.Player1.Conn != nil {
		room.Player1.Conn.WriteMessage(websocket.TextMessage, raw)
	}
	if room.Player2.Conn != nil {
		room.Player2.Conn.WriteMessage(websocket.TextMessage, raw)
	}
}

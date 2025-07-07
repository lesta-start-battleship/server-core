package ws

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/match"
	"github.com/lesta-battleship/server-core/internal/ws/handlers"

	"github.com/lesta-battleship/server-core/internal/wsiface"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(c *gin.Context, dispatcher *event.MatchEventDispatcher) {
	roomID := c.Query("room_id")
	playerID := c.Query("player_id")

	rawRoom, ok := match.Rooms.Load(roomID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	room := rawRoom.(*match.GameRoom)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("[WS] Upgrade error:", err)
		return
	}
	defer conn.Close()

	var player *match.PlayerConn
	if room.Player1.ID == playerID {
		player = room.Player1
		room.Player1.Conn = conn
	} else if room.Player2.ID == playerID {
		player = room.Player2
		room.Player2.Conn = conn
	} else {
		log.Println("[WS] Invalid playerID:", playerID)
		conn.Close()
		return
	}
	player.Conn = conn

	log.Printf("[WS] Player %s connected to room %s\n", playerID, roomID)

	ctx := &wsiface.Context{
		Conn:       conn,
		Player:     player,
		Room:       room,
		Dispatcher: dispatcher,
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[WS] Read error:", err)
			break
		}

		var input wsiface.WSInput
		decoder := json.NewDecoder(bytes.NewReader(msg))
		decoder.DisallowUnknownFields() 

		if err = decoder.Decode(&input); err != nil {
			log.Println("[WS] JSON decode error (invalid fields?):", err)
			handlers.SendError(conn, "invalid input format: "+err.Error())
			continue
		}

		handler, ok := handlers.GetHandler(input.Event)
		if !ok {
			handlers.SendError(conn, "unknown event")
			continue
		}

		log.Printf("[WS] Event received from %s: %s\n", playerID, input.Event)
		if err := handler.Handle(input, ctx); err != nil {
			log.Printf("[WS] Handler error (%s): %v\n", input.Event, err)
		}
	}
}

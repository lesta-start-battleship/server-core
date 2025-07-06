package ws

import (
	"encoding/json"
	"lesta-battleship/server-core/internal/game-core/event"
	"lesta-battleship/server-core/internal/game-core/game"
	"lesta-battleship/server-core/internal/game-core/match"
	"lesta-battleship/server-core/internal/game-core/ws/handlers"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("[WS] Upgrade error:", err)
		return
	}

	room := rawRoom.(*match.GameRoom)

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

	log.Printf("[WS] Player %s connected to room %s\n", playerID, roomID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("[WS] Read error:", err)
			break
		}

		var input struct {
			Event string    `json:"event"`
			Ship  game.Ship `json:"ship"`
			X     int       `json:"x"`
			Y     int       `json:"y"`
			ItemID int 		`json:"itemid"`
		}
		if err = json.Unmarshal(msg, &input); err != nil {
			log.Println("[WS] JSON unmarshal error:", err)
			continue
		}

		log.Printf("[WS] Event received from %s: %s\n", playerID, input.Event)

		switch input.Event {
		case "place_ship":
			if err := handlers.HandlePlaceShip(room, player, conn, input); err != nil {
				log.Printf("[WS] PlaceShip error: %v", err)
				continue
			}

		case "remove_ship":
			if err := handlers.HandleRemoveShip(room, player, conn, input); err != nil {
				log.Printf("[WS] RemoveShip error: %v", err)
				continue
			}

		case "ready":
			if err := handlers.HandleReady(room, player, conn, input); err != nil {
				log.Printf("[WS] Ready error: %v", err)
				continue
			}

		case "shoot":
			if err := handlers.HandleFire(room, player, conn, input, dispatcher); err != nil {
				log.Printf("[WS] Shoot error: %v", err)
				continue
			}

		case "heal_ship": // TODO: прямого доступа снаружи к данной функциональности быть не должно
			if err := handlers.HandleHealShip(room, player, conn, input); err != nil {
				log.Printf("[WS] HealShip error: %v", err)
				continue
			}

		case "open_cell": // TODO: прямого доступа снаружи к данной функциональности быть не должно
			if err := handlers.HandleOpenCell(room, player, conn, input); err != nil {
				log.Printf("[WS] OpenCell error: %v", err)
				continue
			}

		case "item":
			if err := handlers.HandleItem(room, player, conn, input); err != nil {
				log.Printf("[WS] Item error: %v", err)
				continue
			}
			
		default:
			handlers.Send(conn, "error", "unknown event")
		}
	}
}

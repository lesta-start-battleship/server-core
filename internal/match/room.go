package match

import (
	"sync"
	"time"

	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/items"

	"github.com/gorilla/websocket"
)

type PlayerConn struct {
	ID     string
	Ready  bool
	States *game.States // ВАЖНО чтобы при создании GameRoom было создано всего 2 GameState, а не 4
	Conn   *websocket.Conn
	Items  map[items.ItemID]int // хранит количество предметов у юзера
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

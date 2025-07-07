package wsiface

import (
	"github.com/gorilla/websocket"
	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/match"
)

type WSEventHandler interface {
	EventName() string
	Handle(input any, ctx *Context) error
}

type Context struct {
	Conn       *websocket.Conn
	Player     *match.PlayerConn
	Room       *match.GameRoom
	Dispatcher *event.MatchEventDispatcher
}

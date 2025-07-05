package websocket

import (
	"net/http"

	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors"
	"github.com/lesta-battleship/server-core/internal/matchmaking/infra"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	upgrader websocket.Upgrader
}

func NewWebsocketServer() *WebsocketServer {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &WebsocketServer{
		upgrader: upgrader,
	}
}

func (s *WebsocketServer) Connect(w http.ResponseWriter, r *http.Request) (actors.ClientInterfacer, error) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	websocketInterfacer := infra.NewWebsocketInterfacer(infra.GenerateId(), conn)

	return websocketInterfacer, nil
}

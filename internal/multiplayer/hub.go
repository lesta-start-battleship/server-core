package multiplayer

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	players    map[*Player]struct{}
	upgrader   websocket.Upgrader
	matchmaker *Matchmaker

	connectPlayerChan    chan *Player
	disconnectPlayerChan chan *Player
}

func NewHub() *Hub {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	matchmaker := NewMatchmaker()
	go matchmaker.Run()

	return &Hub{
		players: make(map[*Player]struct{}),

		upgrader:   upgrader,
		matchmaker: matchmaker,

		connectPlayerChan:    make(chan *Player),
		disconnectPlayerChan: make(chan *Player),
	}
}

func (h *Hub) ConnectPlayer(w http.ResponseWriter, r *http.Request) error {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := NewPlayer(generateId(), conn, h)
	h.players[client] = struct{}{}

	h.connectPlayerChan <- client

	go client.ReadFromSocket()
	go client.WriteToSocket()

	return nil
}

func (h *Hub) Run() {
	defer h.Close()

	for {
		select {
		case player := <-h.connectPlayerChan:
			h.Connect(player)
		case player := <-h.disconnectPlayerChan:
			h.Disconnect(player)
		}
	}
}

func (h *Hub) Connect(player *Player) {
	h.matchmaker.randomPoolChan <- player

	log.Printf("Hub: Received player %q", player.id)
}

func (h *Hub) Disconnect(player *Player) {
	delete(h.players, player)

	log.Printf("Hub: Disconnected player %q", player.id)
}

// TODO: Probably unsafe
func (h *Hub) Close() {
	for player := range h.players {
		delete(h.players, player)
		player.Close()
	}

	h.matchmaker.Close()
	h.matchmaker = nil

	if _, ok := <-h.connectPlayerChan; !ok {
		close(h.connectPlayerChan)
	}
	if _, ok := <-h.disconnectPlayerChan; !ok {
		close(h.disconnectPlayerChan)
	}

	log.Println("Hub: Closed")
}

// TODO: Create proper ID generator
func generateId() string {
	return rand.Text()
}

package multiplayer

import (
	"encoding/json"
	"lesta-battleship/server-core/pkg/packets"
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	id   string
	conn *websocket.Conn
	hub  *Hub
	room *Room

	messageChan chan packets.Packet
}

func NewPlayer(id string, conn *websocket.Conn, hub *Hub) *Player {
	return &Player{
		id:          id,
		conn:        conn,
		hub:         hub,
		messageChan: make(chan packets.Packet, 256),
	}
}

func (p *Player) ConnectedToSocket() bool {
	return p.conn != nil
}

func (p *Player) ConnectedToRoom() bool {
	return p.room != nil
}

func (p *Player) ConnectToRoom(room *Room) error {
	if p.room != nil {
		return ErrAlreadyConnectedToRoom
	}
	p.room = room

	log.Printf("Player %q: Connected to room %q", p.id, room.id)

	return nil
}

func (p *Player) DisconnectFromRoom() error {
	if p.room == nil {
		return ErrNotConnectedToRoom
	}
	p.room = nil

	log.Printf("Player %q: Disconnected from room %q", p.id, p.room.id)

	return nil
}

func (p *Player) GetMessage(packet packets.Packet) {
	p.messageChan <- packet

	log.Printf("Player %q: Received message", p.id)
}

func (p *Player) SendMessage(packet packets.Packet) {
	p.room.broadcastChan <- packet

	log.Printf("Player %q: Sent message", p.id)
}

func (p *Player) ReadFromSocket() {
	defer p.Close()

	for {
		msgType, msg, err := p.conn.ReadMessage()
		if err != nil {
			log.Printf("Player: %s", err)
			break
		}
		if msgType != websocket.TextMessage {
			log.Printf("Player: %s", err)
			break
		}

		packet := packets.Packet{}
		if err := json.Unmarshal(msg, &packet); err != nil {
			log.Printf("Player: %s", err)
			break
		}

		// TODO: Ensure to add UserID from Group 3
		if packet.SenderId == "" {
			packet.SenderId = p.id
		}

		p.SendMessage(packet)
	}
}

func (c *Player) WriteToSocket() {
	defer c.Close()

	for packet := range c.messageChan {
		err := c.conn.WriteJSON(packet)
		if err != nil {
			log.Printf("Player: %s", err)
			break
		}
	}
}

func (c *Player) Close() {
	c.room.disconnectChan <- c
	c.room = nil

	c.hub.disconnectPlayerChan <- c
	c.hub = nil

	if _, ok := <-c.messageChan; !ok {
		close(c.messageChan)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Player %q: %v", c.id, err)
	}

	log.Printf("Player %q: Disconnected", c.id)
}

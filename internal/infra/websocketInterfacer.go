package infra

import (
	"lesta-battleship/server-core/internal/app/multiplayer/actors"
	"lesta-battleship/server-core/pkg/packets"
	"log"

	"github.com/gorilla/websocket"
)

type WebsocketInterfacer struct {
	id string

	conn   *websocket.Conn
	player actors.Actor

	packetChan chan packets.Packet
}

func NewWebsocketInterfacer(id string, conn *websocket.Conn) *WebsocketInterfacer {
	return &WebsocketInterfacer{
		id: id,

		conn:   conn,
		player: nil,

		packetChan: make(chan packets.Packet),
	}
}

func (i *WebsocketInterfacer) Id() string {
	return i.id
}

func (p *WebsocketInterfacer) ConnectTo(player actors.Actor) {
	p.player = player
}

func (p *WebsocketInterfacer) GetPacket(senderId string, packet packets.Packet) {
	if p.id == senderId {
		p.player.GetPacket(senderId, packet)
	} else {
		p.packetChan <- packet
	}

	log.Printf("WebsocketInterfacer %q: Got packet %T from %q", p.id, packet.Body, senderId)
}

func (p *WebsocketInterfacer) ReadPump() {
	defer func() {
		if _, ok := <-p.packetChan; !ok {
			close(p.packetChan)
		}
		p.Stop()
	}()

	for {
		packet := packets.Packet{}
		if err := p.conn.ReadJSON(&packet); err != nil {
			log.Printf("WebsocketInterfacer %q: %v", p.id, err)

			break
		}

		if packet.SenderId == "" {
			packet.SenderId = p.id
		}

		p.GetPacket(packet.SenderId, packet)
	}
}

func (p *WebsocketInterfacer) WritePump() {
	defer func() {
		if _, ok := <-p.packetChan; !ok {
			close(p.packetChan)
		}
		p.Stop()
	}()

	for packet := range p.packetChan {
		err := p.conn.WriteJSON(packet)
		if err != nil {
			log.Printf("WebsocketInterfacer %q: %v", p.id, err)

			break
		}
	}
}

func (p *WebsocketInterfacer) Stop() {
	if err := p.conn.Close(); err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseMessage) {
			log.Printf("%q", err)
		}
	}

	p.player = nil

	log.Printf("WebsocketInterfacer %q: Stopped", p.id)
}

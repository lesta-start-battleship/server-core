package players

import (
	"lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors"
	"lesta-battleship/server-core/pkg/matchmaking/packets"
	"log"
)

type Player struct {
	id       string
	strategy Strategy

	conn actors.ClientInterfacer

	messageChan chan packets.Packet
}

func NewPlayer(id string, conn actors.ClientInterfacer) *Player {
	player := &Player{
		id:       id,
		strategy: nil,

		conn: conn,

		messageChan: make(chan packets.Packet, 256),
	}

	conn.ConnectTo(player)

	return player
}

func (p *Player) Id() string {
	return p.id
}

func (p *Player) ChangeStrategy(newStrategy Strategy) {
	if p.strategy != nil {
		p.strategy.OnExit()
	}

	p.strategy = newStrategy

	log.Printf("Player %q: Changed strategy to %q", p.id, newStrategy)
}

func (p *Player) GetPacket(senderId string, packet packets.Packet) {
	p.messageChan <- packet

	log.Printf("Player %q: Received packet %T from %q", p.id, packet.Body, senderId)
}

func (p *Player) Start() {
	defer func() {
		if _, ok := <-p.messageChan; !ok {
			close(p.messageChan)
		}
		p.Stop()
	}()

	log.Printf("Player %q: Started", p.id)

	for packet := range p.messageChan {
		p.handlePacket(packet.SenderId, packet)
	}
}

func (p *Player) Stop() {
	if p.strategy != nil {
		p.strategy.OnExit()
	}

	if p.conn != nil {
		p.conn.Stop()
	}

	log.Printf("Player %q: Stopped", p.id)
}

func (p *Player) handlePacket(senderId string, packet packets.Packet) {
	if p.id == senderId {
		p.strategy.HandlePacket(senderId, packet)
	} else {
		p.conn.GetPacket(senderId, packet)
	}
}

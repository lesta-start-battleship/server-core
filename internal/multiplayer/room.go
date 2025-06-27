package multiplayer

import (
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

const MaxPlayers = 2

type Room struct {
	id         string
	players    [MaxPlayers]*Player
	matchmaker *Matchmaker

	connectChan    chan *Player
	disconnectChan chan *Player
	broadcastChan  chan packets.Packet
}

func NewRoom(id string, matchmaker *Matchmaker) *Room {
	return &Room{
		id:         id,
		players:    [MaxPlayers]*Player{nil, nil},
		matchmaker: matchmaker,

		connectChan:    make(chan *Player),
		disconnectChan: make(chan *Player),
		broadcastChan:  make(chan packets.Packet),
	}
}

func (r *Room) Run() {
	defer r.Close()

	for {
		select {
		case client := <-r.connectChan:
			err := r.Connect(client)
			if err != nil {
				log.Printf("Room: %v", err)
			}
		case client := <-r.disconnectChan:
			err := r.Disconnect(client)
			if err != nil {
				log.Printf("Room: %v", err)
			}
		case packet := <-r.broadcastChan:
			r.Broadcast(packet)
		}
	}
}

func (r *Room) Full() bool {
	return r.players[0] != nil && r.players[1] != nil
}

func (r *Room) GetMessage(packet packets.Packet) {
	log.Printf("Room %q: Received message from %q", r.id, packet.SenderId)

	r.broadcastChan <- packet
}

func (r *Room) Connect(player *Player) error {
	log.Printf("Room %q: Received player %q", r.id, player.id)
	for i, position := range r.players {
		if position == nil {
			if err := player.ConnectToRoom(r); err != nil {
				return err
			}
			r.players[i] = player

			log.Printf("Room %q: Connected player %q", r.id, player.id)

			return nil
		}
	}

	return ErrRoomIsFull
}

func (r *Room) Disconnect(player *Player) error {
	for i, position := range r.players {
		if position != nil && position.id == player.id {
			if err := player.DisconnectFromRoom(); err != nil {
				return err
			}
			r.players[i] = nil

			log.Printf("Room %q: Disconnected player %q", r.id, player.id)

			return nil
		}
	}

	return ErrNotConnectedToRoom
}

func (r *Room) Broadcast(packet packets.Packet) {
	log.Printf("Room %q: Broadcasting to %v packet by %q", r.id, r.players, packet.SenderId)
	for _, player := range r.players {
		log.Printf("Room %q: Iterating through player %v", r.id, player)
		if player != nil && player.id != packet.SenderId {
			player.GetMessage(packet)
		}
	}
}

// TODO: Probably unsafe
func (r *Room) Close() {
	for _, player := range r.players {
		if player != nil {
			player.Close()
		}
	}

	r.matchmaker.deleteRoomChan <- r
	r.matchmaker = nil

	if _, ok := <-r.connectChan; !ok {
		close(r.connectChan)
	}
	if _, ok := <-r.disconnectChan; !ok {
		close(r.disconnectChan)
	}
	if _, ok := <-r.broadcastChan; !ok {
		close(r.broadcastChan)
	}

	log.Printf("Room %q: Closed", r.id)
}

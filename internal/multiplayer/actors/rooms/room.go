package rooms

import (
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/internal/multiplayer/actors/players"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

const MaxPlayers = 2

type Room struct {
	id      string
	players [MaxPlayers]*players.Player

	playerRegistry players.PlayerRegistry

	matchmaker actors.Actor

	packetChan chan packets.Packet
}

func NewRoom(id string, playerRegistry players.PlayerRegistry, matchmaker actors.Actor) *Room {
	return &Room{
		id:      id,
		players: [MaxPlayers]*players.Player{nil, nil},

		playerRegistry: playerRegistry,

		matchmaker: matchmaker,

		packetChan: make(chan packets.Packet, 256),
	}
}

func (r *Room) Id() string {
	return r.id
}

func (r *Room) GetPacket(senderId string, packet packets.Packet) {
	r.packetChan <- packet

	log.Printf("Room %q: Received message from %q", r.id, senderId)
}

func (r *Room) Start() {
	defer r.Stop()

	for packet := range r.packetChan {
		r.handlePacket(packet.SenderId, packet)
	}
}

// TODO: Probably unsafe
func (r *Room) Stop() {
	for _, player := range r.players {
		if player != nil {
			player.Stop()
		}
	}

	if _, ok := <-r.packetChan; !ok {
		close(r.packetChan)
	}

	log.Printf("Room %q: Closed", r.id)
}

func (r *Room) handlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case packets.ConnectPlayer:
		r.handleConnect(senderId, packet)
	case packets.Disconnect:
		r.handleDisconnect(senderId, packet)
	case packets.PlayerMessage:
		r.handleBroadcast(senderId, packet)
	default:
		log.Printf("Room %q: Received incorrect packet %t from %q", r.id, packet, senderId)
	}
}

func (r *Room) Full() bool {
	return r.players[0] != nil && r.players[1] != nil
}

func (r *Room) handleConnect(senderId string, packet packets.ConnectPlayer) error {
	player := r.playerRegistry.Find(senderId)

	for i, position := range r.players {
		if position == nil {
			players.SetInRoom(player, r)
			r.players[i] = player

			log.Printf("Room %q: Connected player %q", r.id, player.Id())
			log.Printf("Room %q: %v", r.id, r.players)

			return nil
		}
	}

	return ErrRoomIsFull
}

func (r *Room) handleDisconnect(senderId string, packet packets.Disconnect) error {
	player := r.playerRegistry.Find(senderId)

	for i, position := range r.players {
		if position != nil && position.Id() == player.Id() {
			players.SetInSearch(player, r.matchmaker)
			r.players[i] = nil

			log.Printf("Room %q: Disconnected player %q from room", r.id, player.Id())
			player.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})

			return nil
		}
	}

	return ErrNotConnectedToRoom
}

func (r *Room) handleBroadcast(senderId string, packet packets.PlayerMessage) {
	for _, player := range r.players {
		log.Printf("Room %q: Iterating through player %v", r.id, player)
		if player != nil && player.Id() != senderId {
			player.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})
		}
	}

	log.Printf("Room %q: Broadcasted to players packet by %q", r.id, senderId)
}

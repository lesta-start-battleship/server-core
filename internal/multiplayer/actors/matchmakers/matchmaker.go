package matchmakers

import (
	"lesta-battleship/server-core/internal/infra"
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/internal/multiplayer/actors/hubs"
	"lesta-battleship/server-core/internal/multiplayer/actors/players"
	"lesta-battleship/server-core/internal/multiplayer/actors/rooms"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type MatchType int

const (
	RandomMatch MatchType = iota
	RankedMatch
	GuildMatch
	CustomMatch
)

type Matchmaker struct {
	playerRegistry players.PlayerRegistry
	roomRegistry   rooms.RoomRegistry

	hub actors.Actor

	packetChan chan packets.Packet
}

func NewMatchmaker(playerRegistry players.PlayerRegistry, roomRegistry rooms.RoomRegistry) *Matchmaker {
	return &Matchmaker{
		playerRegistry: playerRegistry,
		roomRegistry:   roomRegistry,

		hub: nil,

		packetChan: make(chan packets.Packet, 256),
	}
}

func (mm *Matchmaker) SetParent(hub *hubs.Hub) {
	mm.hub = hub
}

func (mm *Matchmaker) GetPacket(senderId string, packet packets.Packet) {
	mm.packetChan <- packet

	log.Printf("Matchmaker: Received packet from %q", packet.SenderId)
}

func (mm *Matchmaker) Start() {
	defer mm.Stop()

	for packet := range mm.packetChan {
		mm.handlePacket(packet.SenderId, packet)
	}
}

// TODO: Probably unsafe
func (mm *Matchmaker) Stop() {
	if _, ok := <-mm.packetChan; !ok {
		close(mm.packetChan)
	}

	log.Println("Matchmaker: Closed")
}

func (mm *Matchmaker) handlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case packets.JoinSearch:
		mm.handleJoinRoom(senderId, packet)
	case packets.Disconnect:
		mm.handleDisconnect(senderId, packet)
	}
}

func (mm *Matchmaker) handleJoinRoom(senderId string, packet packets.JoinSearch) {
	// TODO: Add PROPER matchmaking
	rooms := mm.roomRegistry.Rooms()
	for id, room := range rooms {
		if !room.Full() {
			room.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packets.ConnectPlayer{Id: senderId}})

			log.Printf("Matchmaker: Connected to Random room %q", id)

			return
		}
	}

	room := mm.createRandomRoom(infra.GenerateId())
	mm.roomRegistry.Track(room.Id(), room)

	room.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packets.ConnectPlayer{Id: senderId}})
}

func (mm *Matchmaker) handleDisconnect(senderId string, packet packets.Disconnect) {
	mm.hub.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})
}

func (mm *Matchmaker) createRandomRoom(id string) *rooms.Room {
	room := createRoom(id, mm.playerRegistry, mm)
	mm.roomRegistry.Track(id, room)

	log.Printf("Matchmaker: Created random room %q", room.Id())

	return room
}

func (mm *Matchmaker) deleteRoom(room *rooms.Room) {
	mm.roomRegistry.Delete(room.Id())
	room.Stop()

	log.Printf("Matchmaker: Deregistered room %q", room.Id())
}

func createRoom(id string, playerRegistry players.PlayerRegistry, mm *Matchmaker) *rooms.Room {
	room := rooms.NewRoom(id, playerRegistry, mm)
	go room.Start()

	return room
}

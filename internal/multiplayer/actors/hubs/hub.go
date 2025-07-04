package hubs

import (
	"lesta-battleship/server-core/internal/multiplayer/actors/matchmakers"
	"lesta-battleship/server-core/internal/multiplayer/actors/players"
	"lesta-battleship/server-core/internal/multiplayer/actors/rooms"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type Hub struct {
	matchmakerRegistry matchmakers.MatchmakerRegistry
	roomRegistry       rooms.RoomRegistry
	playerRegistry     players.PlayerRegistry

	packetChan chan packets.Packet
}

func NewHub(matchmakerRegistry matchmakers.MatchmakerRegistry, roomRegistry rooms.RoomRegistry, playerRegistry players.PlayerRegistry) *Hub {
	hub := &Hub{
		matchmakerRegistry: matchmakerRegistry,
		roomRegistry:       roomRegistry,
		playerRegistry:     playerRegistry,

		packetChan: make(chan packets.Packet),
	}

	return hub
}

func (h *Hub) Id() string {
	return "Hub"
}

func (h *Hub) GetPacket(senderId string, packet packets.Packet) {
	h.packetChan <- packet

	log.Printf("Hub: Received packet %T from %q", packet.Body, senderId)
}

func (h *Hub) SendPacket(receiverId string, packet packets.Packet) {
}

func (h *Hub) Start() {
	defer func() {
		if _, ok := <-h.packetChan; !ok {
			close(h.packetChan)
		}
		h.Stop()
	}()

	log.Println("Hub: Started")

	for packet := range h.packetChan {
		h.handlePacket(packet.SenderId, packet)
	}
}

// TODO: Probably unsafe
func (h *Hub) Stop() {
	h.matchmakerRegistry = nil
	h.roomRegistry = nil
	h.playerRegistry = nil

	log.Println("Hub: Closed")
}

func (h *Hub) handlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case *packets.JoinSearch:
		h.handleSearch(senderId, packet)
	case *packets.Disconnect:
		h.handleDisconnect(senderId, packet)
	default:
		log.Printf("Hub: Received incorrect packet %T from %q", packet, senderId)
	}
}

func (h *Hub) handleSearch(senderId string, packet *packets.JoinSearch) {
	matchType := packet.MatchType
	matchmaker := h.matchmakerRegistry.Find(matchType)
	if matchmaker == nil {
		log.Printf("Hub: Received incorrect matchType %q", matchType)

		return
	}
	matchmaker.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})

	log.Printf("Hub: Send %q to %q Pool", senderId, matchType)
}

func (h *Hub) handleDisconnect(senderId string, packet *packets.Disconnect) {
	player := h.playerRegistry.Find(senderId)
	player.Stop()
	h.playerRegistry.Delete(player.Id())

	log.Printf("Hub: Disconnected player %q", player.Id())
	log.Printf("Hub: %+v", h.playerRegistry.Players())
}

package hubs

import (
	"lesta-battleship/server-core/internal/multiplayer/actors/matchmakers"
	"lesta-battleship/server-core/internal/multiplayer/actors/players"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type Hub struct {
	playerRegistry     players.PlayerRegistry
	matchmakerRegistry matchmakers.MatchmakerRegistry

	packetChan chan packets.Packet
}

func NewHub(playerRegistry players.PlayerRegistry, matchmakerRegistry matchmakers.MatchmakerRegistry) *Hub {
	hub := &Hub{
		playerRegistry:     playerRegistry,
		matchmakerRegistry: matchmakerRegistry,

		packetChan: make(chan packets.Packet),
	}

	return hub
}

func (h *Hub) Id() string {
	return "Hub"
}

func (h *Hub) GetPacket(senderId string, packet packets.Packet) {
	h.packetChan <- packet

	log.Printf("Hub: Received packet from %q", senderId)
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

	for packet := range h.packetChan {
		h.handlePacket(packet.SenderId, packet)
	}
}

// TODO: Probably unsafe
func (h *Hub) Stop() {
	h.playerRegistry = nil
	h.matchmakerRegistry = nil

	log.Println("Hub: Closed")
}

func (h *Hub) handlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case packets.JoinSearch:
		h.handleConnect(senderId, packet)
	case packets.Disconnect:
		h.handleDisconnect(senderId, packet)
	default:
		log.Printf("Hub: Received incorrect packet %t from %q", packet, senderId)
	}
}

func (h *Hub) handleConnect(senderId string, packet packets.JoinSearch) {
	matchType := packet.MatchType
	matchmaker := h.matchmakerRegistry.Find(matchType)
	if matchmaker == nil {
		log.Printf("Hub: Received incorrect MatchType %q", matchType)

		return
	}
	matchmaker.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})

	log.Printf("Hub: Send %q to %q Pool", senderId, matchType)
}

func (h *Hub) handleDisconnect(senderId string, packet packets.Disconnect) {
	player := h.playerRegistry.Find(senderId)
	player.Stop()
	h.playerRegistry.Delete(player.Id())

	log.Printf("Hub: Disconnected player %q", player.Id())
	log.Printf("Hub: %+v", h.playerRegistry.Players())
}

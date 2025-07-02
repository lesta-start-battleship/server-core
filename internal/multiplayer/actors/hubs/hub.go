package hubs

import (
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/internal/multiplayer/actors/players"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type Hub struct {
	playerRegistry players.PlayerRegistry

	matchmaker actors.Actor

	packetChan chan packets.Packet
}

func NewHub(mm actors.Actor, playerRegistry players.PlayerRegistry) *Hub {
	hub := &Hub{
		playerRegistry: playerRegistry,

		matchmaker: mm,

		packetChan: make(chan packets.Packet),
	}

	return hub
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
	h.matchmaker.Stop()
	h.matchmaker = nil

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
	switch {
	default:
		h.matchmaker.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})

		log.Printf("Hub: Send  %q to Random Pool", senderId)
	}
}

func (h *Hub) handleDisconnect(senderId string, packet packets.Disconnect) {
	log.Printf("Hub: Disconnecting %q", senderId)

	player := h.playerRegistry.Find(senderId)
	player.Stop()
	h.playerRegistry.Delete(player.Id())

	log.Printf("Hub: Disconnected player %q", player.Id())
	log.Printf("Hub: %+v", h.playerRegistry.Players())
}

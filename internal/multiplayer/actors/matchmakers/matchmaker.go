package matchmakers

import (
	"lesta-battleship/server-core/internal/infra"
	"lesta-battleship/server-core/internal/multiplayer/actors"
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
	id    string
	queue map[string]*players.Player

	playerRegistry players.PlayerRegistry
	roomRegistry   rooms.RoomRegistry

	hub      actors.Actor
	strategy Strategy

	packetChan chan packets.Packet
}

func NewMatchmaker(id string, playerRegistry players.PlayerRegistry, roomRegistry rooms.RoomRegistry) *Matchmaker {
	return &Matchmaker{
		id:    id,
		queue: make(map[string]*players.Player),

		playerRegistry: playerRegistry,
		roomRegistry:   roomRegistry,

		hub: nil,

		packetChan: make(chan packets.Packet, 256),
	}
}

func (mm *Matchmaker) Id() string {
	return mm.id
}

func (mm *Matchmaker) ConnectTo(hub actors.Actor) {
	mm.hub = hub
}

func (mm *Matchmaker) ChangeStrategy(newStrategy Strategy) {
	if mm.strategy != nil {
		mm.strategy.OnExit()
	}

	mm.strategy = newStrategy

	log.Printf("Matchmaker %q: Changed strategy to %q", mm.id, newStrategy)
}

func (mm *Matchmaker) GetPacket(senderId string, packet packets.Packet) {
	mm.packetChan <- packet

	log.Printf("Matchmaker %q: Received packet from %q", mm.id, packet.SenderId)
}

func (mm *Matchmaker) Start() {
	defer func() {
		if _, ok := <-mm.packetChan; !ok {
			close(mm.packetChan)
		}
		mm.Stop()
	}()

	log.Printf("Matchmaker %q: Started", mm.id)

	for packet := range mm.packetChan {
		mm.handlePacket(packet.SenderId, packet)
	}
}

func (mm *Matchmaker) Stop() {
	if mm.strategy != nil {
		mm.strategy.OnExit()
	}
	mm.strategy = nil

	mm.hub = nil
	mm.playerRegistry = nil
	mm.roomRegistry = nil

	log.Println("Matchmaker: Stopped")
}

func (mm *Matchmaker) handlePacket(senderId string, packet packets.Packet) {
	mm.strategy.HandlePacket(senderId, packet)
}

func (mm *Matchmaker) CreateRoom() actors.Actor {
	id := infra.GenerateId()
	room := createRoom(id, mm)
	mm.roomRegistry.Track(id, room)

	log.Printf("Matchmaker: Created room %q", room.Id())

	return room
}

func (mm *Matchmaker) DeleteRoom(room *rooms.Room) {
	mm.roomRegistry.Delete(room.Id())
	room.Stop()

	log.Printf("Matchmaker: Deregistered room %q", room.Id())
}

func (mm *Matchmaker) AddToQueue(playerId string) {
	player := mm.playerRegistry.Find(playerId)
	mm.queue[playerId] = player
	players.SetInSearch(player, mm)

	log.Printf("Matchmaker %q: Added to queue player %q", mm.id, playerId)
	log.Printf("Matchmaker %q: Queue %v", mm.id, mm.queue)
}

func (mm *Matchmaker) RemoveFromQueue(playerId string) {
	delete(mm.queue, playerId)

	log.Printf("Matchmaker %q: Removed from queue player %q", mm.id, playerId)
	log.Printf("Matchmaker %q: Queue %v", mm.id, mm.queue)
}

func createRoom(id string, mm *Matchmaker) *rooms.Room {
	room := rooms.NewRoom(id, mm.playerRegistry, mm)
	go room.Start()

	return room
}

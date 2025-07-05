package matchmakers

import (
	"lesta-battleship/server-core/internal/infra"
	"lesta-battleship/server-core/internal/app/multiplayer/actors"
	"lesta-battleship/server-core/internal/app/multiplayer/actors/players"
	"lesta-battleship/server-core/internal/app/multiplayer/actors/rooms"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type Matchmaker struct {
	id       string
	strategy Strategy
	queue    map[string]*players.Player

	playerRegistry players.PlayerRegistry
	roomRegistry   rooms.RoomRegistry

	hub actors.Actor

	packetChan chan packets.Packet
}

func NewMatchmaker(
	id string,
	playerRegistry players.PlayerRegistry,
	roomRegistry rooms.RoomRegistry,
	hub actors.Actor,
) *Matchmaker {
	return &Matchmaker{
		id:       id,
		strategy: nil,
		queue:    make(map[string]*players.Player),

		playerRegistry: playerRegistry,
		roomRegistry:   roomRegistry,

		hub: hub,

		packetChan: make(chan packets.Packet, 256),
	}
}

func (mm *Matchmaker) Id() string {
	return mm.id
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

	log.Printf("Matchmaker %q: Received packet %T from %q", mm.id, packet.Body, packet.SenderId)
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

	for id, player := range mm.queue {
		mm.playerRegistry.Delete(id)
		player.Stop()
	}
	mm.queue = nil

	mm.playerRegistry = nil
	mm.roomRegistry = nil

	mm.hub = nil

	log.Printf("Matchmaker %q: Stopped", mm.id)
}

func (mm *Matchmaker) handlePacket(senderId string, packet packets.Packet) {
	mm.strategy.HandlePacket(senderId, packet)
}

func (mm *Matchmaker) CreateRoom() actors.Actor {
	id := infra.GenerateId()
	room := createRoom(id, mm)
	mm.roomRegistry.Track(id, room)

	log.Printf("Matchmaker %q: Created room %q", mm.id, room.Id())

	return room
}

func (mm *Matchmaker) ConnectToRoom(roomId, playerId string) {
	room := mm.roomRegistry.Find(roomId)
	if room == nil {
		log.Printf("Matchmaker %q: Room %q is nil", mm.id, roomId)

		return
	}

	room.GetPacket(playerId, packets.NewConnectPlayer(playerId, playerId))

	log.Printf("Matchmaker %q: Send player %q to room %q", mm.id, playerId, roomId)
}

func (mm *Matchmaker) DeleteRoom(room *rooms.Room) {
	mm.roomRegistry.Delete(room.Id())
	room.Stop()

	log.Printf("Matchmaker: Deleted room %q", room.Id())
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

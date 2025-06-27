package multiplayer

import "log"

type Matchmaker struct {
	rooms map[*Room]struct{}

	randomPoolChan chan *Player
	deleteRoomChan chan *Room
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		rooms:          make(map[*Room]struct{}),
		randomPoolChan: make(chan *Player),
		deleteRoomChan: make(chan *Room),
	}
}

func (mm *Matchmaker) Run() {
	for {
		select {
		case player := <-mm.randomPoolChan:
			mm.JoinRandomRoom(player)
		case room := <-mm.deleteRoomChan:
			mm.deleteRoom(room)
		}
	}
}

func (mm *Matchmaker) JoinRandomRoom(player *Player) {
	// TODO: Add PROPER matchmaking
	for room := range mm.rooms {
		if !room.Full() {
			room.connectChan <- player

			return
		}
	}

	room := mm.createRoom(generateId())

	room.connectChan <- player
}

func (mm *Matchmaker) createRoom(id string) *Room {
	room := NewRoom(id, mm)
	mm.registerRoom(room)

	go room.Run()
	log.Printf("Matchmaker: Created room %q", room.id)

	return room
}

func (mm *Matchmaker) registerRoom(room *Room) {
	mm.rooms[room] = struct{}{}

	log.Printf("Matchmaker: Registered room %q", room.id)
}

func (mm *Matchmaker) deleteRoom(room *Room) {
	delete(mm.rooms, room)

	log.Printf("Matchmaker: Deregistered room %q", room.id)
}

// TODO: Probably unsafe
func (mm *Matchmaker) Close() {
	for room := range mm.rooms {
		delete(mm.rooms, room)
		room.Close()
	}

	if _, ok := <-mm.randomPoolChan; !ok {
		close(mm.randomPoolChan)
	}
	if _, ok := <-mm.deleteRoomChan; !ok {
		close(mm.deleteRoomChan)
	}

	log.Println("Matchmaker: Closed")
}

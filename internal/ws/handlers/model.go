package handlers

import "lesta-battleship/server-core/internal/game"

type EventInput struct {
	Event string    `json:"event"`
	Ship  game.Ship `json:"ship"`
	X     int       `json:"x"`
	Y     int       `json:"y"`
}

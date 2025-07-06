package handlers

import "lesta-battleship/server-core/internal/game-core/game"

const (
	// success
	EventShipPlaced     = "ship_placed"
	EventShipRemoved    = "ship_removed"
	EventReadyConfirmed = "ready_confirmed"
	EventShootResult    = "shoot_result"
	EventGameStart      = "game_start"
	EventGameEnd        = "game_end"
	// error
	EventPlaceShipError  = "place_ship_error"
	EventRemoveShipError = "remove_ship_error"
	EventReadyError      = "ready_error"
	EventShootError      = "shoot_error"
)

type EventInput struct {
	Event  string    `json:"event"`
	Ship   game.Ship `json:"ship"`
	X      int       `json:"x"`
	Y      int       `json:"y"`
	ItemID int       `json:"itemid"`
}

type ShipPlacedResponse struct {
	Coords []game.Coord `json:"coords"`
}

type ReadyConfirmedResponse struct {
	AllReady bool `json:"all_ready"`
}

type GameStartResponse struct {
	FirstTurn string `json:"first_turn"`
}

type ShipRemovedResponse struct {
	Coords []game.Coord `json:"coords"`
}

type ShootResultResponse struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	By       string `json:"by"`
	Hit      bool   `json:"hit"`
	NextTurn string `json:"next_turn"`
	GameOver bool   `json:"game_over"`
}

type GameEndResponse struct {
	Winner string `json:"winner"`
}

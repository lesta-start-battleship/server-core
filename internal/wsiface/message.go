package wsiface

import (
	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/items"
)

const (
	// success для ручек
	EventShipPlaced     = "ship_placed"     // place_ship
	EventShipRemoved    = "ship_removed"    // remove_ship
	EventReadyConfirmed = "ready_confirmed" // ready
	EventShootResult    = "shoot_result"    // shoot
	EventItemUsed       = "item_used"       // use_item

	// success для начала и конца игры
	EventGameStart = "game_start"
	EventGameEnd   = "game_end"

	// error
	EventError = "event_error"
)

// TODO: нужно подумать про инпут предметов
type WSInput struct {
	Event     string         `json:"event"`
	Ship      *game.Ship     `json:"ship,omitempty"`
	X         int            `json:"x"`
	Y         int            `json:"y"`
	ItemID    int            `json:"item_id,omitempty"`
	Direction int            `json:"direction,omitempty"`
}

type WSResponse struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type ShipPlacedResponse struct {
	Coords []game.Coord `json:"coords"`
}

type ShipRemovedResponse struct {
	Coords []game.Coord `json:"coords"`
}

type ReadyConfirmedResponse struct {
	AllReady bool `json:"all_ready"`
}

type ShootResultResponse struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	By       string `json:"by"`
	Hit      bool   `json:"hit"`
	NextTurn string `json:"next_turn"`
	GameOver bool   `json:"game_over"`
}

type GameStartResponse struct {
	FirstTurn string `json:"first_turn"`
}

type GameEndResponse struct {
	Winner string `json:"winner"`
}

// TODO: нужно подумать про респонс айтемов
type ItemUsedResponse struct {
	ItemID  items.ItemID       `json:"item_id"`
	Name    string             `json:"name"`
	By      string             `json:"by"`
	Effects []items.ItemEffect `json:"effects"`
}

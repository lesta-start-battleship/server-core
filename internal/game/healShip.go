package game

import (
	"errors"
)

type HealShipCommand struct {
	Coords Coord // - получаем от юзера
}

func NewHealShipCommand(target Coord) *HealShipCommand {
	return &HealShipCommand{Coords: target}
}

func (c *HealShipCommand) Apply(states *States) error {
	gs := states.PlayerState
	// проверяем валидность координаты
	if !gs.isInside(c.Coords) {
		return errors.New("out of bounds")
	}

	// достаем id корабля и проверяем надо ли чинить данную палубу
	shipID := gs.Field[c.Coords.X][c.Coords.Y].ShipID
	if shipID == Empty {
		return errors.New("cell is empty")
	}
	ship := gs.Ships[shipID]
	if stateCell := ship.Decks[c.Coords]; stateCell == Whole {
		return errors.New("cell is whole")
	}

	if ship.Health == 0 {
		return errors.New("ship is dead")
	}
	ship.Decks[c.Coords] = Whole
	ship.Health += 1

	return nil
}

func (c *HealShipCommand) Undo(states *States) {
	gs := states.PlayerState
	shipID := gs.Field[c.Coords.X][c.Coords.Y].ShipID
	ship := gs.Ships[shipID]
	ship.Decks[c.Coords] = Hit
	ship.Health -= 1
}

func (c *HealShipCommand) GetHealedCoord() Coord {
	return c.Coords
}

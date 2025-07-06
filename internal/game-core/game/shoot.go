package game

import (
	"errors"
)

type ShootCommand struct {
	Target    Coord // получаем при создании
	PrevState CellState
	Success   bool
}

func NewShootCommand(target Coord) *ShootCommand {
	return &ShootCommand{Target: target}
}

func (c *ShootCommand) Apply(states *States) error {
	gs := states.EnemyState
	if !gs.isInside(c.Target) {
		return errors.New("out of bounds")
	}
	c.PrevState = gs.Field[c.Target.X][c.Target.Y]

	if c.PrevState.ShipID != Empty {
		ship := gs.Ships[c.PrevState.ShipID]
		if ship.Decks[c.Target] == Whole {
			ship.Decks[c.Target] = Hit
			ship.Health -= 1
			if ship.Health == 0 {
				gs.NumShips -= 1
			}
			c.Success = true
		}
	}
	gs.Field[c.Target.X][c.Target.Y].State = Open
	return nil
}

func (c *ShootCommand) Undo(states *States) {
	gs := states.EnemyState
	if c.Success {
		ship := gs.Ships[c.PrevState.ShipID]
		ship.Decks[c.Target] = Whole
		if ship.Health == 0 {
			gs.NumShips += 1
		}
		ship.Health += 1
	}
	gs.Field[c.Target.X][c.Target.Y].State = c.PrevState.State
}

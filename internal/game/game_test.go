package game

import (
	//"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShootCommand(t *testing.T) {
	t.Run("shoot out of bounds", func(t *testing.T) {
		states := &States{
			EnemyState: NewGameState(),
		}
		cmd := NewShootCommand(Coord{X: -1, Y: 5})
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "out of bounds", err.Error())
	})

	t.Run("shoot empty cell", func(t *testing.T) {
		states := &States{
			EnemyState: NewGameState(),
		}
		target := Coord{X: 5, Y: 5}
		cmd := NewShootCommand(target)
		err := cmd.Apply(states)
		assert.NoError(t, err)
		assert.Equal(t, Open, states.EnemyState.Field[target.X][target.Y].State)
		assert.False(t, cmd.Success)

		// Test undo
		cmd.Undo(states)
		assert.Equal(t, Empty, states.EnemyState.Field[target.X][target.Y].State)
	})

}

func TestPlaceShipCommand(t *testing.T) {
	t.Run("place ship out of bounds", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewPlaceShipCommand(3, Coord{X: -1, Y: 5}, Vertical)
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "starting coordinate is out of bounds", err.Error())
	})

	t.Run("place ship that goes out of bounds", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewPlaceShipCommand(3, Coord{X: 8, Y: 8}, Horizontal)
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "ship placement goes out of bounds", err.Error())
	})

	t.Run("place ship adjacent to another", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		// Place first ship
		cmd1 := NewPlaceShipCommand(3, Coord{X: 2, Y: 2}, Horizontal)
		assert.NoError(t, cmd1.Apply(states))

		// Try to place adjacent ship
		cmd2 := NewPlaceShipCommand(2, Coord{X: 1, Y: 1}, Vertical)
		err := cmd2.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "ship overlaps or is adjacent to another", err.Error())
	})

	t.Run("place valid ship", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewPlaceShipCommand(3, Coord{X: 2, Y: 2}, Horizontal)
		err := cmd.Apply(states)
		assert.NoError(t, err)

		// Check ship placement
		ship := cmd.ship
		assert.NotNil(t, ship)
		assert.Equal(t, 3, ship.Len)
		assert.Equal(t, 3, ship.Health)

		// Check field
		for i := 0; i < 3; i++ {
			x := 2 + i
			y := 2
			assert.Equal(t, ship.ID, states.PlayerState.Field[x][y].ShipID)
		}

		// Test undo
		cmd.Undo(states)
		for i := 0; i < 3; i++ {
			x := 2 + i
			y := 2
			assert.Equal(t, Empty, states.PlayerState.Field[x][y].ShipID)
		}
		assert.Nil(t, states.PlayerState.Ships[ship.ID])
	})

}

func TestRemoveShipCommand(t *testing.T) {
	t.Run("remove from empty cell", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewRemoveShipCommand(Coord{X: 5, Y: 5})
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "empty cell", err.Error())
	})

	t.Run("remove non-existent ship", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		// Set ship ID in cell but don't create the ship
		states.PlayerState.Field[3][3].ShipID = 1
		cmd := NewRemoveShipCommand(Coord{X: 3, Y: 3})
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "ship not found", err.Error())
	})

	t.Run("remove valid ship", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		// Place a ship first
		placeCmd := NewPlaceShipCommand(3, Coord{X: 2, Y: 2}, Horizontal)
		assert.NoError(t, placeCmd.Apply(states))
		initialShips := states.PlayerState.NumShips

		// Remove it
		cmd := NewRemoveShipCommand(Coord{X: 2, Y: 2})
		err := cmd.Apply(states)
		assert.NoError(t, err)
		assert.Equal(t, initialShips-1, states.PlayerState.NumShips)

		// Check field is cleared
		for i := 0; i < 3; i++ {
			x := 2 + i
			y := 2
			assert.Equal(t, Empty, states.PlayerState.Field[x][y].ShipID)
		}

		// Test undo
		cmd.Undo(states)
		assert.Equal(t, initialShips, states.PlayerState.NumShips)
		for i := 0; i < 3; i++ {
			x := 2 + i
			y := 2
			assert.Equal(t, placeCmd.ship.ID, states.PlayerState.Field[x][y].ShipID)
		}
	})
}

func TestHealShipCommand(t *testing.T) {
	t.Run("heal out of bounds", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewHealShipCommand(Coord{X: -1, Y: 5})
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "out of bounds", err.Error())
	})

	t.Run("heal empty cell", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewHealShipCommand(Coord{X: 5, Y: 5})
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "cell is empty", err.Error())
	})

	t.Run("heal whole cell", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		// Place a ship
		placeCmd := NewPlaceShipCommand(3, Coord{X: 2, Y: 2}, Horizontal)
		assert.NoError(t, placeCmd.Apply(states))

		cmd := NewHealShipCommand(Coord{X: 2, Y: 2})
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "cell is whole", err.Error())
	})

}

func TestOpenCellCommand(t *testing.T) {
	t.Run("open cell out of bounds", func(t *testing.T) {
		states := &States{
			EnemyState: NewGameState(),
		}
		cmd := NewOpenCellCommand(Coord{X: -1, Y: 5})
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "out of bounds", err.Error())
	})

	t.Run("open empty cell", func(t *testing.T) {
		states := &States{
			EnemyState: NewGameState(),
		}
		target := Coord{X: 5, Y: 5}
		cmd := NewOpenCellCommand(target)
		err := cmd.Apply(states)
		assert.NoError(t, err)
		assert.Equal(t, Open, states.EnemyState.Field[target.X][target.Y].State)
		assert.False(t, cmd.ShipFound)

		// Test undo
		cmd.Undo(states)
		assert.Equal(t, Empty, states.EnemyState.Field[target.X][target.Y].State)
	})


}

func TestPlaceSubmarineCommand(t *testing.T) {
	t.Run("place submarine out of bounds", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewPlaceSubmarineCommand(Coord{X: -1, Y: 5}, Vertical)
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "starting coordinate is out of bounds", err.Error())
	})

	t.Run("place submarine that goes out of bounds", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewPlaceSubmarineCommand(Coord{X: 8, Y: 8}, Horizontal)
		err := cmd.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "ship placement goes out of bounds", err.Error())
	})

	t.Run("place submarine adjacent to another ship", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		// Place first ship
		cmd1 := NewPlaceShipCommand(3, Coord{X: 2, Y: 2}, Horizontal)
		assert.NoError(t, cmd1.Apply(states))

		// Try to place adjacent submarine
		cmd2 := NewPlaceSubmarineCommand(Coord{X: 1, Y: 1}, Vertical)
		err := cmd2.Apply(states)
		assert.Error(t, err)
		assert.Equal(t, "ship overlaps or is adjacent to another", err.Error())
	})

	t.Run("place valid submarine", func(t *testing.T) {
		states := &States{
			PlayerState: NewGameState(),
		}
		cmd := NewPlaceSubmarineCommand(Coord{X: 2, Y: 2}, Horizontal)
		err := cmd.Apply(states)
		assert.NoError(t, err)

		// Check submarine placement
		submarine := cmd.ship
		assert.NotNil(t, submarine)
		assert.Equal(t, 11, submarine.ID)
		assert.Equal(t, 3, submarine.Len)
		assert.Equal(t, 3, submarine.Health)

		// Check field
		for i := 0; i < 3; i++ {
			x := 2 + i
			y := 2
			assert.Equal(t, submarine.ID, states.PlayerState.Field[x][y].ShipID)
		}

		// Test undo
		cmd.Undo(states)
		for i := 0; i < 3; i++ {
			x := 2 + i
			y := 2
			assert.Equal(t, Empty, states.PlayerState.Field[x][y].ShipID)
		}
		assert.Nil(t, states.PlayerState.Ships[submarine.ID])
	})
}
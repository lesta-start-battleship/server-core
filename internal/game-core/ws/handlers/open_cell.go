package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game-core/game"
	"lesta-battleship/server-core/internal/game-core/match"
	"lesta-battleship/server-core/internal/game-core/transaction"

	"github.com/gorilla/websocket"
)

func HandleOpenCell(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input WSInput) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if room.Status != "playing" {
		err := errors.New("game not started")
		SendError(conn, err.Error())
		return err
	}

	if room.Turn != player.ID {
		err := errors.New("not your turn")
		SendError(conn, err.Error())
		return err
	}

	if input.X < 0 || input.X >= 10 || input.Y < 0 || input.Y >= 10 {
		err := errors.New("coordinates out of bounds")
		SendError(conn, err.Error())
		return err
	}

	cell := player.States.EnemyState.Field[input.X][input.Y]
	if cell.State == game.Open {
		err := errors.New("this cell is already opened")
		SendError(conn, err.Error())
		return err
	}

	cmd := game.NewOpenCellCommand(game.Coord{X: input.X, Y: input.Y})
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(player.States); err != nil {
		SendError(conn, err.Error())
		return err
	}

	Broadcast(room, "cell_opened", map[string]any{
		"x":          input.X,
		"y":          input.Y,
		"by":         player.ID,
		"ship_found": cmd.ShipFound,
	})

	return nil
}

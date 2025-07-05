package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game"
	"lesta-battleship/server-core/internal/match"
	"lesta-battleship/server-core/internal/transaction"

	"github.com/gorilla/websocket"
)

func HandleRemoveShip(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input EventInput) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if room.Status != "waiting" {
		err := errors.New("cannot remove ship during game")
		Send(conn, "remove_ship_error", err.Error())
		return err
	}

	if player.Ready {
		err := errors.New("you cannot remove ship after ready")
		Send(conn, "remove_ship_error", err.Error())
		return err
	}

	cmd := game.NewRemoveShipCommand(game.Coord{X: input.X, Y: input.Y})
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(player.States); err != nil {
		Send(conn, "remove_ship_error", err.Error())
		return err
	}

	Send(conn, "ship_removed", map[string]any{
		"coords": cmd.GetDeckCoords(),
	})

	return nil
}

package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game"
	"lesta-battleship/server-core/internal/match"
	"lesta-battleship/server-core/internal/transaction"

	"github.com/gorilla/websocket"
)

func HandleHealShip(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input EventInput) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if room.Status != "playing" {
		err := errors.New("game not started")
		Send(conn, "heal_ship_error", err.Error())
		return err
	}

	if room.Turn != player.ID {
		err := errors.New("not your turn")
		Send(conn, "heal_ship_error", err.Error())
		return err
	}

	cmd := game.NewHealShipCommand(game.Coord{X: input.X, Y: input.Y})
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(player.States); err != nil {
		Send(conn, "heal_ship_error", err.Error())
		return err
	}

	Send(conn, "ship_healed", map[string]any{
		"coords": cmd.GetHealedCoord(),
	})
	return nil
}

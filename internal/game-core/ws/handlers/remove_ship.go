package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game-core/game"
	"lesta-battleship/server-core/internal/game-core/match"
	"lesta-battleship/server-core/internal/game-core/transaction"

	"github.com/gorilla/websocket"
)

func HandleRemoveShip(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input WSInput) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if room.Status != "waiting" {
		err := errors.New("cannot remove ship during game")
		SendError(conn, err.Error())
		return err
	}

	if player.Ready {
		err := errors.New("you cannot remove ship after ready")
		SendError(conn, err.Error())
		return err
	}

	cmd := game.NewRemoveShipCommand(game.Coord{X: input.X, Y: input.Y})
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(player.States); err != nil {
		SendError(conn, err.Error())
		return err
	}

	SendSuccess(conn, EventShipRemoved, ShipRemovedResponse{
		Coords: cmd.GetDeckCoords(),
	})

	return nil
}

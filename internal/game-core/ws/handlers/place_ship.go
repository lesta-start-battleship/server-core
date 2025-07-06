package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game-core/game"
	"lesta-battleship/server-core/internal/game-core/match"
	"lesta-battleship/server-core/internal/game-core/transaction"

	"github.com/gorilla/websocket"
)

func HandlePlaceShip(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input WSInput) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if room.Status != "waiting" {
		err := errors.New("game already started")
		Send(conn, EventPlaceShipError, err.Error())
		return err
	}

	if player.Ready {
		err := errors.New("cannot place ship after ready")
		Send(conn, EventPlaceShipError, err.Error())
		return err
	}

	if player.States.PlayerState.NumShips >= 10 {
		err := errors.New("maximum 10 ships allowed")
		Send(conn, EventPlaceShipError, err.Error())
		return err
	}

	cmd := game.NewPlaceShipCommand(input.Ship.Len, input.Ship.Coords, input.Ship.Bearings)
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(player.States); err != nil {
		Send(conn, EventPlaceShipError, err.Error())
		return err
	}

	Send(conn, EventShipPlaced, ShipPlacedResponse{
		Coords: cmd.GetDeckCoords(),
	})

	return nil
}

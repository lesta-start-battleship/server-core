package handlers

import (
	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/transaction"
	"github.com/lesta-battleship/server-core/internal/wsiface"
)

type PlaceShipHandler struct{}

func (h *PlaceShipHandler) EventName() string {
	return "place_ship"
}

func (h *PlaceShipHandler) Handle(input any, ctx *wsiface.Context) error {
	ctx.Room.Mutex.Lock()
	defer ctx.Room.Mutex.Unlock()

	wsInput, ok := input.(wsiface.WSInput)
	if !ok {
		return SendError(ctx.Conn, "invalid input format for place_ship")
	}

	if ctx.Room.Status != "waiting" {
		return SendError(ctx.Conn, "game already started")
	}
	if ctx.Player.Ready {
		return SendError(ctx.Conn, "cannot place ship after ready")
	}
	if ctx.Player.States.PlayerState.NumShips >= 10 {
		return SendError(ctx.Conn, "maximum 10 ships allowed")
	}

	cmd := game.NewPlaceShipCommand(wsInput.Ship.Len, wsInput.Ship.Coords, wsInput.Ship.Bearings)
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(ctx.Player.States); err != nil {
		return SendError(ctx.Conn, err.Error())
	}

	return SendSuccess(ctx.Conn, wsiface.EventShipPlaced, wsiface.ShipPlacedResponse{
		Coords: cmd.GetDeckCoords(),
	})
}

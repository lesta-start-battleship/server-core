package handlers

import (
	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/transaction"
	"github.com/lesta-battleship/server-core/internal/wsiface"
)

type RemoveShipHandler struct{}

func (h *RemoveShipHandler) EventName() string {
	return "remove_ship"
}

func (h *RemoveShipHandler) Handle(input any, ctx *wsiface.Context) error {
	ctx.Room.Mutex.Lock()
	defer ctx.Room.Mutex.Unlock()

	wsInput, ok := input.(wsiface.WSInput)
	if !ok {
		return SendError(ctx.Conn, "invalid input format for remove_ship")
	}

	if ctx.Room.Status != "waiting" {
		return SendError(ctx.Conn, "cannot remove ship during game")
	}

	if ctx.Player.Ready {
		return SendError(ctx.Conn, "you cannot remove ship after ready")
	}

	cmd := game.NewRemoveShipCommand(game.Coord{X: wsInput.X, Y: wsInput.Y})
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(ctx.Player.States); err != nil {
		return SendError(ctx.Conn, err.Error())
	}

	return SendSuccess(ctx.Conn, wsiface.EventShipRemoved, wsiface.ShipRemovedResponse{
		Coords: cmd.GetDeckCoords(),
	})
}

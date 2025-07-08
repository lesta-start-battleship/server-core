package handlers

import (
	"fmt"

	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/transaction"
	"github.com/lesta-battleship/server-core/internal/wsiface"
)

type MoveSubmarineHandler struct{}

func (h *MoveSubmarineHandler) EventName() string {
	return "move_submarine"
}

func (h *MoveSubmarineHandler) Handle(input any, ctx *wsiface.Context) error {
	ctx.Room.Mutex.Lock()
	defer ctx.Room.Mutex.Unlock()

	wsInput, ok := input.(wsiface.WSInput)
	if !ok {
		return SendError(ctx.Conn, "invalid input format for move_submarine")
	}

	if ctx.Room.Status != "playing" {
		return SendError(ctx.Conn, "game not started")
	}

	if ctx.Room.Turn != ctx.Player.ID {
		return SendError(ctx.Conn, "not your turn")
	}

	// Проверка кулдауна подлодки
	// (нужно сохранить информацию о последнем ходе подлодки — например, в `PlayerConn`)
	const SubmarineCooldown = 3
	if ctx.Player.LastSubmarineTurn > 0 &&
		ctx.Player.MoveCount-ctx.Player.LastSubmarineTurn < SubmarineCooldown {
		wait := SubmarineCooldown - (ctx.Player.MoveCount - ctx.Player.LastSubmarineTurn)
		return SendError(ctx.Conn, "submarine is on cooldown, wait "+fmt.Sprint(wait)+" turns")
	}

	// Проверка что по координате (x, y) стоит подлодка
	from := game.Coord{X: wsInput.X, Y: wsInput.Y}
	ship := ctx.Player.States.PlayerState.FindShipByCoord(from)
	if ship == nil || ship.ID != 11 {
		return SendError(ctx.Conn, "no submarine at given coordinate")
	}

	// Перемещение
	to := game.Coord{X: wsInput.X2, Y: wsInput.Y2}
	bearings := wsInput.Direction != 0

	tx := transaction.NewTransaction()
	tx.Add(game.NewRemoveShipCommand(from))

	place := game.NewPlaceSubmarineCommand(to, bearings)
	tx.Add(place)

	if err := tx.Execute(ctx.Player.States); err != nil {
		return SendError(ctx.Conn, err.Error())
	}

	ctx.Player.LastSubmarineTurn = ctx.Player.MoveCount
	ctx.Player.MoveCount++

	if ctx.Player == ctx.Room.Player1 {
		ctx.Room.Turn = ctx.Room.Player2.ID
	} else {
		ctx.Room.Turn = ctx.Room.Player1.ID
	}

	return SendSuccess(ctx.Conn, "submarine_moved", wsiface.ShipPlacedResponse{
		Coords: place.GetDeckCoords(),
	})
}

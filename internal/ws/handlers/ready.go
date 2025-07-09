package handlers

import (
	"github.com/lesta-battleship/server-core/internal/wsiface"
)

type ReadyHandler struct{}

func (h *ReadyHandler) EventName() string {
	return "ready"
}

func (h *ReadyHandler) Handle(input any, ctx *wsiface.Context) error {
	ctx.Room.Mutex.Lock()
	defer ctx.Room.Mutex.Unlock()
	
	if ctx.Player.States.PlayerState.NumShips < 10 {
		return SendError(ctx.Conn, "you must place 10 ships before ready")
	}

	ctx.Player.Ready = true
	allReady := ctx.Room.Player1.Ready && ctx.Room.Player2.Ready
	shouldStart := false

	if allReady && ctx.Room.Status == "waiting" {
		ctx.Room.Status = "playing"
		ctx.Room.Turn = ctx.Room.Player1.ID
		shouldStart = true
	}

	if err := SendSuccess(ctx.Conn, wsiface.EventReadyConfirmed, wsiface.ReadyConfirmedResponse{
		AllReady: allReady,
	}); err != nil {
		return err
	}

	if shouldStart {
		if err := Broadcast(ctx.Room, wsiface.EventGameStart, wsiface.GameStartResponse{
			FirstTurn: ctx.Room.Turn,
		}); err != nil {
			return err
		}
	}

	return nil
}

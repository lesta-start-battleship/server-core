package handlers

import (
	"log"

	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/items"
	"github.com/lesta-battleship/server-core/internal/wsiface"
)

type UseItemHandler struct{}

func (h *UseItemHandler) EventName() string {
	return "use_item"
}

func (h *UseItemHandler) Handle(input any, ctx *wsiface.Context) error {
	ctx.Room.Mutex.Lock()
	defer ctx.Room.Mutex.Unlock()

	wsInput, ok := input.(wsiface.WSInput)
	if !ok {
		return SendError(ctx.Conn, "invalid input format for use_item")
	}

	if ctx.Room.Status != "playing" {
		return SendError(ctx.Conn, "game not started")
	}

	if ctx.Room.Turn != ctx.Player.ID {
		return SendError(ctx.Conn, "not your turn")
	}

	itemID := items.ItemID(wsInput.ItemID)
	if ctx.Player.Items[itemID] <= 0 {
		return SendError(ctx.Conn, "item not available or already used")
	}
	log.Println("item dostupen")

	itemData, ok := ctx.Room.Items[itemID]
	if !ok {
		return SendError(ctx.Conn, "item metadata not found")
	}
	log.Println("nashli iz obshey kollekcii")

	result, err := items.UseItem(itemID, ctx.Player.States, ctx.Room.Items, wsInput.Params)
	if err != nil {
		log.Printf("[WS] Use item error: %v", err)
		return SendError(ctx.Conn, err.Error())
	}
	log.Println("usenuli item")

	ctx.Player.Items[itemID]--

	if err := ctx.Dispatcher.DispatchUsedItem(event.Item{
		PlayerID: ctx.Player.ID,
		ItemID:   wsInput.ItemID,
	}); err != nil {
		log.Printf("[KAFKA] Failed to dispatch used item: %v", err)
	}

	return Broadcast(ctx.Room, wsiface.EventItemUsed, wsiface.ItemUsedResponse{
		ItemID: itemID,
		Name:   itemData.Name,
		By:     ctx.Player.ID,
		Result: result,
	})
}

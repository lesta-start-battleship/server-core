package handlers

import (
	"fmt"
	"log"

	"github.com/lesta-battleship/server-core/internal/config"
	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/items"
	"github.com/lesta-battleship/server-core/internal/wsiface"

	"bytes"
	"encoding/json"
	"net/http"
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
		return SendError(ctx.Conn, "item not available")
	}

	itemData, ok := ctx.Room.Items[itemID]
	if !ok {
		return SendError(ctx.Conn, "item metadata not found")
	}

	// чекаем Cooldown и UseLimit
	usage := ctx.Player.ItemUsage[itemID]
	if usage == nil {
		usage = &items.ItemUsageData{}
		ctx.Player.ItemUsage[itemID] = usage
	}

	if itemData.UseLimit > 0 && usage.UsedTimes >= itemData.UseLimit {
		return SendError(ctx.Conn, "item use limit reached")
	}

	if itemData.Cooldown > 0 && (ctx.Player.MoveCount-usage.LastUsedTurn) < itemData.Cooldown {
		waitTurns := itemData.Cooldown - (ctx.Player.MoveCount - usage.LastUsedTurn)
		return SendError(ctx.Conn, "item on cooldown, wait more turns: "+fmt.Sprint(waitTurns))
	}

	if itemData.Name == "Конь" || itemData.Name == "Ладья" || itemData.Name == "Ферзь" || itemData.Name == "Слон" {
		ctx.Player.ChessFigureCount++
	}

	if ctx.Player.ChessFigureCount > 2 {
		return SendError(ctx.Conn, "chess figure use limit reached")
	}

	itemInput := items.ItemInput{
		X:         wsInput.X,
		Y:         wsInput.Y,
		X2:        wsInput.X2,
		Y2:        wsInput.Y2,
		X3:        wsInput.X3,
		Y3:        wsInput.Y3,
		Direction: wsInput.Direction,
		ItemID:    wsInput.ItemID,
	}

	effect, err := items.RunScript(itemData.Script, ctx.Player.States, itemInput)
	if err != nil {
		return SendError(ctx.Conn, err.Error())
	}

	ctx.Player.Items[itemID]--
	usage.UsedTimes++
	usage.LastUsedTurn = ctx.Player.MoveCount

	ctx.Player.MoveCount++

	if err := ctx.Dispatcher.DispatchUsedItem(event.Item{
		PlayerID: ctx.Player.ID,
		ItemID:   wsInput.ItemID,
	}); err != nil {
		log.Printf("[KAFKA] Failed to dispatch used item: %v", err)
	}
	go ForTestInventoryTheyDidNotMakeConsumerWeWillDeleteThisFuncLater(wsInput, ctx)

	return Broadcast(ctx.Room, wsiface.EventItemUsed, wsiface.ItemUsedResponse{
		ItemID:  itemID,
		Name:    itemData.Name,
		By:      ctx.Player.ID,
		Effects: effect,
	})
}

func ForTestInventoryTheyDidNotMakeConsumerWeWillDeleteThisFuncLater(wsInput wsiface.WSInput, ctx *wsiface.Context) {
	// сообщить об использовании предмета еще бл раз, ну вот зачем им делать ручку и говорить что должен все это делать серв кор, когда есть брокер??????? (за 10 часов до сдачи)
	// в итоге получился вот такой говнокод)
	// ----------------------------------
	url := config.InventoryUseItem

	requestBody := map[string]any{
		"item_id": wsInput.ItemID,
		"amount":  1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", ctx.Player.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode <= 200 || resp.StatusCode > 300 {
		return
	}
}

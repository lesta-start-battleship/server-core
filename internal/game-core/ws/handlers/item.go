package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game-core/event"
	"lesta-battleship/server-core/internal/game-core/items"
	"lesta-battleship/server-core/internal/game-core/match"

	"github.com/gorilla/websocket"
)

func HandleItem(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input EventInput, dispatcher *event.MatchEventDispatcher) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if room.Status != "playing" {
		err := errors.New("game not started")
		Send(conn, "use_item_error", err.Error())
		return err
	}

	if room.Turn != player.ID {
		err := errors.New("not your turn")
		Send(conn, "use_item_error", err.Error())
		return err
	}

	if player.Items[items.ItemID(input.ItemID)] == 0 {
		err := errors.New("the user does not have this item")
		Send(conn, "use_item_error", err.Error())
		return err
	}
	// TODO: можем ли использовать данных предмет(ограничения на применения в 1 игре/ограничения на применения несколько ходов подряд)

	// TODO: пока нет идей как это делать в общем случае, предлагайте варианты
	params := map[string]interface{}{"x": 5, "y": 5}

	_, err := items.UseItem(items.ItemID(input.ItemID), player.States, room.Items, params)
	if err != nil {
		Send(conn, "use_item_error", err.Error())
		return err
	}

	// TODO: сообщить об использовании предмета

	Send(conn, "item_used", map[string]any{
		// "coords": cmd.GetHealedCoord(),
	}) // TODO: придумать как в общем виде выдавать результат исполнения командд
	return nil
}

package handlers

import "github.com/lesta-battleship/server-core/internal/wsiface"

var handlerRegistry = map[string]wsiface.WSEventHandler{}

func RegisterHandler(handler wsiface.WSEventHandler) {
	handlerRegistry[handler.EventName()] = handler
}

func GetHandler(event string) (wsiface.WSEventHandler, bool) {
	h, ok := handlerRegistry[event]
	return h, ok
}

func RegisterAllHandlers() {
	RegisterHandler(&PlaceShipHandler{})
	RegisterHandler(&RemoveShipHandler{})
	RegisterHandler(&ShootHandler{})
	RegisterHandler(&ReadyHandler{})
	RegisterHandler(&UseItemHandler{})
	RegisterHandler(&MoveSubmarineHandler{})
}

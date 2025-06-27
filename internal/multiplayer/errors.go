package multiplayer

import "errors"

var (
	ErrRoomIsFull             = errors.New("Room is full")
	ErrAlreadyConnectedToRoom = errors.New("Player already connected to room")
	ErrNotConnectedToRoom     = errors.New("Player not connected to room")
)

package packets

import (
	"encoding/json"
	"errors"
	"reflect"
)

type Packet struct {
	SenderId string     `json:"sender_id"`
	Type     string     `json:"event"`
	Body     PacketBody `json:"-"`
}

func (p *Packet) UnmarshalJSON(b []byte) error {
	data := struct {
		SenderId string          `json:"sender_id"`
		Type     string          `json:"event"`
		Body     json.RawMessage `json:"data"`
	}{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	var body PacketBody
	for _, bodyType := range packetBodyTypes {
		reflectType := reflect.TypeOf(bodyType).Elem()
		if data.Type == bodyType.Type() {
			tempBody := reflect.New(reflectType).Interface()
			if err := json.Unmarshal(data.Body, &tempBody); err != nil {
				return err
			}
			body = tempBody.(PacketBody)
			break
		}
	}

	if body == nil {
		return errors.New("unknown packet type: " + data.Type)
	}

	p.SenderId = data.SenderId
	p.Type = data.Type
	p.Body = body
	return nil
}

func (p Packet) MarshalJSON() ([]byte, error) {
	var data any
	if p.Body != nil {
		data = p.Body
	}
	return json.Marshal(struct {
		SenderId string `json:"sender_id"`
		Type     string `json:"event"`
		Data     any    `json:"data,omitempty"`
	}{
		SenderId: p.SenderId,
		Type:     p.Type,
		Data:     data,
	})
}

type PacketBody interface {
	Type() string
	isPacketBody()
}

func NewPlaceShip(senderId string, ship Ship) Packet {
	body := &PlaceShip{Ship: ship}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewRemoveShip(senderId string, x, y int) Packet {
	body := &RemoveShip{X: x, Y: y}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewReady(senderId string) Packet {
	body := &Ready{}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewShoot(senderId string, x, y int) Packet {
	body := &Shoot{X: x, Y: y}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewShootResult(senderId string, x, y int, by string, hit bool, nextTurn string, gameOver bool) Packet {
	body := &ShootResult{X: x, Y: y, By: by, Hit: hit, NextTurn: nextTurn, GameOver: gameOver}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewShipPlaced(senderId string, coords []Coord) Packet {
	body := &ShipPlaced{Coords: coords}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewShipRemoved(senderId string, coords []Coord) Packet {
	body := &ShipRemoved{Coords: coords}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewReadyConfirmed(senderId string, allReady bool) Packet {
	body := &ReadyConfirmed{AllReady: allReady}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewGameStart(senderId string, firstTurn string) Packet {
	body := &GameStart{FirstTurn: firstTurn}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewGameEnd(senderId string, winner string) Packet {
	body := &GameEnd{Winner: winner}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewError(senderId string, message string) Packet {
	body := &Error{Message: message}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

// x2, y2 and x3, y3 depend on item
func NewUseItem(senderId string, itemID, x, y int, x2, y2, x3, y3, direction int) Packet {
	body := &UseItem{
		ItemID:    itemID,
		X:         x,
		Y:         y,
		X2:        x2,
		Y2:        y2,
		X3:        x3,
		Y3:        y3,
		Direction: direction,
	}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

func NewItemUsed(senderId string, itemID int, name string, by string, effects []ItemEffect) Packet {
	body := &ItemUsed{
		ItemID:  itemID,
		Name:    name,
		By:      by,
		Effects: effects,
	}
	return Packet{SenderId: senderId, Type: body.Type(), Body: body}
}

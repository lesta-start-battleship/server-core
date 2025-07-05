package packets

import (
	"encoding/json"
	"reflect"
)

type Packet struct {
	SenderId string
	Type     string
	Body     PacketBody
}

func (p *Packet) UnmarshalJSON(b []byte) error {
	data := struct {
		SenderId string
		Type     string
		Body     json.RawMessage
	}{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	var body PacketBody
	for _, bodyType := range packetBodyTypes {
		reflectType := reflect.TypeOf(bodyType).Elem()
		typeName := reflectType.Name()
		if data.Type == typeName {
			tempBody := reflect.New(reflectType).Interface()
			if err := json.Unmarshal(data.Body, &tempBody); err != nil {
				return err
			}
			body = tempBody.(PacketBody)

			break
		}
	}

	p.SenderId = data.SenderId
	p.Type = data.Type
	p.Body = body

	return nil
}

type PacketBody interface {
	Type() string
	isPacketBody()
}

func NewPlayerMessage(senderId string, message string) Packet {
	body := &PlayerMessage{Msg: message}
	return Packet{
		SenderId: senderId,
		Type:     body.Type(),
		Body:     body,
	}
}

func NewJoinSearch(senderId string, matchType string) Packet {
	body := &JoinSearch{MatchType: matchType}
	return Packet{
		SenderId: senderId,
		Type:     body.Type(),
		Body:     body,
	}
}

func NewCreateRoom(senderId string) Packet {
	body := &CreateRoom{}
	return Packet{
		SenderId: senderId,
		Type:     body.Type(),
		Body:     body,
	}
}

func NewJoinRoom(senderId string, roomId string) Packet {
	body := &JoinRoom{Id: roomId}
	return Packet{
		SenderId: senderId,
		Type:     body.Type(),
		Body:     body,
	}
}

func NewConnectPlayer(senderId string, playerId string) Packet {
	body := &ConnectPlayer{Id: playerId}
	return Packet{
		SenderId: senderId,
		Type:     body.Type(),
		Body:     body,
	}
}

func NewDisconnect(senderId string) Packet {
	body := &Disconnect{}
	return Packet{
		SenderId: senderId,
		Type:     body.Type(),
		Body:     body,
	}
}

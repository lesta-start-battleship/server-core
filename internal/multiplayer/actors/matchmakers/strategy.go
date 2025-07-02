package matchmakers

import (
	"fmt"
	"lesta-battleship/server-core/pkg/packets"
)

type Strategy interface {
	HandlePacket(senderId string, packet packets.Packet)
	OnExit()

	fmt.Stringer
}

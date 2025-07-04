package packets

var packetBodyTypes = []PacketBody{
	(*PlayerMessage)(nil),
	(*JoinSearch)(nil),
	(*CreateRoom)(nil),
	(*JoinRoom)(nil),
	(*PlayerMessage)(nil),
	(*Disconnect)(nil),
}

type PlayerMessage struct {
	Msg string
}

func (PlayerMessage) Type() string {
	return "PlayerMessage"
}
func (PlayerMessage) isPacketBody() {}

type JoinSearch struct {
	MatchType string
}

func (JoinSearch) Type() string {
	return "JoinSearch"
}
func (JoinSearch) isPacketBody() {}

type CreateRoom struct{}

func (CreateRoom) Type() string {
	return "CreateRoom"
}
func (CreateRoom) isPacketBody() {}

type JoinRoom struct {
	Id string
}

func (JoinRoom) Type() string {
	return "JoinRoom"
}
func (JoinRoom) isPacketBody() {}

type ConnectPlayer struct {
	Id string
}

func (ConnectPlayer) Type() string {
	return "ConnectPlayer"
}
func (ConnectPlayer) isPacketBody() {}

type Disconnect struct{}

func (Disconnect) Type() string {
	return "Disconnect"
}
func (Disconnect) isPacketBody() {}

package packets

type Packet struct {
	SenderId string `json:"sender_id"`
	Body     any    `json:"body"`
}

type PlayerMessage struct {
	Msg string `json:"msg"`
}

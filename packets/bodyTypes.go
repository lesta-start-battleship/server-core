package packets

var packetBodyTypes = []PacketBody{
	(*PlaceShip)(nil),
	(*RemoveShip)(nil),
	(*Ready)(nil),
	(*Shoot)(nil),
	(*ShipPlaced)(nil),
	(*ShipRemoved)(nil),
	(*ReadyConfirmed)(nil),
	(*ShootResult)(nil),
	(*GameStart)(nil),
	(*GameEnd)(nil),
	(*Error)(nil),
	(*UseItem)(nil),
	(*ItemUsed)(nil),
}

type PlaceShip struct {
	Ship Ship `json:"ship"`
}

func (PlaceShip) Type() string  { return "place_ship" }
func (PlaceShip) isPacketBody() {}

type RemoveShip struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (RemoveShip) Type() string  { return "remove_ship" }
func (RemoveShip) isPacketBody() {}

type Ready struct{}

func (Ready) Type() string  { return "ready" }
func (Ready) isPacketBody() {}

type Shoot struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (Shoot) Type() string  { return "shoot" }
func (Shoot) isPacketBody() {}

type ShipPlaced struct {
	Coords []Coord `json:"coords"`
}

func (ShipPlaced) Type() string  { return "ship_placed" }
func (ShipPlaced) isPacketBody() {}

type ShipRemoved struct {
	Coords []Coord `json:"coords"`
}

func (ShipRemoved) Type() string  { return "ship_removed" }
func (ShipRemoved) isPacketBody() {}

type ReadyConfirmed struct {
	AllReady bool `json:"all_ready"`
}

func (ReadyConfirmed) Type() string  { return "ready_confirmed" }
func (ReadyConfirmed) isPacketBody() {}

type ShootResult struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	By       string `json:"by"`
	Hit      bool   `json:"hit"`
	NextTurn string `json:"next_turn"`
	GameOver bool   `json:"game_over"`
}

func (ShootResult) Type() string  { return "shoot_result" }
func (ShootResult) isPacketBody() {}

type GameStart struct {
	FirstTurn string `json:"first_turn"`
}

func (GameStart) Type() string  { return "game_start" }
func (GameStart) isPacketBody() {}

type GameEnd struct {
	Winner string `json:"winner"`
}

func (GameEnd) Type() string  { return "game_end" }
func (GameEnd) isPacketBody() {}

type Error struct {
	Message string `json:"message"`
}

func (Error) Type() string  { return "error" }
func (Error) isPacketBody() {}

type UseItem struct {
	ItemID    int `json:"item_id"`
	X         int `json:"x"`
	Y         int `json:"y"`
	X2        int `json:"x2,omitempty"`
	Y2        int `json:"y2,omitempty"`
	X3        int `json:"x3,omitempty"`
	Y3        int `json:"y3,omitempty"`
	Direction int `json:"direction,omitempty"`
}

func (UseItem) Type() string  { return "use_item" }
func (UseItem) isPacketBody() {}

type ItemUsed struct {
	ItemID  int          `json:"item_id"`
	Name    string       `json:"name"`
	By      string       `json:"by"`
	Effects []ItemEffect `json:"effects"`
}

func (ItemUsed) Type() string  { return "item_used" }
func (ItemUsed) isPacketBody() {}

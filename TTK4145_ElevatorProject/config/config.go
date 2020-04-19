package config

// Scaleable declaration of #floors and #elevators
const (
	NumFloors    = 4
	NumElevators = 3
	NumButtons   = 3
)

type Button int

const (
	BtnUp Button = iota
	BtnDown
	BtnInside
)

type Direction int

const (
	DirDown Direction = iota - 1
	DirStop
	DirUp
)

type Acknowledge int

const (
	Finished Acknowledge = iota - 1
	NotAcked
	Acked
)

type ElevState int

const (
	Undefined ElevState = iota - 1
	Idle
	Moving
	DoorOpen
)

type Keypress struct {
	Floor              int
	Btn                Button
	DesignatedElevator int
	Finished           bool
}

type Elev struct {
	State  ElevState
	Dir    Direction
	Floor  int
	Queue  [NumFloors][NumButtons]bool
	Online bool
}

type AckList struct {
	DesignatedElevator int
	ImplicitAcks       [NumElevators]Acknowledge
}

type Message struct {
	Elevator         [NumElevators]Elev
	RegisteredOrders [NumFloors][NumButtons - 1]AckList
	ID               int
}

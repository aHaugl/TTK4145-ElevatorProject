package elevatorstatemachine

import . "../config"

func shouldStop(elevator Elev) bool {
	switch elevator.Dir {
	case DirUp:
		return elevator.Queue[elevator.Floor][BtnUp] ||
			elevator.Queue[elevator.Floor][BtnInside] ||
			!ordersAbove(elevator)
	case DirDown:
		return elevator.Queue[elevator.Floor][BtnDown] ||
			elevator.Queue[elevator.Floor][BtnInside] ||
			!ordersBelow(elevator)
	case DirStop:
	default:
	}
	return false
}

func chooseDirection(elevator Elev) Direction {
	switch elevator.Dir {
	case DirStop:
		if ordersAbove(elevator) {
			return DirUp
		} else if ordersBelow(elevator) {
			return DirDown
		} else {
			return DirStop
		}

	case DirDown:
		if ordersBelow(elevator) {
			return DirDown
		} else if ordersAbove(elevator) {
			return DirUp
		} else {
			return DirStop
		}
	}
	return DirStop
}

func ordersAbove(elevator Elev) bool {
	for floor := elevator.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator Elev) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

package orderhandler

import (
	"math"

	. "../config"
)

func numberOfOrders(elev Elev) int {
	orders := 0
	for floor := 0; floor < NumFloors; floor++ {
		ordersAtFloor := 0
		for button := 0; button < NumButtons; button++ {
			if elev.Queue[floor][button] {
				ordersAtFloor = 1
			}
		}
		orders += ordersAtFloor
	}
	return orders
}

func costofElev(order Keypress, elev Elev) float64 {
	if !elev.Online {
		return 1000
	}
	sum := float64(order.Floor - elev.Floor)
	sum = math.Abs(sum)
	sum += float64(numberOfOrders(elev))
	if elev.Dir == DirStop && sum == 0 {
		sum += 5
	} else if elev.Dir == DirDown {
		if order.Floor > elev.Floor {
			sum += 3
		}
	} else if elev.Dir == DirUp {
		if order.Floor < elev.Floor {
			sum += 3
		}
	} else {
		if sum == 0 {
			return sum
		}
	}
	return sum
}

// func costCalculator(order Keypress, elevList [NumElevators]Elev, nrElev int, elevatorCh chan Elev) int {
// 	var (
// 		elevators [NumElevators]Elev
// 	)

// 	elevators[nrElev] = <-elevatorCh

// 	if order.Btn == BtnInside {
// 		return nrElev
// 	}
// 	minCost := (NumButtons * NumFloors) * NumElevators
// 	bestElevator := nrElev
// 	for elevator := 0; elevator < NumElevators; elevator++ {
// 		if !elevators[nrElev].Online {
// 			// Disregarding offline elevators
// 			continue
// 		}
// 		cost := order.Floor - elevList[elevator].Floor

// 		if cost == 0 && elevList[elevator].State != Moving {
// 			bestElevator = elevator
// 			return bestElevator
// 		}

// 		if cost < 0 {
// 			cost = -cost
// 			if elevList[elevator].Dir == DirUp {
// 				cost += 3
// 			}
// 		} else if cost > 0 {
// 			if elevList[elevator].Dir == DirDown {
// 				cost += 3
// 			}
// 		}
// 		if cost == 0 && elevList[elevator].State == Moving {
// 			cost += 4
// 		}

// 		if elevList[elevator].State == DoorOpen {
// 			cost++
// 		}

// 		if cost < minCost {
// 			minCost = cost
// 			bestElevator = elevator
// 		}
// 	}
// 	fmt.Println("Cost of elevator",  newOrder.Btn, "is", newOrder.Floor+1)
// 	return bestElevator
// }

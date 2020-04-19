package queue

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

	if elev.Dir == DirDown {
		if order.Floor > elev.Floor {
			sum += 3
		}
	} else if elev.Dir == DirUp {
		if order.Floor < elev.Floor {
			sum += 3
		}
	} else if elev.Dir != DirStop && sum == 0 {
		sum += 5
	} else {
		if sum == 0 {
			return sum
		}
	}
	return sum
}

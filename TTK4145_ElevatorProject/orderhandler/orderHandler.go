package orderhandler

import (
	. "../config"
	hw "../hardware"
)

func getOrder(floor int, buttonType int, elevator Elev) bool {
	return elevator.Queue[floor][buttonType]
}

func setOrder(floor int, buttonType int, elevator Elev, set bool) {
	elevator.Queue[floor][buttonType] = set
}

// en goroutine som mottar knappetrykk og states og sørger for at riktig heis går til riktig plass.
// trenger vel ikke returnere noe??

func OrderHandler(order chan Keypress, nrElev int, completedOrderCh chan int, updateLightsCh chan [NumElevators]Elev,
	newOrderCh chan Keypress, elevatorCh chan Elev, updateQueueCh chan [NumElevators]Elev, updateSyncCh chan Elev,
	orderUpdateCh chan Keypress) {

	var (
		elevators      [NumElevators]Elev
		completedOrder Keypress
	)

	elevators[nrElev] = <-elevatorCh
	updateSyncCh <- elevators[nrElev]
	for {
		select {
		case orderLocal := <-order:
			if !elevators[nrElev].Online && orderLocal.Btn == BtnInside {
				elevators[nrElev].Queue[orderLocal.Floor][BtnInside] = true
				updateLightsCh <- elevators
				go func() { newOrderCh <- orderLocal }()
			} else if !elevators[nrElev].Online && orderLocal.Btn == BtnDown {
				elevators[nrElev].Queue[orderLocal.Floor][BtnDown] = true
				updateLightsCh <- elevators
				go func() { newOrderCh <- orderLocal }()
			} else if !elevators[nrElev].Online && orderLocal.Btn == BtnUp {
				elevators[nrElev].Queue[orderLocal.Floor][BtnUp] = true
				updateLightsCh <- elevators
				go func() { newOrderCh <- orderLocal }()

			} else {
				if orderLocal.Floor == elevators[nrElev].Floor && elevators[nrElev].State != Moving {
					newOrderCh <- orderLocal
				} else {
					if !existingOrder(orderLocal, elevators, nrElev) {
						var sums [NumElevators]float64
						for elevator := 0; elevator < NumElevators; elevator++ {
							//sums[elevator] = costCalculator(orderLocal, elevators, nrElev, elevatorCh)
							sums[elevator] = costofElev(orderLocal, elevators[elevator])
							if elevator != 0 {
								if sums[elevator] < sums[orderLocal.DesignatedElevator-1] {
									orderLocal.DesignatedElevator = elevator + 1
									orderUpdateCh <- orderLocal
								}
							} else {
								orderLocal.DesignatedElevator = 1
							}
						}
					}
				}
			}
		case completedFloor := <-completedOrderCh:
			var button Button
			for btn := BtnUp; btn < NumButtons; btn++ {
				if elevators[nrElev].Queue[completedFloor][btn] {
					button = btn
				}
				for elevator := 0; elevator < NumElevators; elevator++ {
					if button != BtnInside || elevator == nrElev {
						elevators[elevator].Queue[completedFloor][button] = false
					}
				}
			}
			if elevators[nrElev].Online {
				orderUpdateCh <- completedOrder
			}
			updateLightsCh <- elevators

		case newElevator := <-elevatorCh:
			newQueue := elevators[nrElev].Queue
			if elevators[nrElev].State == Undefined && newElevator.State != Undefined {
				elevators[nrElev].Online = true
			}
			elevators[nrElev] = newElevator
			elevators[nrElev].Queue = newQueue
			if elevators[nrElev].Online {
				updateSyncCh <- elevators[nrElev]
			}
		case tempElevList := <-updateQueueCh:
			newOrder := false
			for elevator := 0; elevator < NumElevators; elevator++ {
				if nrElev == elevator {
					continue
				}
				if elevators[elevator].Queue != tempElevList[elevator].Queue {
					newOrder = true
				}
				elevators[elevator] = tempElevList[elevator]
			}
			for button := BtnUp; button < NumButtons; button++ {
				for floor := 0; floor < NumFloors; floor++ {
					if !elevators[nrElev].Queue[floor][button] && tempElevList[nrElev].Queue[floor][button] {
						elevators[nrElev].Queue[floor][button] = true
						order := Keypress{Floor: floor, Btn: button, DesignatedElevator: nrElev, Finished: false}
						go func() { newOrderCh <- order }()
						newOrder = true
					} else if !elevators[nrElev].Queue[floor][button] && tempElevList[nrElev].Queue[floor][button] {
						elevators[nrElev].Queue[floor][button] = false
						order := Keypress{Floor: floor, Btn: button, DesignatedElevator: nrElev, Finished: true}
						go func() { newOrderCh <- order }()
						newOrder = true
					}
				}
			}
			if newOrder {
				updateLightsCh <- elevators
			}
		}
	}
}

func SetLights(updateLightsCh <-chan [NumElevators]Elev, nrElev int) {
	var orders [NumElevators]bool

	for {
		elevators := <-updateLightsCh
		for floor := 0; floor < NumFloors; floor++ {
			for button := BtnUp; button < NumButtons; button++ {
				for elevator := 0; elevator < NumElevators; elevator++ {
					orders[elevator] = false
					if elevator != nrElev && (button == BtnInside || button == BtnDown || button == BtnUp) {
						continue
					}
					if elevators[elevator].Queue[floor][button] {
						hw.SetButtonLamp(button, floor, 1)
						orders[elevator] = true
					} else {
						hw.SetButtonLamp(button, floor, 0)
					}
				}
			}
		}
	}
}

func existingOrder(order Keypress, elevators [NumElevators]Elev, nrElev int) bool {
	if elevators[nrElev].Queue[order.Floor][BtnInside] && order.Btn == BtnInside {
		return true
	}
	for elev := 0; elev < NumElevators; elev++ {
		if elevators[nrElev].Queue[order.Floor][order.Btn] {
			return true
		}
	}
	return false
}

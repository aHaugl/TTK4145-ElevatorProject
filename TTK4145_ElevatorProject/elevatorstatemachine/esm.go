package elevatorstatemachine

import (
	"fmt"
	"time"

	. "../config"
	hw "../hardware"
)

// StateMachineChannels contains all channels between governor - esm and hardware - esm
type StateMachineChannels struct {
	OrderComplete  chan int
	Elevator       chan Elev
	StateError     chan error
	NewOrder       chan Keypress
	ArrivedAtFloor chan int
}

func RunElevator(ch StateMachineChannels) {
	elevator := Elev{
		State:  Idle,
		Dir:    DirStop,
		Floor:  hw.GetFloorSensorSignal(),
		Queue:  [NumFloors][NumButtons]bool{},
		Online: true,
	}
	doorTimedOut := time.NewTimer(3 * time.Second)
	engineErrorTimer := time.NewTimer(3 * time.Second)
	doorTimedOut.Stop()
	engineErrorTimer.Stop()
	orderCleared := false
	ch.Elevator <- elevator

	for {
		select {
		case newOrder := <-ch.NewOrder:
			if newOrder.Finished {
				elevator.Queue[newOrder.Floor][BtnUp] = false
				elevator.Queue[newOrder.Floor][BtnDown] = false
				orderCleared = true
			} else {
				elevator.Queue[newOrder.Floor][newOrder.Btn] = true
			}

			switch elevator.State {
			case Idle:
				elevator.Dir = chooseDirection(elevator)
				hw.SetMotorDirection(elevator.Dir)
				if elevator.Dir == DirStop {
					elevator.State = DoorOpen
					hw.SetDoorOpenLamp(1)
					doorTimedOut.Reset(3 * time.Second)
					go func() { ch.OrderComplete <- newOrder.Floor }()
					elevator.Queue[elevator.Floor] = [NumButtons]bool{}
				} else {
					elevator.State = Moving
					engineErrorTimer.Reset(3 * time.Second)
				}

			case Moving:
			case DoorOpen:
				if elevator.Floor == newOrder.Floor {
					doorTimedOut.Reset(3 * time.Second)
					go func() { ch.OrderComplete <- newOrder.Floor }()
					elevator.Queue[elevator.Floor] = [NumButtons]bool{}
				}

			case Undefined:
			default:
				fmt.Println("Fatal error: Reboot system")
			}
			ch.Elevator <- elevator

		case elevator.Floor = <-ch.ArrivedAtFloor:
			fmt.Println("Arrived at floor", elevator.Floor+1)
			if shouldStop(elevator) ||
				(!shouldStop(elevator) && elevator.Queue == [NumFloors][NumButtons]bool{} && orderCleared) {
				orderCleared = false
				hw.SetDoorOpenLamp(1)
				engineErrorTimer.Stop()
				elevator.State = DoorOpen
				hw.SetMotorDirection(DirStop)
				doorTimedOut.Reset(3 * time.Second)
				elevator.Queue[elevator.Floor] = [NumButtons]bool{}
				go func() { ch.OrderComplete <- elevator.Floor }()

			} else if elevator.State == Moving {
				engineErrorTimer.Reset(3 * time.Second)
			}
			ch.Elevator <- elevator

		case <-doorTimedOut.C:
			hw.SetDoorOpenLamp(0)
			elevator.Dir = chooseDirection(elevator)
			if elevator.Dir == DirStop {
				elevator.State = Idle
				engineErrorTimer.Stop()
			} else {
				elevator.State = Moving
				engineErrorTimer.Reset(3 * time.Second)
				hw.SetMotorDirection(elevator.Dir)
			}
			ch.Elevator <- elevator

		case <-engineErrorTimer.C:
			hw.SetMotorDirection(DirStop)
			elevator.State = Undefined
			fmt.Println("\x1b[1;1;33m", "Engine Error - Go offline", "\x1b[0m")
			for i := 0; i < 10; i++ {
				if i%2 == 0 {
					hw.SetStopLamp(1)
				} else {
					hw.SetStopLamp(0)
				}
				time.Sleep(time.Millisecond * 200)
			}
			hw.SetMotorDirection(elevator.Dir)
			ch.Elevator <- elevator
			engineErrorTimer.Reset(5 * time.Second)
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	. "./config"
	esm "./elevatorstatemachine"
	hw "./hardware"
	"./networkCommunication/bcast"
	"./networkCommunication/peers"
	ordH "./orderhandler"
	sync "./syncElevators"
	// . "github.com/perkjelsvik/TTK4145-sanntid/project/config"
	// gov "github.com/perkjelsvik/TTK4145-sanntid/project/elevatorGovernor"
	// esm "github.com/perkjelsvik/TTK4145-sanntid/project/elevatorStateMachine"
	// hw "github.com/perkjelsvik/TTK4145-sanntid/project/hardware"
	// "github.com/perkjelsvik/TTK4145-sanntid/project/networkCommunication/network/bcast"
	// "github.com/perkjelsvik/TTK4145-sanntid/project/networkCommunication/network/peers"
	// sync "github.com/perkjelsvik/TTK4145-sanntid/project/syncElevators"
)

func main() {
	var (
		runType string
		id      string
		e       hw.Elev_type
		ID      int
		simPort string
	)

	flag.StringVar(&runType, "run", "", "run type")
	flag.StringVar(&id, "id", "0", "id of this peer")
	flag.StringVar(&simPort, "simPort", "44523", "simulation server port")
	flag.Parse()
	ID, _ = strconv.Atoi(id)

	if runType == "sim" {
		e = hw.ET_Simulation
		fmt.Println("Running in simulation mode!")
	}

	esmChans := esm.StateMachineChannels{
		OrderComplete:  make(chan int),
		Elevator:       make(chan Elev),
		NewOrder:       make(chan Keypress),
		ArrivedAtFloor: make(chan int),
	}
	syncChans := sync.SyncChannels{
		UpdateQueue:     make(chan [NumElevators]Elev),
		UpdateSync:      make(chan Elev),
		OrderUpdate:     make(chan Keypress),
		OnlineElevators: make(chan [NumElevators]bool),
		IncomingMsg:     make(chan Message),
		OutgoingMsg:     make(chan Message),
		PeerUpdate:      make(chan peers.PeerUpdate),
		PeerTxEnable:    make(chan bool),
	}
	var (
		btnsPressedChan = make(chan Keypress)
		updateLightsCh  = make(chan [NumElevators]Elev)
	)

	hw.Init(e, btnsPressedChan, esmChans.ArrivedAtFloor, simPort)

	go hw.ButtonPoller(btnsPressedChan)
	go hw.FloorIndicatorLoop(esmChans.ArrivedAtFloor)
	go esm.RunElevator(esmChans)
	go ordH.OrderHandler(btnsPressedChan, ID, esmChans.OrderComplete, updateLightsCh, esmChans.NewOrder, esmChans.Elevator,
		syncChans.UpdateQueue, syncChans.UpdateSync, syncChans.OrderUpdate)
	go ordH.SetLights(updateLightsCh, ID)
	go sync.Synchronise(syncChans, ID)
	go bcast.Transmitter(42034, syncChans.OutgoingMsg)
	go bcast.Receiver(42034, syncChans.IncomingMsg)
	go peers.Transmitter(42035, id, syncChans.PeerTxEnable)
	go peers.Receiver(42035, syncChans.PeerUpdate)
	go killSwitch()

	select {}
}

func killSwitch() {
	// killSwitch turns the motor off if the program is killed with CTRL+C.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	hw.SetMotorDirection(DirStop)
	fmt.Println("\x1b[31;1m", "User terminated program.", "\x1b[0m")
	for i := 0; i < 10; i++ {
		hw.SetMotorDirection(DirStop)
		if i%2 == 0 {
			hw.SetStopLamp(1)
		} else {
			hw.SetStopLamp(0)
		}
		time.Sleep(200 * time.Millisecond)
	}
	hw.SetMotorDirection(DirStop)
	os.Exit(1)
}

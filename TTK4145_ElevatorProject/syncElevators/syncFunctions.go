package syncElevators

import (
	."../config"
	//. "github.com/perkjelsvik/TTK4145-sanntid/project/config"
)

func copyAckList(msg Message, registeredOrders [NumFloors][NumButtons - 1]AckList, elevator, floor, id int, btn Button) [NumFloors][NumButtons - 1]AckList {
	registeredOrders[floor][btn].ImplicitAcks[id] = msg.RegisteredOrders[floor][btn].ImplicitAcks[elevator]
	registeredOrders[floor][btn].ImplicitAcks[elevator] = msg.RegisteredOrders[floor][btn].ImplicitAcks[elevator]
	registeredOrders[floor][btn].DesignatedElevator = msg.RegisteredOrders[floor][btn].DesignatedElevator
	return registeredOrders
}

func checkAllAckStatus(onlineList [NumElevators]bool, ImplicitAcks [NumElevators]Acknowledge, status Acknowledge) bool {
	for elev := 0; elev < NumElevators; elev++ {
		if !onlineList[elev] {
			continue
		}
		if ImplicitAcks[elev] != status {
			return false
		}
	}
	return true
}

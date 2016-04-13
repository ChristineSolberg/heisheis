package messageHandler

import (
	//"../cost"
	"../elevatorControl/driver"
	"../elevatorControl/elevatorStatus"
	"../elevatorControl/orderHandling"
	"../eventHandler"
	"../network"
	"./message"
	"fmt"
	"time"
)

func MessageHandler(recvChan chan message.UpdateMessage, sendNetwork chan message.UpdateMessage, newOrderToFSM chan elevatorStatus.Elevator, elevObject chan elevatorStatus.Elevator, 
	placedOrder chan message.UpdateMessage, elevatorMap chan map[string]*elevatorStatus.Elevator, setExternalLightsOn chan [2]int, setExternalLightsOff chan [4]int, abortElev chan bool, shouldAbort chan bool) {
	var msg message.UpdateMessage
	myIP := network.GetIpAddress()
	var masterMatrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int
	elevators := make(map[string]*elevatorStatus.Elevator)
	elevatorTimers := make(map[string]*time.Timer)
	for {
		msg = <-recvChan
		msgType := msg.MessageType
		switch msgType {
		case message.IAmAlive:
			var shouldAddToMap bool = true
			for ip, _ := range elevatorTimers {
				if ip == msg.ElevatorStatus.IP {
					shouldAddToMap = false
					elevatorTimers[ip].Reset(time.Second * 2)
				}
			}
			if shouldAddToMap == true {
				elevators[msg.ElevatorStatus.IP] = new(elevatorStatus.Elevator)
				elevators[msg.ElevatorStatus.IP].Direction = msg.ElevatorStatus.Direction
				elevators[msg.ElevatorStatus.IP].CurrentFloor = msg.ElevatorStatus.CurrentFloor
				elevators[msg.ElevatorStatus.IP].PreviousFloor = msg.ElevatorStatus.PreviousFloor
				elevators[msg.ElevatorStatus.IP].State = msg.ElevatorStatus.State
				elevators[msg.ElevatorStatus.IP].IP = msg.ElevatorStatus.IP
				//updateElevator(elevators, msg.ElevatorStatus.IP ,msg)
				ip := msg.ElevatorStatus.IP
				elevatorTimers[msg.ElevatorStatus.IP] = time.AfterFunc(time.Second*2, func() {
					eventHandler.DeleteElevator(elevators, ip, sendNetwork, elevatorTimers, myIP, elevObject, abortElev, masterMatrix, shouldAbort)
				})
				eventHandler.SelectMaster(elevators)
			}
		case message.PlacedOrder:
			if elevators[myIP].Master == network.GetIpAddress(){
				floor := msg.Order[0]
				button := msg.Order[1]
				if masterMatrix[floor][button] == 0 {
					if button < 2 {
						masterMatrix[floor][button] = 1
					}
				}
				placedOrder <- msg
				elevatorMap <- elevators
			}			
		case message.AssignedOrder:
			button := msg.Order[1]
			if button < 2 {
				setExternalLightsOn <- msg.Order
			}
			if msg.ReceiverIP == network.GetIpAddress() {
				*elevators[msg.ReceiverIP] = orderHandling.AddOrderToQueue(*elevators[msg.ReceiverIP], msg.Order)
				newOrderToFSM <- *elevators[msg.ReceiverIP]
			}
		case message.CompletedOrder:
			completedOrder := msg.DelOrder
			floor := msg.DelOrder[3]
			for button := 0; button < driver.NUM_BUTTONS; button++ {
				if completedOrder[button] == 1 {
					masterMatrix[floor][button] = 0
					setExternalLightsOff <- completedOrder
				}
			}
		case message.StateUpdate:
			for _, elev := range elevators {
				fmt.Println("Elevators in map in StateUpdate: ", elev)
			}
			if elevators[msg.ElevatorStatus.IP] != nil {
				updateElevator(elevators, msg.ElevatorStatus.IP ,msg)
			}
		case message.LightUpdate:
			completedOrder := msg.DelOrder
			floor := msg.DelOrder[3]
			for button := 0; button < driver.NUM_BUTTONS-1; button++ {
				if completedOrder[button] == 1 {
					driver.Set_button_lamp(button, floor, 0)
				}
			}
			f := msg.Order[0]
			b := msg.Order[1]
			if b < 2 {
				driver.Set_button_lamp(b, f, 1)
			}
		}
	}
}

func updateElevator(elevators map[string]*elevatorStatus.Elevator, IP string, msg message.UpdateMessage){
	elevators[IP].Direction = msg.ElevatorStatus.Direction
	elevators[IP].CurrentFloor = msg.ElevatorStatus.CurrentFloor
	elevators[IP].PreviousFloor = msg.ElevatorStatus.PreviousFloor
	elevators[IP].State = msg.ElevatorStatus.State
	elevators[IP].IP = msg.ElevatorStatus.IP
	elevators[IP].OrderMatrix = msg.ElevatorStatus.OrderMatrix
	fmt.Println("OrderMatrix i updateElevator: ", elevators[IP].OrderMatrix)
}
package eventHandler

import (
	"fmt"
	"time"
	//"../network"
	"../elevatorControl/elevatorStatus"
	"../elevatorControl/orderHandling"
	"../messageHandler/message"
	"../cost"
	"../elevatorControl/driver"
)

func SelectMaster(elevators map[string]*elevatorStatus.Elevator) {
	var minimumIP string = "129.241.187.999"
	for ip, _ := range elevators {
		fmt.Println("IP: ", ip)
		if ip < minimumIP {
			minimumIP = ip
		}
	}
	for _, elev := range elevators {
		elev.Master = minimumIP
	}
	fmt.Println("New master: ", minimumIP)
}

func DeleteElevator(elevators map[string]*elevatorStatus.Elevator, IP string, sendNetwork chan message.UpdateMessage, elevatorTimers map[string]*time.Timer, myIP string, elevObject chan elevatorStatus.Elevator, abortElev chan bool, masterMatrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int, shouldAbort chan bool) {
	elev := elevators[IP]
	fmt.Println("Delete this elevator: ", elev)
	delete(elevatorTimers, IP)
	delete(elevators, IP)
	SelectMaster(elevators)
	//If the deleted elevator had any external orders they will be reassigned here
	if IP != myIP {
		for floor := 0; floor < driver.NUM_FLOORS; floor++ {
			for button := 0; button < driver.NUM_BUTTONS-1; button++ {
				if elev.OrderMatrix[floor][button] == 1 {
					var order [2]int
					order[0] = floor
					order[1] = button
					sendNetwork <- message.UpdateMessage{MessageType: message.PlacedOrder, Order: order, ReceiverIP: elevators[myIP].Master}
				}
			}
		}
	} else {
		e := <-elevObject
		for floor := 0; floor < driver.NUM_FLOORS; floor++ {
			for button := 0; button < driver.NUM_BUTTONS-1; button++ {
				e.OrderMatrix[floor][button] = masterMatrix[floor][button]
			}
		}
		elevObject <- e
		abort := <-shouldAbort
		if abort == true {
			abortElev <- true
		}
	}
}

func EventHandler(newStateUpdate chan bool, buttonPushed chan [2]int, powerOffDetected chan bool, deleteCompletedOrder chan [4]int, elevObject chan elevatorStatus.Elevator, placedOrder chan message.UpdateMessage, 
	elevatorMap chan map[string]*elevatorStatus.Elevator, setExternalLightsOn chan [2]int, setExternalLightsOff chan [4]int, sendNetwork chan message.UpdateMessage, notAlive chan bool, shouldAbort chan bool) {
	for {
		select {
		case <-newStateUpdate:
			elev := elevatorStatus.MakeCopyOfElevator(elevObject)
			sendNetwork <- message.UpdateMessage{MessageType: message.StateUpdate, ElevatorStatus: elev}
		case button := <-buttonPushed:
			elev := elevatorStatus.MakeCopyOfElevator(elevObject)
			orderHandling.WriteInternalsToFile(elev.OrderMatrix)
			if button[1] == driver.BUTTON_COMMAND {
				e := <-elevObject
				e = orderHandling.AddOrderToQueue(e, button)
				fmt.Println("OrderMatrix etter AddOrderToQueue: ", e.OrderMatrix)
				elevObject <- e
			} else {
				sendNetwork <- message.UpdateMessage{MessageType: message.PlacedOrder, ReceiverIP: elev.Master, Order: button,
					ElevatorStatus: elev}
			}
		case orderMsg := <-placedOrder:
			elevators := <- elevatorMap
			//floor := orderMsg.Order[0]
			button := orderMsg.Order[1]
			if button < 2 {
				assignedElev := cost.AssignOrders(orderMsg, elevators)
				sendNetwork <- message.UpdateMessage{MessageType: message.AssignedOrder, Order: orderMsg.Order, ReceiverIP: assignedElev}
			} else if button == driver.BUTTON_COMMAND {
				sendNetwork <- message.UpdateMessage{MessageType: message.AssignedOrder, ReceiverIP: orderMsg.ElevatorStatus.IP, Order: orderMsg.Order}
			}	
		case lightsOn := <- setExternalLightsOn:
			sendNetwork <- message.UpdateMessage{MessageType: message.LightUpdate, Order: lightsOn}
		case lightsOff := <- setExternalLightsOff:
			var noOrder [2]int
			noOrder[0] = 0
			noOrder[1] = 1
			sendNetwork <- message.UpdateMessage{MessageType: message.LightUpdate, DelOrder: lightsOff, Order: noOrder} 
		case completedOrder := <-deleteCompletedOrder:
			elev := elevatorStatus.MakeCopyOfElevator(elevObject)
			orderHandling.WriteInternalsToFile(elev.OrderMatrix)
			sendNetwork <- message.UpdateMessage{MessageType: message.CompletedOrder, DelOrder: completedOrder, ElevatorStatus: elev}
		case <-powerOffDetected:
			notAlive <- true
			shouldAbort <- true
		}
	}
}
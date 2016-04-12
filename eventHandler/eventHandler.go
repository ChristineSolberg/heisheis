package eventHandler

import (
	"fmt"
	"time"
	//"../network"
	"../elevatorControl/elevatorStatus"
	"../elevatorControl/orderHandling"
	"../messageHandler/message"
	//"../cost"
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

func DeleteElevator(elevators map[string]*elevatorStatus.Elevator, IP string, sendChan chan message.UpdateMessage, elevatorTimers map[string]*time.Timer, myIP string, elevObject chan elevatorStatus.Elevator, abortElev chan bool, masterMatrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int, shouldAbort chan bool) {
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
					sendChan <- message.UpdateMessage{MessageType: message.PlacedOrder, Order: order, RecieverIP: elevators[myIP].Master}
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

func EventHandler(newStateUpdate chan bool, buttonPushed chan [2]int, powerOffDetected chan bool, deleteCompletedOrder chan [4]int, elevObject chan elevatorStatus.Elevator, sendNetwork chan message.UpdateMessage, notAlive chan bool, shouldAbort chan bool) {
	for {
		select {
		case <-newStateUpdate:
			elev := elevatorStatus.MakeCopyOfElevator(elevObject)
			sendNetwork <- message.UpdateMessage{MessageType: message.StateUpdate, ElevatorStatus: elev}
		case order := <-buttonPushed:
			elev := elevatorStatus.MakeCopyOfElevator(elevObject)
			orderHandling.WriteInternalsToFile(elev.OrderMatrix)
			if order[1] == driver.BUTTON_COMMAND {
				e := <-elevObject
				e = orderHandling.AddOrderToQueue(e, order)
				elevObject <- e
			} else {
				sendNetwork <- message.UpdateMessage{MessageType: message.PlacedOrder, RecieverIP: elev.Master, Order: order,
					ElevatorStatus: elev}
			}
		case <-powerOffDetected:
			notAlive <- true
			shouldAbort <- true
		case completedOrder := <-deleteCompletedOrder:
			elev := elevatorStatus.MakeCopyOfElevator(elevObject)
			orderHandling.WriteInternalsToFile(elev.OrderMatrix)
			sendNetwork <- message.UpdateMessage{MessageType: message.CompletedOrder, DelOrder: completedOrder, ElevatorStatus: elev}
		}
	}
}

// func MessageHandler(recvChan chan message.UpdateMessage, sendChan chan message.UpdateMessage, newOrderToFSM chan elevatorStatus.Elevator, elevObject chan elevatorStatus.Elevator, abortElev chan bool, abortChan chan bool){ // +order message.UpdateMessage
// 	var msg message.UpdateMessage
// 	myIP := network.GetIpAddress()
// 	var MasterMatrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int
// 	elevators := make(map[string]*elevatorStatus.Elevator)
// 	elevatorTimers := make(map[string]*time.Timer)
// 	for{
// 		msg = <-recvChan
// 		msgType := msg.MessageType
// 		switch(msgType){
// 		case message.IAmAlive:
// 			var shouldAppend bool = true
// 			for ip,_ := range elevatorTimers{
// 				if ip == msg.ElevatorStatus.IP{
// 					shouldAppend = false
// 					elevatorTimers[ip].Reset(time.Second*2)
// 				}
// 			}
// 			if shouldAppend == true{
// 				elevators[msg.ElevatorStatus.IP] = new(elevatorStatus.Elevator)
// 				elevators[msg.ElevatorStatus.IP].Dir = msg.ElevatorStatus.Dir
// 				elevators[msg.ElevatorStatus.IP].CurrentFloor = msg.ElevatorStatus.CurrentFloor
// 				elevators[msg.ElevatorStatus.IP].PreviousFloor = msg.ElevatorStatus.PreviousFloor
// 				elevators[msg.ElevatorStatus.IP].State = msg.ElevatorStatus.State
// 				elevators[msg.ElevatorStatus.IP].IP = msg.ElevatorStatus.IP

// 				ip := msg.ElevatorStatus.IP
// 				elevatorTimers[msg.ElevatorStatus.IP] = time.AfterFunc(time.Second*2, func() { deleteElevator(elevators, ip, sendChan, elevatorTimers, myIP, elevObject, abortElev, MasterMatrix, abortChan)})
// 				selectMaster(elevators)
// 			}
// 		case message.PlacedOrder:
// 			floor := msg.Order[0]
// 			button := msg.Order[1]
// 			if MasterMatrix[floor][button] == 0{
// 				if elevators[myIP].Master == network.GetIpAddress(){
// 					if button < 2{
// 						MasterMatrix[floor][button] = 1
// 						AssignedElev := cost.AssignOrdersToElevator(msg, elevators)
// 						sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, RecieverIP: AssignedElev}
// 					} else if button == driver.BUTTON_COMMAND{
// 						sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, RecieverIP: msg.ElevatorStatus.IP, Order: msg.Order}
// 					}
// 				}
// 			}
// 		case message.AssignedOrder:
// 			button := msg.Order[1]
// 			if button < 2{
// 				sendChan <-message.UpdateMessage{MessageType: message.LightUpdate, Order: msg.Order}
// 			}
// 			if msg.RecieverIP == network.GetIpAddress(){
// 				*elevators[msg.RecieverIP] = orderHandling.AddOrderToQueue(*elevators[msg.RecieverIP], msg.Order)
// 				newOrderToFSM <- *elevators[msg.RecieverIP]
// 			}
// 		case message.CompletedOrder:
// 			completed := msg.DelOrder
// 			floor := msg.DelOrder[3]
// 			var noOrder [2]int
// 			noOrder[0] = 0
// 			noOrder[1] = 1
// 			for button := 0; button < driver.NUM_BUTTONS; button++{
// 				if completed[button] == 1{
// 					MasterMatrix[floor][button] = 0
// 					if button < 2{
// 						sendChan <-message.UpdateMessage{MessageType: message.LightUpdate, DelOrder: completed, Order: noOrder}
// 					} else {
// 						sendChan <-message.UpdateMessage{MessageType: message.LightUpdate, RecieverIP: msg.ElevatorStatus.IP, DelOrder: completed, Order: noOrder}
// 					}
// 				}
// 			}
// 		case message.StateUpdate:
// 			for _,elev := range elevators{
// 				fmt.Println("Elevators in map in StateUpdate: ", elev)
// 			}
// 			if elevators[msg.ElevatorStatus.IP] != nil{
// 				elevators[msg.ElevatorStatus.IP].Dir = msg.ElevatorStatus.Dir
// 				elevators[msg.ElevatorStatus.IP].CurrentFloor = msg.ElevatorStatus.CurrentFloor
// 				elevators[msg.ElevatorStatus.IP].PreviousFloor = msg.ElevatorStatus.PreviousFloor
// 				elevators[msg.ElevatorStatus.IP].State = msg.ElevatorStatus.State
// 				elevators[msg.ElevatorStatus.IP].IP = msg.ElevatorStatus.IP
// 				elevators[msg.ElevatorStatus.IP].OrderMatrix = msg.ElevatorStatus.OrderMatrix
// 			}
// 		case message.LightUpdate:
// 			completed := msg.DelOrder
// 			floor := msg.DelOrder[3]
// 			for button := 0; button < driver.NUM_BUTTONS-1; button++{
// 				if completed[button] == 1{
// 					driver.Set_button_lamp(button,floor,0)
// 				}
// 			}
// 			f := msg.Order[0]
// 			b := msg.Order[1]
// 			if b < 2{
// 			 	driver.Set_button_lamp(b, f, 1)
// 			}
// 		}
// 	}
// }

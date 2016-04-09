package eventHandler
import(
	"fmt"
	//"net"
	"time"
	"../network"
	"../message"
	"../elevatorControl/elevatorStatus"
	"../elevatorControl/orderHandling"
	"../cost"
	"../elevatorControl/driver"
)


func selectMaster(elevs map[string]*elevatorStatus.Elevator){
	//fmt.Println("select between: ", elevs)
	var minimumIP string = "129.241.187.999"
	for ip,_ := range elevs{
		fmt.Println("IP: ", ip)
		if ip < minimumIP {
			//fmt.Println("Ny master: ", ip)
			minimumIP = ip
		}
	}
	for _,elev := range elevs{
		//fmt.Println("Elev: ", elev, " Master: ", elev.Master)
		elev.Master = minimumIP
		//fmt.Println("Registrert master: ", elev.Master)
	} 
	fmt.Println("New master: ", minimumIP)
}

func deleteElevator(elevs map[string]*elevatorStatus.Elevator,IP string, sendChan chan message.UpdateMessage, elevatorTimers map[string]*time.Timer, myIP string, elevChan chan elevatorStatus.Elevator, abortElev chan bool){
	elev := elevs[IP]
	fmt.Println("IP: ", IP, "myIP: ", myIP)
	fmt.Println("delete this elevator: ", elev)
	delete(elevatorTimers, IP)
	delete(elevs, IP)
	selectMaster(elevs)
	//If the deleted elevator had any external orders they need to be reassigned
	if IP != myIP{
		for floor := 0; floor < driver.NUM_FLOORS; floor++{
			for button := 0; button < driver.NUM_BUTTONS-1; button++{
				if elev.OrderMatrix[floor][button] == 1{
					var order [2]int
					order[0] = floor
					order[1] = button
					fmt.Println("Order: ", order, "RecieverIP: ", elevs[IP].Master )
					sendChan <-message.UpdateMessage{MessageType: message.PlacedOrder, Order: order, RecieverIP: elevs[myIP].Master}
				}
			}
		}
	} else{
		// slett egne eksterne ordre
		e := <- elevChan
		for floor := 0; floor < driver.NUM_FLOORS; floor++{
			for button := 0; button < driver.NUM_BUTTONS-1; button++{
				e.OrderMatrix[floor][button] = 0
				driver.Set_button_lamp(button,floor,0)

			}
		}
		elevChan <-e
	}


	abortElev <- true

}
func EventHandler(newStateUpdate chan bool, buttonChan chan [2]int, powerChan chan bool, deleteChan chan [4]int, elevChan chan elevatorStatus.Elevator, sendNetwork chan message.UpdateMessage,  notAlive chan bool){
	
	for{
		select{
			case <-newStateUpdate:
				fmt.Println("Sender ny state update") 
				elev := message.MakeCopyOfElevator(elevChan)
				sendNetwork <-message.UpdateMessage{MessageType: message.StateUpdate, ElevatorStatus: elev}
			case order:= <-buttonChan:
				fmt.Println("sent new order on network")
				elev := message.MakeCopyOfElevator(elevChan)
				orderHandling.WriteInternals(elev.OrderMatrix)
				sendNetwork <-message.UpdateMessage{MessageType: message.PlacedOrder, RecieverIP: elev.Master, Order: order,
				ElevatorStatus: elev}
			case <-powerChan:
				notAlive <- true
				fmt.Println("Det har blitt lagt noe på powerChan")
				
			case deleted := <-deleteChan:
				elev := message.MakeCopyOfElevator(elevChan)
				orderHandling.WriteInternals(elev.OrderMatrix)
				sendNetwork <-message.UpdateMessage{MessageType: message.CompletedOrder, DelOrder: deleted, ElevatorStatus: elev}

		}
	}
}

 

func MessageHandler(recvChan chan message.UpdateMessage, sendChan chan message.UpdateMessage, newOrderToFSM chan elevatorStatus.Elevator, elevChan chan elevatorStatus.Elevator, abortElev chan bool){ // +order message.UpdateMessage
	var msg message.UpdateMessage
	// sjekk om mottatt melding er sent fra en av våre heiser, før det legges ut på channel til main
	
	myIP := network.GetIpAddress()
	var MasterMatrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int
	//fmt.Println("MasterMatrix: ",MasterMatrix)
	elevs := make(map[string]*elevatorStatus.Elevator)
	elevatorTimers := make(map[string]*time.Timer)
	//fmt.Println("Map: ", elevs)

	
	for{
		msg = <-recvChan
		msgType := msg.MessageType
		//fmt.Println("MAP :: %v", elevs)
		switch(msgType){
			case message.IAmAlive:
				var shouldAppend bool = true
				for ip,_ := range elevatorTimers{
					if ip == msg.ElevatorStatus.IP{
						//fmt.Println("1. IAmAlive")
						shouldAppend = false
						elevatorTimers[ip].Reset(time.Second*2)
						//fmt.Println("inne i IAmAlive", shouldAppend)
					}
				}	
				if shouldAppend == true{
					//fmt.Println("Oppdager heis for første gang: ")

					elevs[msg.ElevatorStatus.IP] = new(elevatorStatus.Elevator)
					elevs[msg.ElevatorStatus.IP].Dir = msg.ElevatorStatus.Dir
					elevs[msg.ElevatorStatus.IP].CurrentFloor = msg.ElevatorStatus.CurrentFloor
					elevs[msg.ElevatorStatus.IP].PreviousFloor = msg.ElevatorStatus.PreviousFloor
					elevs[msg.ElevatorStatus.IP].State = msg.ElevatorStatus.State
					elevs[msg.ElevatorStatus.IP].IP = msg.ElevatorStatus.IP
					
					
					

					for _,elev := range elevs{
						fmt.Println("Elevators in map: ", elev)
					}
					ip := msg.ElevatorStatus.IP
					elevatorTimers[msg.ElevatorStatus.IP] = time.AfterFunc(time.Second*2, func() { deleteElevator(elevs, ip, sendChan, elevatorTimers, myIP, elevChan, abortElev) })
					
					

					selectMaster(elevs)
					 
				}			

			case message.PlacedOrder:
				//for _,elev := range elevs{
				//		fmt.Println("Elevators in map: ", elev)
				//} 
				fmt.Println("fått placedorder")
				floor := msg.Order[0]
				button := msg.Order[1]
				//fmt.Println("Master før if: ", MasterMatrix[floor][button])
				if MasterMatrix[floor][button] == 0{
					//fmt.Println("floor: ", floor, "button: ", button)
					fmt.Println("Mastermatrix: ", MasterMatrix)
					//fmt.Println("master: ", elevs[msg.ElevatorStatus.IP].Master, " GetIpAddress: ",network.GetIpAddress() )
					if elevs[msg.ElevatorStatus.IP].Master == network.GetIpAddress(){ 
						
						// Kall kostfunksjon og legg bestillingen (+valgt heis) på en channel - mellomledd før nettverket tar bestillingen videre herfra?
						if button < 2{
							MasterMatrix[floor][button] = 1
							for _,elev := range elevs{
								fmt.Println("Elevators in map før kostfunksjon: ", elev)
							}
							AssignedElev := cost.AssignOrdersToElevator(msg, elevs)
							//fmt.Println("AssignedElev: ", AssignedElev)
							sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, RecieverIP: AssignedElev}
							//fmt.Println("Sender AssignedOrder fra master")
							//meldingmelding := <-sendChan
							//fmt.Println ("AssignedOrder:", meldingmelding)
								
								//ElevatorStatus: elevatorStatus.Elevator{RecieverIP: AssignedElev}}
						} else if button == 2{
							sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, RecieverIP: msg.ElevatorStatus.IP, Order: msg.Order}
								
								//ElevatorStatus: elevatorStatus.Elevator{RecieverIP: msg.ElevatorStatus.SenderIP}}
						}
					}
				}	


			case message.AssignedOrder:
				// Ta imot og legg til bestillinger for master (må ha en lik case i func Slave)
				//fmt.Println("RecieverIP - Assigned: ", msg.RecieverIP)
				
				button := msg.Order[1]
				if button < 2{
					sendChan <-message.UpdateMessage{MessageType: message.LightUpdate, Order: msg.Order}
				}
				if msg.RecieverIP == network.GetIpAddress(){
					*elevs[msg.RecieverIP] = orderHandling.AddOrderToQueue(*elevs[msg.RecieverIP], msg.Order)
					for _,elev := range elevs{
						fmt.Println("Elevators in map i AssignedOrder: ", elev)
					}
					newOrderToFSM <- *elevs[msg.RecieverIP]
					

				}
			case message.CompletedOrder:
				// Slett ordre i MasterMatrix
				completed := msg.DelOrder
				floor := msg.DelOrder[3]
				var noOrder [2]int 
				noOrder[0] = 0
				noOrder[1] = 1
				

				for button := 0; button < driver.NUM_BUTTONS; button++{
					if completed[button] == 1{
						MasterMatrix[floor][button] = 0
						if button < 2{
							sendChan <-message.UpdateMessage{MessageType: message.LightUpdate, DelOrder: completed, Order: noOrder}
						} else {
							sendChan <-message.UpdateMessage{MessageType: message.LightUpdate, RecieverIP: msg.ElevatorStatus.IP, DelOrder: completed, Order: noOrder}
						}
					}
				}

			case message.StateUpdate:
				//fmt.Println("Mottar StateUpdate: ", msg.ElevatorStatus.IP)

				for _,elev := range elevs{
					fmt.Println("Elevators in map in StateUpdate: ", elev)
				}
				if elevs[msg.ElevatorStatus.IP] != nil{
					elevs[msg.ElevatorStatus.IP].Dir = msg.ElevatorStatus.Dir
					elevs[msg.ElevatorStatus.IP].CurrentFloor = msg.ElevatorStatus.CurrentFloor
					elevs[msg.ElevatorStatus.IP].PreviousFloor = msg.ElevatorStatus.PreviousFloor
					elevs[msg.ElevatorStatus.IP].State = msg.ElevatorStatus.State
					elevs[msg.ElevatorStatus.IP].IP = msg.ElevatorStatus.IP
					elevs[msg.ElevatorStatus.IP].OrderMatrix = msg.ElevatorStatus.OrderMatrix
				}

			case message.LightUpdate:
				//To turn off buttonlights
				completed := msg.DelOrder
				floor := msg.DelOrder[3]
				for button := 0; button < driver.NUM_BUTTONS-1; button++{
					if completed[button] == 1{
						fmt.Println("floor i lights: ", floor)
						driver.Set_button_lamp(button,floor,0)
					}
				}
				f := msg.Order[0]
				b := msg.Order[1]
				fmt.Println("f: ", f, "b: ", b)
				if b < 2{
				 	driver.Set_button_lamp(b, f, 1)
				}


		}
	}
}

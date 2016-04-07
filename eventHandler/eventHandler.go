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

func deleteElevator(elevs map[string]*elevatorStatus.Elevator,IP string){
	delete(elevs,IP)
	selectMaster(elevs)
	//Husk å fordele ordre hos denne heisen til de resterende heisene
}


func MessageHandler(recvChan chan message.UpdateMessage, sendChan chan message.UpdateMessage, newOrderToFSM chan elevatorStatus.Elevator){ // +order message.UpdateMessage
	var msg message.UpdateMessage
	// sjekk om mottatt melding er sent fra en av våre heiser, før det legges ut på channel til main
	
	var MasterMatrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int
	//fmt.Println("MasterMatrix: ",MasterMatrix)
	elevs := make(map[string]*elevatorStatus.Elevator)
	elevatorTimers := make(map[string]*time.Timer)
	//fmt.Println("Map: ", elevs)

	
	for{
		msg = <-recvChan
		//fmt.Println("Recieved message: ", msg)
		msgType := msg.MessageType

	
		switch(msgType){
			case message.IAmAlive:
				var shouldAppend bool = true
				for ip,_ := range elevatorTimers{
					if ip == msg.ElevatorStatus.IP{
						//fmt.Println("1. IAmAlive")
						elevatorTimers[ip].Reset(time.Second*2)
						shouldAppend = false
						//fmt.Println("inne i IAmAlive", shouldAppend)
					}
				}	
				if shouldAppend == true{
					//fmt.Println("Oppdager heis for første gang: ")
					//var e elevatorStatus.Elevator// bør fungere, spørr mathias 

					elevs[msg.ElevatorStatus.IP] = new(elevatorStatus.Elevator)
					elevs[msg.ElevatorStatus.IP].Dir = msg.ElevatorStatus.Dir
					elevs[msg.ElevatorStatus.IP].CurrentFloor = msg.ElevatorStatus.CurrentFloor
					elevs[msg.ElevatorStatus.IP].PreviousFloor = msg.ElevatorStatus.PreviousFloor
					elevs[msg.ElevatorStatus.IP].State = msg.ElevatorStatus.State
					elevs[msg.ElevatorStatus.IP].IP = msg.ElevatorStatus.IP
					// e.Dir = msg.ElevatorStatus.Dir
					// e.CurrentFloor = msg.ElevatorStatus.CurrentFloor
					// e.PreviousFloor = msg.ElevatorStatus.PreviousFloor
					// e.State = msg.ElevatorStatus.State
					// e.IP = msg.ElevatorStatus.IP

					// elevs[msg.ElevatorStatus.IP] = &e

					
					

					for _,elev := range elevs{
						fmt.Println("Elevators in map: ", elev)
					} 
			
					elevatorTimers[msg.ElevatorStatus.IP] = time.AfterFunc(time.Second*2, func() {deleteElevator(elevs,msg.ElevatorStatus.IP)})
					selectMaster(elevs)
					//for _,elev := range elevs{
					//	fmt.Println("Elevators in map after master: ", elev)
					//} 
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
					//fmt.Println("Mastermatrix: ", MasterMatrix)
					fmt.Println("master: ", elevs[msg.ElevatorStatus.IP].Master, " GetIpAddress: ",network.GetIpAddress() )
					if msg.RecieverIP == network.GetIpAddress(){ 
						
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
							sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, RecieverIP: msg.ElevatorStatus.IP}
								
								//ElevatorStatus: elevatorStatus.Elevator{RecieverIP: msg.ElevatorStatus.SenderIP}}
						}
					}
				}	


			case message.AssignedOrder:
				// Ta imot og legg til bestillinger for master (må ha en lik case i func Slave)
				//fmt.Println("RecieverIP - Assigned: ", msg.RecieverIP)
				if msg.RecieverIP == network.GetIpAddress(){
					*elevs[msg.RecieverIP] = orderHandling.AddOrderToQueue(*elevs[msg.RecieverIP], msg.Order)
					for _,elev := range elevs{
						fmt.Println("Elevators in map i AssignedOrder: ", elev)
					}
					newOrderToFSM <- *elevs[msg.RecieverIP]
				}
			case message.CompletedOrder:
				// Slett ordre i MasterMatrix
				del := msg.DelOrder
				floor := msg.DelOrder[3]

				for button := 0; button < driver.NUM_BUTTONS; button++{
					if del[button] == 1{
						MasterMatrix[floor][button] = 0
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
				}

				// var e elevatorStatus.Elevator// bør fungere, spørr matias 

				// e.Dir = msg.ElevatorStatus.Dir
				// e.CurrentFloor = msg.ElevatorStatus.CurrentFloor
				// e.PreviousFloor = msg.ElevatorStatus.PreviousFloor
				// e.State = msg.ElevatorStatus.State
				// e.IP = msg.ElevatorStatus.IP

				// elevs[msg.ElevatorStatus.IP] = &e

		}
	}
}

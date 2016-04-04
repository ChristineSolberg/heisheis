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

type AllElevators struct{

	elevs map[string]*elevatorStatus.Elevator
	MasterMatrix [4][3]int
	
}


/*func InitMasterSlave(msgChan chan int, recvChan chan message.UpdateMessage, e elevatorStatus.Elevator, elevators AllElevators){
	//HJÆLP
	var recv int = 0
	var msg message.UpdateMessage
	timer := time.NewTimer(time.Second*2)
	counter := time.After(time.Second*2)

	for {
		msg = <-recvChan
		recv = msg.MessageType
		if recv != 0{
			msgChan <- recv
			break
		}
		select{
		case <-counter:
			break
		}
	}
	select{
		case <-timer.C:
			msg.ElevatorStatus.Master = network.GetIpAddress()
		case <-msgChan:
			sendNetwork <-message.UpdateMessage{MessageType: IAmAlive, RecieverIP: msg.ElevatorStatus.Master
					ElevatorStatus: elevatorStatus.Elevator{IP: network.GetIpAddress()}}
			
	}
}*/

func selectMaster(elevators AllElevators){
	fmt.Println("select between: ", elevators.elevs)
	var minimumIP string = "129.241.187.256"
	for ip,elev := range elevators.elevs{

		//if len(elevators.elevs) != 0{
		if ip < minimumIP {
			ip = elev.IP
		}
		//}
	}
	for _,elev := range elevators.elevs{
		elev.Master = minimumIP
	} 
	fmt.Println("New master: ", minimumIP)
}

func deleteElevator(elevators AllElevators,IP string){
	delete(elevators.elevs,IP)
	selectMaster(elevators)
}


func EventHandler(recvChan chan message.UpdateMessage, sendChan chan message.UpdateMessage, inToFSM chan elevatorStatus.Elevator){ // +order message.UpdateMessage
	var msg message.UpdateMessage
	//var elevators AllElevators
	elevators.elevs = make(map[string]*elevatorStatus.Elevator)
	elevatorTimers := make(map[string]*time.Timer)
	
	for{
		msg = <-recvChan
		fmt.Println("Recieved message: ", msg)
		msgType := msg.MessageType

	
		switch(msgType){
			case message.IAmAlive:
				var shouldAppend bool = true
				for ip,_ := range elevators.elevs{
					if ip == msg.ElevatorStatus.IP{
						elevatorTimers[ip].Reset(time.Second*2)
						shouldAppend = false
						fmt.Println("inne i IAmAlive", shouldAppend)
					}
				}	
				if shouldAppend == true{
					fmt.Println("Oppdager heis for første gang: ")
					elevators.elevs[msg.ElevatorStatus.IP] = new(elevatorStatus.Elevator)
					elevators.elevs[msg.ElevatorStatus.IP] = msg.ElevatorStatus
					fmt.Println("Elevators: ", elevators)
					elevatorTimers[msg.ElevatorStatus.IP] = time.AfterFunc(time.Second*2, func() {deleteElevator(elevators,msg.ElevatorStatus.IP)})
					selectMaster(elevators)
					
				}			

			case message.PlacedOrder:
				floor := msg.Order[0]
				button := msg.Order[1]
				elevators.MasterMatrix[floor][button] = 1
				if msg.ElevatorStatus.Master == network.GetIpAddress(){ 
					
					// Kall kostfunksjon og legg bestillingen (+valgt heis) på en channel - mellomledd før nettverket tar bestillingen videre herfra?
					if button < 2{
						AssignedElev := cost.AssignOrdersToElevator(msg, elevators.elevs)
						sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, RecieverIP: AssignedElev}
							
							//ElevatorStatus: elevatorStatus.Elevator{RecieverIP: AssignedElev}}
					} else if button == 2{
						sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, RecieverIP: msg.ElevatorStatus.IP}
							
							//ElevatorStatus: elevatorStatus.Elevator{RecieverIP: msg.ElevatorStatus.SenderIP}}
					}
				}


			case message.AssignedOrder:
				// Ta imot og legg til bestillinger for master (må ha en lik case i func Slave)
				if msg.RecieverIP == network.GetIpAddress(){
					*elevators.elevs[msg.RecieverIP] = orderHandling.AddOrderToQueue(*elevators.elevs[msg.RecieverIP])
						inToFSM <- *elevators.elevs[msg.RecieverIP]

				}
			case message.CompletedOrder:
				// Slett ordre i MasterMatrix
				del := msg.DelOrder
				floor := msg.DelOrder[3]

				for button := 0; button < driver.NUM_BUTTONS; button++{
					if del[button] == 1{
						elevators.MasterMatrix[floor][button] = 0
					}
				}

			case message.StateUpdate:
				*elevators.elevs[msg.ElevatorStatus.IP] = msg.ElevatorStatus

		}
	}
}

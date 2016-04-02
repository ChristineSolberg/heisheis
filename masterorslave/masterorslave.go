package masterorslave
import(
	//"fmt"
	//"net"
	"time"
	"../network"
	"../message"
	"../elevatorControl/elevatorStatus"
	"../elevatorControl/orderHandling"
	"../cost"
)



func InitMasterSlave(msgChan chan int, recvChan chan message.UpdateMessage, e elevatorStatus.Elevator){
	var recv int = 0
	var msg message.UpdateMessage
	timer := time.NewTimer(time.Second*2)

	for{
		msg = <-recvChan
		recv = msg.MessageType
		if recv != 0{
			msgChan <- recv
			break
		}
	}
	select{
		case <-timer.C:
			e.MasterOrSlave = message.Master
		case <-msgChan:
			e.MasterOrSlave = message.Slave
	}
}

func Master(recvChan chan message.UpdateMessage, sendChan chan message.UpdateMessage, inToFSM chan elevatorStatus.Elevator){ // +order message.UpdateMessage
	elevators := make(map[string]*elevatorStatus.Elevator)
	var msg message.UpdateMessage
	msg = <-recvChan
	msgType := msg.MessageType

	elevatorTimers := make(map[string]*time.Timer)

	//trengs det en channel for sending og en annen for mottak? (slik som i main)

	switch(msgType){
		case message.IAmAlive:
			// Starte teller hos master hver gang den får en IAmAlive. Hvis det har gått x antall sekunder uten IAmAlive - anta heisen er død og fjern den fra Elevators. 
			// Slaves skal høre etter fra master også. En av slavene skal bli master dersom master dør
			var shouldAppend bool = true
			if msg.ElevatorStatus.IP != network.GetIpAddress(){
				for ip,_ := range elevators{
					if ip == msg.ElevatorStatus.IP{
						elevatorTimers[ip].Reset(time.Second*2)
						shouldAppend = false
					}
				}	
				if shouldAppend == true{
					elevators[msg.ElevatorStatus.IP] = append(elevators[msg.ElevatorStatus.IP],msg.ElevatorStatus) // HER AVSLUTTET VI 02.04.16
					elevatorTimers[msg.ElevatorStatus.IP] = time.AfterFunc(time.Second*2, deleteElevator) //VI MÅ LAGE DELETEELEVATOR!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!:)
				}			

			}
		case message.PlacedOrder: 
			button := msg.Order[0]
			floor := msg.Order[1]
			msg.MasterMatrix[floor][button] = 1
			
			// Kall kostfunksjon og legg bestillingen (+valgt heis) på en channel - mellomledd før nettverket tar bestillingen videre herfra?
			if button < 2{
				AssignedElev := cost.AssignOrdersToElevator(msg, elevators) //-- Finn på nytt navn på channel - trenger vi channel her egentlig?
				sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, RecieverIP: AssignedElev}
					
					//ElevatorStatus: elevatorStatus.Elevator{RecieverIP: AssignedElev}}
			} else if button == 2{
				sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, RecieverIP: msg.ElevatorStatus.IP}
					
					//ElevatorStatus: elevatorStatus.Elevator{RecieverIP: msg.ElevatorStatus.SenderIP}}
			}


		case message.AssignedOrder:
			// Ta imot og legg til bestillinger for master (må ha en lik case i func Slave)
			if msg.RecieverIP == network.GetIpAddress(){
				elevators[msg.RecieverIP] = orderHandling.AddOrderToQueue(elevators[msg.RecieverIP])
					inToFSM <- elevators[msg.RecieverIP]

			}
		case message.CompletedOrder:
			// Slett ordre i MasterMatrix
			button := msg.Order[0]
			floor := msg.Order[1]
			msg.MasterMatrix[floor][button] = 0


		case message.StateUpdate:
			// Hver gang heisene endrer state - oppdater Elevators

	}
}

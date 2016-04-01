package masterorslave

import(
	//"fmt"
	"net"
	"time"
	"./network"
	"./message"
	"./elevatorControl/elevatorStatus"
	"./cost"
)



func InitMasterSlave(conn *net.UDPConn, msgChan chan int, recvChan chan message.UpdateMessage, e elevatorStatus.Elevator){
	var recv int = 0
	var msg message.UpdateMessage
	timer := time.NewTimer(time.Second*2)

	for{
		msg = message.RecvMsg(conn, recvChan)
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

func Master(conn *net.UDPConn, recvChan chan message.UpdateMessage, sendChan chan message.UpdateMessage){ // +order message.UpdateMessage
	elevators := make([]elevatorStatus.Elevator, 0)
	var msg message.UpdateMessage
	message.RecvMsg(conn, recvChan)
	msg = <-recvChan
	msgType := msg.MessageType

	//trengs det en channel for sending og en annen for mottak? (slik som i main)

	switch(msgType){
		case message.IAmAlive:
			// Starte teller hos master hver gang den får en IAmAlive. Hvis det har gått x antall sekunder uten IAmAlive - anta heisen er død og fjern den fra Elevators. 
			// Slaves skal høre etter fra master også. En av slavene skal bli master dersom master dør
		case message.PlacedOrder: 
			// Husk legge inn i MasterMatrix
			button := msg.Order[0]
			floor := msg.Order[1]
			msg.MasterMatrix[floor][button] = 1
			
			// Kall kostfunksjon og legg bestillingen (+valgt heis) på en channel - mellomledd før nettverket tar bestillingen videre herfra?
			if button < 2{
				AssignedElev := cost.AssignOrdersToElevator(msg, elevators) //-- Finn på nytt navn på channel - trenger vi channel her egentlig?
				sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, 
					ElevatorStatus: elevatorStatus.Elevator{RecieverIP: AssignedElev}}
			} else if button == 2{
				sendChan <-message.UpdateMessage{MessageType: message.AssignedOrder, Order: msg.Order, 
					ElevatorStatus: elevatorStatus.Elevator{RecieverIP: msg.ElevatorStatus.SenderIP}}
			}


		case message.AssignedOrder:
			// Ta imot og legg til bestillinger for master (må ha en lik case i func Slave)
			if msg.RecieverIP == network.GetIpAddress(){
				for _, elev := range elevators {
					if elev.ElevatorStatus.SenderIP == network.GetIpAddress(){
						elev = orderHandling.AddOrderToQueue(elev)
						// skal den oppdaterte elev sendes via channel tilbake til main?
					}
				}

			}
		case message.CompletedOrder:
			// Slett ordre i MasterMatrix
		case message.StateUpdate:
			// Hver gang heisene endrer state - oppdater Elevators
	}
}

package masterorslave

import(
	"fmt"
	"net"
	"time"
	"./network"
	"./message"
)



func InitMasterSlave(msgChan chan message.UpdateMessage, e elevatorStatus.Elevator){
	var recv int := 0
	var msg message.UpdateMessage
	timer := time.NewTimer(time.Second*2)

	for{
		msg = message.RecvMsg()
		recv = msg.MessageType
		if recv != 0{
			msgChan <- recv
			break
		}
	}
	select{
		case <- timer.C
			e.MasterOrSlave = message.Master
		case <- msgChan
			e.MasterOrSlave = message.Slave
	}
}

func Master(conn *net.UDPConn, msgChan chan UpdateMessage){ // +order message.UpdateMessage
	elevators := make(elevatorStatus.Elevator[], 0)
	var msg UpdateMessage
	message.RecvMsg(conn, msgChan)
	msg := <- msgChan
	msgType := msg.MessageType

	switch(message.MessageType){
		case msgType == IAmAlive:
			//
		case msgType == PlacedOrder: 
			// Husk legge inn i MasterMatrix
			// button := msg.NewOrder[0]
			// floor := msg.NewOrder[1]
			// msg.MasterMatrix[floor][button] = 1
			
			// Kall kostfunksjon og legg bestillingen (+valgt heis) på en channel - mellomledd før nettverket tar bestillingen videre herfra?
			// if button < 1{
			// 		AssignedElev := AssignOrdersToElevator(order, elevators, networkChan) //-- Finn på nytt navn på channel - trenger vi channel her egentlig?
			//		msgChan<- UpdateMessage{MessageType: AssignedOrder}
			// }


		case msgType == CompletedOrder:
			//
		case msgType == StateUpdate:
			// 
	}
}

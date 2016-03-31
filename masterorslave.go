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

func Master(conn *net.UDPConn, msgChan chan UpdateMessage){
	elevators := make(elevatorStatus.Elevator[], 0)
	var msg UpdateMessage
	message.RecvMsg(conn, msgChan)
	msg := <- msgChan
	msgType := msg.MessageType 

	switch(message.MessageType){
		case msgType == 1:
			//
		case msgType == 2: 
			// Husk legge inn i MasterMatrix
			// Kall kostfunksjon
		case msgType == 3:
			//
		case msgType == 4:
			// 
	}
}

package masterorslave

import(
	"fmt"
	"net"
	"time"
	"./network"
	"./message"
)



func InitMasterSlave() (int, int){
	var recv int := 0
	//timerChan := make(chan string)
	msgChan := make(chan string)

	var msg message.UpdateMessage
	timer := time.NewTimer(time.Second*2)

	for{
		msg := message.RecvMsg()
		recv := msg.MessageType
		if recv != 0{
			msgChan <- recv
			break
		}
	}

	var master int := 0
	var slave int := 0


	select{
		case <- timer.C
			// bli master
			master = 1

		case <- msgChan
			//bli slave
			slave = 1

	}

	return master, slave
}

func Master(conn *net.UDPConn, msgChan chan UpdateMessage){
	var msg UpdateMessage
	message.RecvMsg(conn, msgChan)
	msg := <- msgChan
	msgType := msg.MessageType 

	switch{
		case msgType == 1:
			//
		case msgType == 2: 
			// Kall kostfunksjon
		case msgType == 3:
			//
	}
}
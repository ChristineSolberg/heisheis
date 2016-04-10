package message

import(
	"net"
	"fmt"
	"time"
	"encoding/json"
	"../../network"
	"../../elevatorControl/elevatorStatus"
	)

const(
	IAmAlive 		= 1
	PlacedOrder 	= 2
	AssignedOrder 	= 3
	CompletedOrder 	= 4
	StateUpdate 	= 5
	LightUpdate 	= 6
)

type UpdateMessage struct{
	MessageType int
	RecieverIP string
	Order [2] int   //[button, floor]
	DelOrder [4] int
	ElevatorStatus elevatorStatus.Elevator
}

func RecvMsg(conn *net.UDPConn, msgChan chan UpdateMessage) {
	buffer := make([]byte, 1024) 
	for{
		var msg UpdateMessage
		msgSize := network.UDPListen(conn,buffer)
		array := buffer[0:msgSize]
		err := json.Unmarshal(array, &msg)
		if err == nil{
			msgChan <- msg
		}
	}
	if (conn != nil){
		defer conn.Close()
	}
}

func SendMsg(conn *net.UDPConn, msgChan chan UpdateMessage, elevChan chan elevatorStatus.Elevator, notAlive chan bool){
	if (conn != nil){
		defer conn.Close()
	}
	ticker := time.NewTicker(time.Millisecond*500)
	var shouldSend bool = false
	if network.GetIpAddress() != "::1"{
		for {
			select{
			case message := <-msgChan:
				encoded := encodeUDPmsg(message)
				buffer := []byte(encoded)
				network.UDPWrite(conn, buffer)
			case <- notAlive:
				shouldSend = true
			case <-ticker.C:
				if shouldSend != true{
					e := elevatorStatus.MakeCopyOfElevator(elevChan)
					alive := UpdateMessage{MessageType: IAmAlive, ElevatorStatus: e}
					encoded := encodeUDPmsg(alive)
					buffer := []byte(encoded)
					network.UDPWrite(conn, buffer)
				}
			}
		}
	}
}

func encodeUDPmsg(message UpdateMessage)[]byte{
	encoded,err := json.Marshal(message)
	if err != nil{
		fmt.Println("error in encodeUDPmsg", err)
	}
	return encoded
}
package message

import(
	"net"
	"fmt"
	"time"
	"encoding/json"
	"../network"
	"../elevatorControl/elevatorStatus"

	)

const(
	IAmAlive = 1
	PlacedOrder = 2
	AssignedOrder = 3
	CompletedOrder = 4
	StateUpdate = 5
	LightUpdate = 6
)

type UpdateMessage struct{
	MessageType int
	RecieverIP string
	Order [2] int   // [button, floor]
	DelOrder [4] int
	ElevatorStatus elevatorStatus.Elevator
	CheckSum float64
	
}

func RecvMsg(conn *net.UDPConn, msgChan chan UpdateMessage) {
	// må kjøre serverConnection() for denne funksjonen kjøres
	buffer := make([]byte, 1024) 
	
	for{
		var msg UpdateMessage
		size := network.UDPListen(conn,buffer)
		//fmt.Println("size i recvmsg:", size)
		array := buffer[0:size]
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
	// må kjøre clientConnection() for denne funksjonen kjøres
	if (conn != nil){
		defer conn.Close()
	}
	ticker := time.NewTicker(time.Millisecond*500)
	var dontSend bool = false
	if network.GetIpAddress() != "::1"{
		e := MakeCopyOfElevator(elevChan)
		for {
			select{
			case message := <-msgChan:
				encoded := encodeUDPmsg(message)
				buffer := []byte(encoded)
				network.UDPWrite(conn, buffer)
			
			case <- notAlive:
				dontSend = true
				fmt.Println("Du blir snart slått av, notAlive")
			case <-ticker.C:
				if dontSend != true{
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
		fmt.Println("error", err)
	}
	return encoded
}
  
func MakeCopyOfElevator(elevChan chan elevatorStatus.Elevator)elevatorStatus.Elevator{
	e := <- elevChan
	elevChan <- e
	return e
}


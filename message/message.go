package message

import(
	"net"
	"fmt"
	//"time"
	"../network"
	"encoding/json"

	)

const(
	IAmAlive = 1
	PlacedOrder = 2
	StateUpdate = 3

)


const(
	Master		= 1
	Slave 		= 0
)

const(
	Down 		= -1
	Idle 		= 0
	Up 			= 1
)

const(
	First		= 1
	Second		= 2
	Third		= 3
	Fourth		= 4
)




type UpdateMessage struct{
	MessageType int
	ElevatorId int
	MasterOrSlave int
	NewOrder [2] float64   // [button, floor]
	CurrentState [2] float64 //[current floor, current direction]
	OrderMatrix [3][4] int
}

func RecvMsg(conn *net.UDPConn, msgChan chan UpdateMessage) UpdateMessage{
	// må kjøre serverConnection() for denne funksjonen kjøres
	buffer := make([]byte, 1024) 
	var msg UpdateMessage

	for{
		fmt.Println("inne i forløkken i recv")
		size := network.UDPListen(conn,buffer)
		fmt.Println(size)
		array := buffer[0:size]
		err := json.Unmarshal(array, &msg)
		if err == nil{
			msgChan <- msg
		}
	}
	fmt.Println("Recv2")
	defer conn.Close()

	return msg
}

func SendMsg(conn *net.UDPConn, msgChan chan UpdateMessage){
	// må kjøre clientConnection() for denne funksjonen kjøres
	defer conn.Close()
	for {
		encoded,_ := json.Marshal(<-msgChan)
		
		buf := []byte(encoded)
		network.UDPWrite(conn, buf)
		fmt.Println("Alive")

		fmt.Println("Send ferdig")
	}
 
	// enc,_ := json.Marshal(msg)
	// buffer := []byte(enc)
	// network.UDPWrite(conn, buffer)

	// ny kode med channel
	

	// for{
	// 	enc, err := json.Marshal(<-msgChan)
	// 	if err == nil{
	// 		network.UDPWrite(conn, enc)
	// 	}
	//}

	
}


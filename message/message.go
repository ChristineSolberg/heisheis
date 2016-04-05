package message

import(
	"net"
	"fmt"
	//"time"
	"../network"
	"encoding/json"
	"../elevatorControl/elevatorStatus"

	)

const(
	IAmAlive = 1
	PlacedOrder = 2
	AssignedOrder = 3
	CompletedOrder = 4
	StateUpdate = 5
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
	RecieverIP string
	Order [2] int   // [button, floor]
	DelOrder [4] int
	ElevatorStatus elevatorStatus.Elevator
	
}

func RecvMsg(conn *net.UDPConn, msgChan chan UpdateMessage) {
	// må kjøre serverConnection() for denne funksjonen kjøres
	buffer := make([]byte, 1024) 
	
	

	for{
		var msg UpdateMessage
		size := network.UDPListen(conn,buffer)
		fmt.Println("size i recvmsg:", size)
		array := buffer[0:size]
		err := json.Unmarshal(array, &msg)
		if err == nil{
			msgChan <- msg
		}
		fmt.Println("Mottatt melding: ", msg)
	}
	fmt.Println("Recv2")
	defer conn.Close()

	//return msg
}

func SendMsg(conn *net.UDPConn, msgChan chan UpdateMessage){
	// må kjøre clientConnection() for denne funksjonen kjøres
	defer conn.Close()
	for {
		v := <-msgChan
		fmt.Println("Melding via network: ",v)
		encoded,err := json.Marshal(v)
		fmt.Println("error", err)
		buf := []byte(encoded)
		network.UDPWrite(conn, buf)

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

func MessageManager(toElev chan UpdateMessage, fromElev chan UpdateMessage){
	conn1 := network.ServerConnection()
	conn2 := network.ClientConnection()
	
	go RecvMsg(conn1,toElev)
	go SendMsg(conn2,fromElev)


    // sjekk om mottatt melding er sent fra en av våre heiser, før det legges ut på channel til main
     

}
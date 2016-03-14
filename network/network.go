package network

import(
	"net"
	"fmt"
	)




func ServerConnection()*net.UDPConn{ //er det bra funksjonsnavn??
	port := ":20012"
	udpAddress, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// error handling here 

	conn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return conn
}

func ClientConnection() *net.UDPConn{
	port := ":20012"
	serverAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255" + port)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	conn, err := net.DialUDP("udp",nil,serverAddress)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	return conn
	// when to use conn.close()??
}



func UDPListen(conn *net.UDPConn, buffer []byte) int{
	fmt.Println("Waiting for msg")
	size,_,err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("YEEEEES")	
	return size
}

func UDPWrite(conn *net.UDPConn, buffer []byte){
	fmt.Println("Sending msg")
	conn.Write(buffer)
}
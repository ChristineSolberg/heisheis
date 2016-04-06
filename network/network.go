package network

import(
	"net"
	"fmt"
	"strings"
	)




func ServerConnection()*net.UDPConn{ //er det bra funksjonsnavn??
	port := ":20023"
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
	port := ":20023"
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
	//fmt.Println("Waiting for msg")
	size,_,err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	//fmt.Println("size: ", size)	
	return size
}

func UDPWrite(conn *net.UDPConn, buffer []byte){
	//fmt.Println("Sending msg")
	conn.Write(buffer)
}

func GetIpAddress()string{
	ipAdd,_ := net.InterfaceAddrs() 
    ip:=strings.Split(ipAdd[1].String(),"/")[0]
    return ip
}




package network

import(
	"net"
	"fmt"
	"strings"
	)

func ServerConnection()*net.UDPConn{ 
	port := ":5555"
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
	port := ":5555"
	serverAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255" + port)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	conn, err := net.DialUDP("udp",nil,serverAddress)
	if err != nil {
		fmt.Println("Error: ", err, conn)
	}

	return conn
	// when to use conn.close()??
}



func UDPListen(conn *net.UDPConn, buffer []byte) int{
	//fmt.Println("Waiting for msg")
	size,_,err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error: ", err)
		conn.Close()
		newConn := ServerConnection()
		if (newConn != nil){
			*conn = *newConn
		}
	}
	//fmt.Println("size: ", size)	
	return size
}

func UDPWrite(conn *net.UDPConn, buffer []byte){
	//fmt.Println("Sending msg")
	_, err := conn.Write(buffer)
	if err != nil {
		fmt.Println("Error: ", err)
		conn.Close()
		newConn := ClientConnection()
		if (newConn != nil){
			*conn = *newConn
		}
	}
}

func GetIpAddress()string{
	ipAdd,_ := net.InterfaceAddrs() 
    ip:=strings.Split(ipAdd[1].String(),"/")[0]
    return ip
}




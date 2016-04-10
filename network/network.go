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
		fmt.Println("Error in ServerConnection: ", err)
	}
	conn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		fmt.Println("Error in ServerConnection: ", err)
	}
	return conn
}

func ClientConnection() *net.UDPConn{
	port := ":5555"
	serverAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255" + port)
	if err != nil {
		fmt.Println("Error in ClientConnection: ", err)
	}
	conn, err := net.DialUDP("udp",nil,serverAddress)
	if err != nil {
		fmt.Println("Error in ClientConnection: ", err, conn)
	}
	return conn
}

func UDPListen(conn *net.UDPConn, buffer []byte) int{
	size,_,err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error in UDPListen: ", err)
		conn.Close()
		newConn := ServerConnection()
		if (newConn != nil){
			*conn = *newConn
		}
	}
	return size
}

func UDPWrite(conn *net.UDPConn, buffer []byte){
	_, err := conn.Write(buffer)
	if err != nil {
		fmt.Println("Error in UDPWrite: ", err)
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
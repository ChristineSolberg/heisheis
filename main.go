package main

import (
    "./elevatorControl/driver"
    "./elevatorControl"
    "./elevatorControl/orderHandling"
    "./elevatorControl/elevatorStatus"

    //"./network"
    //"./message"
    "fmt"
    //"time"
)



func main() {
	//Kjører en heis
	driver.Init()
	var e elevatorStatus.Elevator // er det godkjent med globale variabler i main?
	e = elevatorControl.StartUp(e)
	for{
		e = orderHandling.AddOrderToQueue(e)
		e = elevatorControl.UpdateFSM(e)
	}


// Tester nettverksmodulen
	// fmt.Println("Start main")
 //    recvChan := make(chan message.UpdateMessage)
 //    sendChan := make(chan message.UpdateMessage)

 //    conn1 := network.ServerConnection()
 //    conn2 := network.ClientConnection()
 //    fmt.Println("main2")


 //    var Alive message.UpdateMessage
 //    ticker := time.NewTicker(time.Millisecond*500)
	// fmt.Println("ticker started")
	// Alive.MessageType = message.IAmAlive
    
 //    go message.RecvMsg(conn1,recvChan)
 //    go message.SendMsg(conn2,sendChan)
	// for{
	// 	select{
	// 		case <-ticker.C:
	// 			fmt.Println("Legger på kanalen")
	// 			sendChan<-Alive

	// 		case msg := <-recvChan:
	// 			fmt.Println(msg)
	// 	}
	// }
 


    
    fmt.Println("main3")
    


}
package main

import (
    "./elevatorControl/driver"
    "./elevatorControl"
    "./elevatorControl/orderHandling"
    "./elevatorControl/elevatorStatus"
    "./eventHandler"

    "./network"
    "./message"
    //"fmt"
    //"net"
    //"time"
)

//Når skal vi "kill ourselves". Sverre mente vel dette var viktig?

func main() {
	//Initiating elevator
	elevChan := make(chan elevatorStatus.Elevator)
	driver.Init()
	elevatorControl.StartUp(elevChan)
	//fmt.Println("StartUp values: ", e)
	//fmt.Println("På channel: ", <-elevChan)

	
	// Initiating network
	recvNetwork := make(chan message.UpdateMessage, 100)
	sendNetwork := make(chan message.UpdateMessage, 100)
	notAlive := make(chan bool)
	conn1 := network.ServerConnection()
	conn2 := network.ClientConnection()
	go message.RecvMsg(conn1,recvNetwork)
	go message.SendMsg(conn2,sendNetwork,elevChan, notAlive)


	newOrderToFSM := make(chan elevatorStatus.Elevator,100)
	newStateUpdate := make(chan bool,100)
	buttonChan := make(chan [2]int, 100)
	powerChan := make(chan bool)
	abortElev := make(chan bool)
	deleteChan := make(chan [4]int, 100)

	go eventHandler.MessageHandler(recvNetwork,sendNetwork,newOrderToFSM, elevChan ,abortElev)
	go eventHandler.EventHandler(newStateUpdate, buttonChan, powerChan, deleteChan, elevChan, sendNetwork, notAlive)




	go orderHandling.ReadButtons(buttonChan, elevChan)
	
	// må lage en updatemessage med knappetrykket som sendes over nettverket til master

	//FSM

	
	elevatorControl.UpdateFSM(newOrderToFSM,newStateUpdate,deleteChan,elevChan,powerChan, abortElev)

	
}



	//Kjører en heis
	// driver.Init()
	// var e elevatorStatus.Elevator // er det godkjent med globale variabler i main?
	// e = elevatorControl.StartUp(e)
	// for{
	// 	e = orderHandling.AddOrderToQueue(e)
	// 	e = elevatorControl.UpdateFSM(e)
	
	// }


 //Tester nettverksmodulen
	// fmt.Println("Start main")
 //    recvChan := make(chan message.UpdateMessage)
 //    sendChan := make(chan message.UpdateMessage)
	


 //    conn1 := network.ServerConnection()
 //    conn2 := network.ClientConnection()
 //    fmt.Println("main2")


 //    //var Alive message.UpdateMessage
 //    ticker := time.NewTicker(time.Millisecond*500)
	// fmt.Println("ticker started")
	// //Alive.MessageType = message.IAmAlive
	// //Alive.Order[1] = 10
    

 //    //fmt.Println(Alive)
 //    go message.RecvMsg(conn1,recvChan)
 //    go message.SendMsg(conn2,sendChan)
	// for{
	// 	select{
	// 		case <-ticker.C:
	// 			fmt.Println("Legger på kanalen")
	// 			sendChan<- message.UpdateMessage{MessageType: message.IAmAlive}

	// 		case msg := <-recvChan:
	// 			fmt.Println("mottatt mld: " ,msg)
	// 	}
	// }
 
package main

import (
    "./elevatorControl/driver"
    "./elevatorControl"
    "./elevatorControl/orderHandling"
    "./elevatorControl/elevatorStatus"
    "./eventHandler"

    "./network"
    "./message"
    "fmt"
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
	conn1 := network.ServerConnection()
	conn2 := network.ClientConnection()
	go message.RecvMsg(conn1,recvNetwork)
	go message.SendMsg(conn2,sendNetwork,elevChan)


	newOrderToFSM := make(chan elevatorStatus.Elevator,100)
	newStateUpdate := make(chan bool,100)
	go eventHandler.MessageHandler(recvNetwork,sendNetwork,newOrderToFSM)


	buttonChan := make(chan [2]int, 20)
	go orderHandling.ReadButtons(buttonChan)
	
	// må lage en updatemessage med knappetrykket som sendes over nettverket til master

	//FSM
	deleteChan := make(chan [4]int, 10)
	go elevatorControl.UpdateFSM(newOrderToFSM,newStateUpdate,deleteChan,elevChan)

	for{
		select{
			case <-newStateUpdate:
				fmt.Println("Sender ny state update") 
				elev := message.MakeCopyOfElevator(elevChan)
				fmt.Println("Elev in new state update: ", elev/*.IP*/)
				sendNetwork <-message.UpdateMessage{MessageType: message.StateUpdate, ElevatorStatus: elev}
			case order:= <-buttonChan:
				fmt.Println("sent new order on network")
				elev := message.MakeCopyOfElevator(elevChan)
				sendNetwork <-message.UpdateMessage{MessageType: message.PlacedOrder, RecieverIP: elev.Master, Order: order,
				ElevatorStatus: elev}
				
			case deleted := <-deleteChan:
				elev := message.MakeCopyOfElevator(elevChan)
				sendNetwork <-message.UpdateMessage{MessageType: message.CompletedOrder, DelOrder: deleted, ElevatorStatus: elev}

		}
	}
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
 
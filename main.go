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
    "time"
)

//Når skal vi "kill ourselves". Sverre mente vel dette var viktig?

func main() {




	// Pseudokode for hvordan main skal gå:

	var e elevatorStatus.Elevator
	
	driver.Init()

	e = elevatorControl.StartUp(e)

	
	// initialiser nettverksdel her

	recvNetwork := make(chan message.UpdateMessage, 10)
	sendNetwork := make(chan message.UpdateMessage, 10)

	go Alive(sendNetwork) 
	go message.MessageManager(recvNetwork,sendNetwork)
	




	inToFSM := make(chan elevatorStatus.Elevator,10)
	outOfFSM := make(chan elevatorStatus.Elevator,10)

	go eventHandler.EventHandler(recvNetwork,sendNetwork,inToFSM)


	buttonChan := make(chan [2]int, 30)
	go orderHandling.ReadButtons(buttonChan,e)
	
	// må lage en updatemessage med knappetrykket som sendes over nettverket til master

	//FSM
	deleteChan := make(chan [4]int, 10)

	go elevatorControl.UpdateFSM(e,inToFSM,outOfFSM, deleteChan)

	


	for{
		select{
			case <-outOfFSM:
				e = <- outOfFSM
				sendNetwork <- message.UpdateMessage{MessageType: message.StateUpdate, ElevatorStatus: e}
			case <- buttonChan:
				sendNetwork <-message.UpdateMessage{MessageType: message.PlacedOrder, Order: <-buttonChan,
				ElevatorStatus: elevatorStatus.Elevator{IP: network.GetIpAddress()}}
			case <- deleteChan:
				sendNetwork <- message.UpdateMessage{MessageType: message.CompletedOrder, DelOrder: <- deleteChan, ElevatorStatus: e}

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
 

 }


func Alive(sendNetwork chan message.UpdateMessage){
	ticker := time.NewTicker(time.Millisecond*500)
    
	for{
		select{
			case <-ticker.C:
				sendNetwork <-message.UpdateMessage{MessageType: message.IAmAlive,
					ElevatorStatus: elevatorStatus.Elevator{IP: network.GetIpAddress()}}
		}
	}
}
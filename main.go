package main

import (
    "./elevatorControl/driver"
    "./elevatorControl"
    "./elevatorControl/orderHandling"
    "./elevatorControl/elevatorStatus"
    "./eventHandler"
    "./messageHandler"
    "./network"
    "./messageHandler/message"
)

func main() {
	//Initiating elevator
	elevChan := make(chan elevatorStatus.Elevator)
	driver.Init()
	elevatorControl.StartUp(elevChan)

	//Initiating network
	recvNetwork := make(chan message.UpdateMessage, 100)
	sendNetwork := make(chan message.UpdateMessage, 100)
	notAlive := make(chan bool)
	
	conn1 := network.ServerConnection()
	conn2 := network.ClientConnection()
	go message.RecvMsg(conn1,recvNetwork)
	go message.SendMsg(conn2,sendNetwork,elevChan, notAlive)

	//Running elevator
	newOrderToFSM := make(chan elevatorStatus.Elevator,100)
	newStateUpdate := make(chan bool,100)
	buttonChan := make(chan [2]int, 100)
	deleteChan := make(chan [4]int, 100)
	powerChan := make(chan bool)
	shouldAbort := make(chan bool)
	abortElev := make(chan bool)

	go messageHandler.MessageHandler(recvNetwork,sendNetwork,newOrderToFSM, elevChan ,abortElev, shouldAbort)
	go eventHandler.EventHandler(newStateUpdate, buttonChan, powerChan, deleteChan, elevChan, sendNetwork, notAlive, shouldAbort)
	go orderHandling.ReadButtons(buttonChan, elevChan)
	elevatorControl.RunFSM(newOrderToFSM,newStateUpdate,deleteChan,elevChan,powerChan, abortElev)	
}
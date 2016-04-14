package main

import (
	"./elevatorControl"
	"./elevatorControl/driver"
	"./elevatorControl/elevatorStatus"
	"./elevatorControl/orderHandling"
	"./eventHandler"
	"./messageHandler"
	"./messageHandler/message"
	"./network"
)

func main() {
	//Initiating elevator
	elevObject := make(chan elevatorStatus.Elevator)
	driver.Init()
	elevatorControl.StartUp(elevObject)

	//Initiating network
	recvNetwork := make(chan message.UpdateMessage, 100)
	sendNetwork := make(chan message.UpdateMessage, 100)
	notAlive := make(chan bool)

	conn1 := network.ServerConnection()
	conn2 := network.ClientConnection()
	go message.RecvMsg(conn1, recvNetwork)
	go message.SendMsg(conn2, sendNetwork, elevObject, notAlive)

	//Running elevator
	newOrderToFSM := make(chan elevatorStatus.Elevator, 100)
	newStateUpdate := make(chan bool, 100)
	buttonPushed := make(chan [2]int, 100)
	placedOrder := make(chan message.UpdateMessage, 100)
	setExternalLightsOn := make(chan [2]int, 50)
	setExternalLightsOff := make(chan [4]int, 50)
	deleteCompletedOrder := make(chan [4]int, 30)
	elevatorMap := make(chan map[string]*elevatorStatus.Elevator)
	powerOffDetected := make(chan bool)
	shouldAbort := make(chan bool)
	abortElev := make(chan bool)

	go messageHandler.MessageHandler(recvNetwork, sendNetwork, newOrderToFSM, elevObject, placedOrder, elevatorMap, setExternalLightsOn, setExternalLightsOff, abortElev, shouldAbort)
	go eventHandler.EventHandler(newStateUpdate, buttonPushed, powerOffDetected, deleteCompletedOrder, elevObject, placedOrder, elevatorMap, setExternalLightsOn, setExternalLightsOff, sendNetwork, notAlive, shouldAbort)
	go orderHandling.ReadButtons(buttonPushed, elevObject)
	elevatorControl.RunFSM(newOrderToFSM, newStateUpdate, deleteCompletedOrder, elevObject, powerOffDetected, abortElev)
}

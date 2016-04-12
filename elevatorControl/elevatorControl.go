package elevatorControl

import (
	"../network"
	"./driver"
	"./elevatorStatus"
	"./orderHandling"
	"fmt"
	"time"
)

func RunFSM(newOrderToFSM chan elevatorStatus.Elevator, newStateUpdate chan bool, deleteCompletedOrder chan [4]int, elevObject chan elevatorStatus.Elevator, powerOffDetected chan bool, abortElev chan bool) {
	doorTimer := time.NewTimer(time.Nanosecond)
	DoorTimeout := doorTimer.C
	powerTimer := time.NewTimer(time.Nanosecond)
	PowerTimeout := powerTimer.C
	eventChan := make(chan elevatorStatus.Event, 100)
	go getNextEvent(elevObject, eventChan, powerOffDetected, PowerTimeout, DoorTimeout)
	var turnOffPower bool
	for {
		select {
		case newOrder := <-newOrderToFSM:
			orderHandling.UpdateOrderMatrix(newOrder.OrderMatrix, elevObject)
		case turnOffPower = <-abortElev:
			fmt.Println("Detected power off, elevator will be turned off shortly")
		case event := <-eventChan:
			state := getElevatorState(elevObject)
			switch state {
			case elevatorStatus.IDLE:
				fmt.Println("State: IDLE")
				updateFSM_IDLE(event, elevObject, deleteCompletedOrder, powerTimer, doorTimer)
				break
			case elevatorStatus.GO_UP:
				fmt.Println("State: GO UP")
				updateFSM_GO_UP(event, elevObject, deleteCompletedOrder, powerOffDetected, powerTimer, doorTimer)
				break
			case elevatorStatus.GO_DOWN:
				fmt.Println("State: GO DOWN")
				updateFSM_GO_DOWN(event, elevObject, deleteCompletedOrder, powerOffDetected, powerTimer, doorTimer)
				break
			case elevatorStatus.DOOR_OPEN:
				fmt.Println("State: DOOR OPEN")
				updateFSM_DOOR_OPEN(event, elevObject, deleteCompletedOrder, doorTimer)
				break
			default:
				fmt.Println("Error: No valid state in UpdateFSM")
			}
			newStateUpdate <- true
		}
		if turnOffPower == true && network.GetIpAddress() != "::1" {
			fmt.Println("The power is off: Shutting down elevator")
			break
		}
	}
}

func getElevatorState(elevObject chan elevatorStatus.Elevator) elevatorStatus.State {
	e := <-elevObject
	state := e.State
	elevObject <- e
	return state
}

func updateFSM_IDLE(event elevatorStatus.Event, elevObject chan elevatorStatus.Elevator, deleteCompletedOrder chan [4]int, powerTimer *time.Timer, doorTimer *time.Timer) {
	e := <-elevObject
	switch event {
	case elevatorStatus.NEXT_ORDER:
		getNextDirection(&e)
		driver.Set_motor_speed(e.Direction)
		if e.Direction == driver.MDIR_UP {
			e.State = elevatorStatus.GO_UP
			powerTimer.Reset(time.Second * 20)
		} else if e.Direction == driver.MDIR_DOWN {
			e.State = elevatorStatus.GO_DOWN
			powerTimer.Reset(time.Second * 20)
		} else {
			e.State = elevatorStatus.IDLE
		}
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		driver.Set_door_open_lamp(1)
		doorTimer.Reset(time.Second * 3)
		orderHandling.DeleteCompletedOrders(&e, deleteCompletedOrder)
		e.State = elevatorStatus.DOOR_OPEN
		break
	case elevatorStatus.TIMER_OUT:
		driver.Set_door_open_lamp(0)
	case elevatorStatus.NO_EVENT:
		//do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_IDLE!")
	}
	elevObject <- e
}

func updateFSM_GO_UP(event elevatorStatus.Event, elevObject chan elevatorStatus.Elevator, deleteCompletedOrder chan [4]int, powerOffDetected chan bool, powerTimer *time.Timer, DoorTimeout *time.Timer) {
	e := <-elevObject
	switch event {
	case elevatorStatus.FLOOR_REACHED:
		if orderHandling.ShouldStop(e) == true {
			driver.Set_motor_speed(driver.MDIR_STOP)
			driver.Set_door_open_lamp(1)
			DoorTimeout.Reset(time.Second * 3)
			powerTimer.Stop()
			orderHandling.DeleteCompletedOrders(&e, deleteCompletedOrder)
			e.PreviousFloor = driver.Get_floor_sensor_signal()
			getNextDirection(&e)
			e.State = elevatorStatus.DOOR_OPEN
		} else {
			e.State = elevatorStatus.GO_UP
		}
		break
	case elevatorStatus.POWER_OFF:
		fmt.Println("Power off detected")
		powerOffDetected <- true
	case elevatorStatus.NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_GO_UP!")
	}
	elevObject <- e
}

func updateFSM_GO_DOWN(event elevatorStatus.Event, elevObject chan elevatorStatus.Elevator, deleteCompletedOrder chan [4]int, powerOffDetected chan bool, powerTimer *time.Timer, DoorTimeout *time.Timer) {
	e := <-elevObject
	switch event {
	case elevatorStatus.FLOOR_REACHED:
		if orderHandling.ShouldStop(e) == true {
			driver.Set_motor_speed(driver.MDIR_STOP)
			driver.Set_door_open_lamp(1)
			e.State = elevatorStatus.DOOR_OPEN
			DoorTimeout.Reset(time.Second * 3)
			powerTimer.Stop()
			orderHandling.DeleteCompletedOrders(&e, deleteCompletedOrder)
			e.PreviousFloor = driver.Get_floor_sensor_signal()
			getNextDirection(&e)
		} else {
			e.State = elevatorStatus.GO_DOWN
		}
		break
	case elevatorStatus.POWER_OFF:
		powerOffDetected <- true
	case elevatorStatus.NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_GO_DOWN!")
	}
	elevObject <- e
}

func updateFSM_DOOR_OPEN(event elevatorStatus.Event, elevObject chan elevatorStatus.Elevator, deleteCompletedOrder chan [4]int, DoorTimeout *time.Timer) {
	e := <-elevObject
	switch event {
	case elevatorStatus.TIMER_OUT:
		driver.Set_door_open_lamp(0)
		e.State = elevatorStatus.IDLE
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		DoorTimeout.Reset(time.Second * 3)
		driver.Set_door_open_lamp(1)
		e.State = elevatorStatus.DOOR_OPEN
		orderHandling.DeleteCompletedOrders(&e, deleteCompletedOrder)
		break
	case elevatorStatus.NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_DOOR_OPEN!")
	}
	elevObject <- e
}

func getNextEvent(elevObject chan elevatorStatus.Elevator, eventChan chan elevatorStatus.Event, powerOffDetected chan bool, PowerTimeout <-chan time.Time, DoorTimeout <-chan time.Time) {
	var nextEvent elevatorStatus.Event
	var prevEvent elevatorStatus.Event
	for {
		e := <-elevObject
		go changeElev(e, elevObject)
		currentFloor := driver.Get_floor_sensor_signal()
		select {
		case <-DoorTimeout:
			nextEvent = elevatorStatus.TIMER_OUT
			if prevEvent != nextEvent {
				fmt.Println("Event: ", nextEvent)
				eventChan <- nextEvent
				prevEvent = nextEvent
			}
		case <-PowerTimeout:
			nextEvent = elevatorStatus.POWER_OFF
			eventChan <- nextEvent
		default:
			if (currentFloor != -1) && (e.PreviousFloor != currentFloor) {
				nextEvent = elevatorStatus.FLOOR_REACHED
			} else if ((e.State == elevatorStatus.IDLE) || (e.State == elevatorStatus.DOOR_OPEN)) && orderHandling.NewOrderAtCurrentFloor(e) == true {
				nextEvent = elevatorStatus.NEW_ORDER_AT_CURRENT
			} else if orderHandling.LengthOfQueue(e) > 0 && e.State == elevatorStatus.IDLE {
				nextEvent = elevatorStatus.NEXT_ORDER
			} else {
				nextEvent = elevatorStatus.NO_EVENT
			}
			if prevEvent != nextEvent {
				fmt.Println("Event: ", nextEvent)
				if nextEvent == elevatorStatus.FLOOR_REACHED {
					driver.Set_floor_indicator(driver.Get_floor_sensor_signal())
				}
				eventChan <- nextEvent
				prevEvent = nextEvent
			}
		}
	}
}

func getNextDirection(e *elevatorStatus.Elevator) {
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	if e.CurrentFloor == 0 || e.CurrentFloor == 3 {
		e.Direction = driver.MDIR_STOP
	}
	if e.Direction != driver.MDIR_DOWN && e.CurrentFloor != 3 {
		if orderHandling.CheckUpOrdersAbove(*e) == true {
			e.Direction = driver.MDIR_UP
		} else {
			if orderHandling.CheckDownOrdersAbove(*e) == true {
				e.Direction = driver.MDIR_UP
			} else if e.Direction != driver.MDIR_STOP {
				e.Direction = driver.MDIR_STOP
			}
		}
	}
	if e.Direction != driver.MDIR_UP && e.CurrentFloor != 0 {
		if orderHandling.CheckDownOrdersBelow(*e) == true {
			e.Direction = driver.MDIR_DOWN
		} else {
			if orderHandling.CheckUpOrdersBelow(*e) == true {
				e.Direction = driver.MDIR_DOWN
			} else if e.Direction != driver.MDIR_STOP {
				e.Direction = driver.MDIR_STOP
			}
		}
	}
}

func StartUp(elevObject chan elevatorStatus.Elevator) {
	var e elevatorStatus.Elevator
	fmt.Println("Ready to start when network connection is established")
	for {
		if network.GetIpAddress() != "::1" {
			fmt.Println("Have established a network connection")
			break
		}
	}
	for driver.Get_floor_sensor_signal() == -1 {
		driver.Set_motor_speed(driver.MDIR_DOWN)
	}
	for driver.Get_floor_sensor_signal() != -1 {
		driver.Set_floor_indicator(driver.Get_floor_sensor_signal())
		driver.Set_motor_speed(driver.MDIR_STOP)
		break
	}
	for floor := 0; floor < driver.NUM_FLOORS; floor++ {
		for button := 0; button < driver.NUM_BUTTONS; button++ {
			if (floor == 0 && button == 1) || (floor == 3 && button == 0) {
			} else {
				driver.Set_button_lamp(button, floor, 0)
			}
		}
	}
	e.Direction = driver.MDIR_STOP
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	e.PreviousFloor = driver.Get_floor_sensor_signal()
	e.State = elevatorStatus.IDLE
	e.IP = network.GetIpAddress()
	//orderHandling.WriteInternalsToFile(e.OrderMatrix)
	e.OrderMatrix = orderHandling.ReadInternalsToFile()
	fmt.Println("StartUp values: ", e)
	go changeElev(e, elevObject)
}

func changeElev(e elevatorStatus.Elevator, elevObject chan elevatorStatus.Elevator) {
	elevObject <- e
}

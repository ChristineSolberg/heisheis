package elevatorControl

import(
	
	"./driver"
	"fmt"
	"time"
	"./orderHandling"
	"./elevatorStatus"
	"../network"
)



func UpdateFSM(e *elevatorStatus.Elevator, inToFSM chan elevatorStatus.Elevator, outOfFSM chan bool, DelOrder chan [4]int){
	eventChan := make(chan elevatorStatus.Event,100)
	//DoorTimeout := make(chan time.Time)
	//var DoorTimeout <-chan time.Time

	doorTimer := time.NewTimer(time.Millisecond)
	DoorTimeout := doorTimer.C
	fmt.Println("door timeout ",DoorTimeout)

	go getNextEvent(e, DoorTimeout, eventChan)
	for{	
		select{
		case newOrder := <-inToFSM:
			fmt.Println("newOrder: ", newOrder)
			//fmt.Println(e)
			//fmt.Println(e.OrderMatrix)
			e.OrderMatrix = orderHandling.UpdateOrderMatrix(newOrder.OrderMatrix,e.OrderMatrix)
			fmt.Println("newOrderMatrix: ", e.OrderMatrix)
		case event := <- eventChan:
			fmt.Println("Direction: ", e.Dir)
			switch(e.State){
			case elevatorStatus.IDLE:
				fmt.Println("State: IDLE")
				*e = updateFSM_IDLE(event, e, DelOrder, doorTimer)
				break
			case elevatorStatus.GO_UP:
				fmt.Println("State: GO UP")
				*e = updateFSM_GO_UP(event, e, DelOrder, doorTimer)
				break
			case elevatorStatus.GO_DOWN:
				fmt.Println("State: GO DOWN")
				*e = updateFSM_GO_DOWN(event, e, DelOrder, doorTimer)
				break
			case elevatorStatus.DOOR_OPEN:
				//time.Sleep(time.Second*3)
				fmt.Println("State: DOOR OPEN")
				

				*e = updateFSM_DOOR_OPEN(event, e, DelOrder, doorTimer)

				break
			default:
				fmt.Println("Error: No valid state in UpdateFSM")
			
			}
			outOfFSM <- true

			
		}
		
		time.Sleep(time.Millisecond * 100)
		
		
	}
}
// Nå vil vi legge ut en "ny" e på outOfFSM selv om Staten ikke nødvendigvis er oppdatert.
// Blir dette for ofte? Burde vi heller sjekke om staten er forskjellig fra forrige gang, og kun oppdatere da?


func updateFSM_IDLE(event elevatorStatus.Event, e *elevatorStatus.Elevator, DelOrder chan [4]int, DoorTimeout *time.Timer)elevatorStatus.Elevator{
	switch(event){
	case elevatorStatus.NEXT_ORDER:
		e.Dir = GetNextDirection(*e)
		driver.Set_motor_speed(e.Dir)
		if (e.Dir == driver.MDIR_UP){
			e.State = elevatorStatus.GO_UP
			fmt.Println("Ny state: GO_UP")
		} else if(e.Dir == driver.MDIR_DOWN){
			e.State = elevatorStatus.GO_DOWN
			fmt.Println("Ny state: GO_DOWN")
		} else{
			e.State = elevatorStatus.IDLE
		}
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		fmt.Println("UpdateFSM_IDLE: new order at current")
		DoorTimeout.Reset(time.Second*3)
		//DoorTimeout = time.Tick(time.Second * 3)
		orderHandling.DeleteCompletedOrders(e, DelOrder)
		e.State = elevatorStatus.DOOR_OPEN
		break
	case elevatorStatus.TIMER_OUT:
		fmt.Println("UpdateFSM_IDLE: timer out")
	case elevatorStatus.NO_EVENT:
		fmt.Println("UpdateFSM_IDLE: no event")
		//fmt.Println("Length of queue", orderHandling.LengthOfQueue(*e))
		//do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_IDLE!")
	}
	return *e
}



func updateFSM_GO_UP(event elevatorStatus.Event,e *elevatorStatus.Elevator,DelOrder chan [4]int, DoorTimeout *time.Timer)elevatorStatus.Elevator{
	fmt.Println("inne i updateFSM_GO_UP")
	switch (event){
	case elevatorStatus.FLOOR_REACHED:
		if (orderHandling.ShouldStop(*e) == 1){
			driver.Set_motor_speed(driver.MDIR_STOP)
			//Start timer, og legg true på doorTimeout når det har gått 3 sek.
			DoorTimeout.Reset(time.Second*3)
			//DoorTimeout = time.Tick(time.Second * 3)
			fmt.Println("før delete i FLOOR_REACHED")
			orderHandling.DeleteCompletedOrders(e, DelOrder)
			fmt.Println("etter delete i FLOOR_REACHED")
			e.PreviousFloor = driver.Get_floor_sensor_signal()
			e.Dir = GetNextDirection(*e)
			e.State = elevatorStatus.DOOR_OPEN
		} else{
			e.State = elevatorStatus.GO_UP
		}
		break
	case elevatorStatus.NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_GO_UP!")
	}
	return *e
}

func updateFSM_GO_DOWN(event elevatorStatus.Event, e *elevatorStatus.Elevator, DelOrder chan [4]int, DoorTimeout *time.Timer)elevatorStatus.Elevator{
	switch (event){
	case elevatorStatus.FLOOR_REACHED:
		fmt.Println("stop: ", orderHandling.ShouldStop(*e))
		if (orderHandling.ShouldStop(*e) == 1){
			driver.Set_motor_speed(driver.MDIR_STOP)
			e.State = elevatorStatus.DOOR_OPEN
			//Start timer, og legg true på doorTimeout når det har gått 3 sek.
			DoorTimeout.Reset(time.Second*3)
			fmt.Println("door timeout ",DoorTimeout)
			orderHandling.DeleteCompletedOrders(e, DelOrder)
			e.PreviousFloor = driver.Get_floor_sensor_signal()
			e.Dir = GetNextDirection(*e)
		} else{
			e.State = elevatorStatus.GO_DOWN
		}
		break
	case elevatorStatus.NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_GO_DOWN!")
	}
	return *e
}

func updateFSM_DOOR_OPEN(event elevatorStatus.Event, e *elevatorStatus.Elevator, DelOrder chan [4]int, DoorTimeout *time.Timer)elevatorStatus.Elevator{
	switch(event){
	case elevatorStatus.TIMER_OUT:
		driver.Set_door_open_lamp(0)
		e.State = elevatorStatus.IDLE
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		DoorTimeout.Reset(time.Second*3)
		//DoorTimeout = time.Tick(time.Second * 3)
		e.State = elevatorStatus.DOOR_OPEN
		orderHandling.DeleteCompletedOrders(e, DelOrder)
		break
	case elevatorStatus.NO_EVENT:
		fmt.Println("Length of queue", orderHandling.LengthOfQueue(*e))
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_DOOR_OPEN!")
	}
	return *e
}


func getNextEvent(e *elevatorStatus.Elevator, DoorTimeout <-chan time.Time, eventChan chan elevatorStatus.Event){
	fmt.Println("CurrentFloor: ", e.CurrentFloor)
	
	// får vi første no event nå? 
	var nextEvent elevatorStatus.Event
	var prevEvent elevatorStatus.Event

	for{
		eCopy := *e
		//fmt.Println("door timeout ",DoorTimeout)
		eCopy.CurrentFloor = driver.Get_floor_sensor_signal()
		select{
		case <-DoorTimeout:
			nextEvent = elevatorStatus.TIMER_OUT
			fmt.Println("Event: TIMER_OUT")

			if prevEvent != nextEvent{
				fmt.Println("Event: ", nextEvent)
				eventChan <-nextEvent
				prevEvent = nextEvent
			}
				//Dette vil også legges på channelen, ja?
		default:
			//fmt.Println("NewOrderAtCurrentFloor: ", orderHandling.NewOrderAtCurrentFloor(e))
			//fmt.Println("Queue: ", e.OrderMatrix)
			if (eCopy.CurrentFloor != -1) && (eCopy.PreviousFloor != eCopy.CurrentFloor){
				nextEvent = elevatorStatus.FLOOR_REACHED
				//fmt.Println("Event: FLOOR_REACHED")
			} else if ((eCopy.State == elevatorStatus.IDLE)||(eCopy.State == elevatorStatus.DOOR_OPEN)) && orderHandling.NewOrderAtCurrentFloor(eCopy) == 1{
				nextEvent = elevatorStatus.NEW_ORDER_AT_CURRENT
				//fmt.Println("Event: NEW_ORDER_AT_CURRENT")
			} else if (orderHandling.LengthOfQueue(eCopy) > 0 && eCopy.State == elevatorStatus.IDLE){  //&& (e.State == elevatorStatus.IDLE) || event == elevatorStatus.TIMER_OUT){
				//fmt.Println("Length of queue: ", orderHandling.LengthOfQueue(*e))
				nextEvent = elevatorStatus.NEXT_ORDER
				//fmt.Println("Event: NEXT_ORDER")
			} else{
				nextEvent = elevatorStatus.NO_EVENT
				//fmt.Println("Event = NO_EVENT")
			}
			
			if prevEvent != nextEvent{
				fmt.Println("Event: ", nextEvent)
				eventChan <-nextEvent
				prevEvent = nextEvent
			}
		}
	}
}




func GetNextDirection(e elevatorStatus.Elevator)driver.MotorDirection{
	e.CurrentFloor = driver.Get_floor_sensor_signal()

	if e.CurrentFloor == 0 || e.CurrentFloor == 3{
		e.Dir = driver.MDIR_STOP
	}
	fmt.Println("Retning etter første if")
	if(e.Dir != driver.MDIR_DOWN && e.CurrentFloor != 3){
		fmt.Println("Inne i getnextdirOPP")
		if(orderHandling.CheckUpOrdersAbove(e) == 1){
			e.Dir = driver.MDIR_UP
		} else {
			if(orderHandling.CheckDownOrdersAbove(e) == 1){
			e.Dir = driver.MDIR_UP
			} else if (e.Dir != driver.MDIR_STOP){
			e.Dir = driver.MDIR_STOP
			}
		}
	}

	if(e.Dir != driver.MDIR_UP && e.CurrentFloor != 0){
		fmt.Println("Inne i getnextdirNED")
		if orderHandling.CheckDownOrdersBelow(e) == 1{
			e.Dir = driver.MDIR_DOWN
		} else{
			if orderHandling.CheckUpOrdersBelow(e) == 1{
				e.Dir = driver.MDIR_DOWN
			} else if e.Dir != driver.MDIR_STOP{
				e.Dir = driver.MDIR_STOP
			}
		}
	}
	fmt.Println("Retning etter andre if")
	return e.Dir
}



func StartUp(e elevatorStatus.Elevator)elevatorStatus.Elevator{
	fmt.Println("Floor: ", driver.Get_floor_sensor_signal())
	for (driver.Get_floor_sensor_signal() == -1){
		driver.Set_motor_speed(driver.MDIR_DOWN)
	}
	for (driver.Get_floor_sensor_signal() != -1){
		driver.Set_motor_speed(driver.MDIR_STOP)
		break
	}

	for floor := 0; floor <driver.NUM_FLOORS; floor++{
		for button := 0; button < driver.NUM_BUTTONS; button++{
			if(floor == 0 && button == 1) || (floor == 3 && button == 0){
			} else{
				driver.Set_button_lamp(button,floor, 0)
			}
		}
	}

	//e.Dir = driver.MDIR_STOP
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	e.PreviousFloor = driver.Get_floor_sensor_signal()
	e.State = elevatorStatus.IDLE
	e.IP = network.GetIpAddress()
	return e
}
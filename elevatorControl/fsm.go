package elevatorControl

import(
	
	"./driver"
	"fmt"
	"time"
	"./orderHandling"
	"./elevatorStatus"
)


func UpdateFSM(e elevatorStatus.Elevator, inToFSM chan elevatorStatus.Elevator, outOfFSM chan elevatorStatus.Elevator, DelOrder chan [4]int){
	time.Sleep(time.Millisecond * 100)
	fmt.Println("Inne i updateFSM")
	e = <- inToFSM
	var DoorTimeout <-chan time.Time
	event := getNextEvent(e, DoorTimeout)
	
	
	fmt.Println("Direction: ", e.Dir)
	switch(e.State){
	case elevatorStatus.IDLE:
		fmt.Println("State: IDLE")
		e = updateFSM_IDLE(event, e, DelOrder, DoorTimeout)
		break
	case elevatorStatus.GO_UP:
		fmt.Println("State: GO UP")
		e = updateFSM_GO_UP(event, e, DelOrder, DoorTimeout)
		break
	case elevatorStatus.GO_DOWN:
		fmt.Println("State: GO DOWN")
		e = updateFSM_GO_DOWN(event, e, DelOrder, DoorTimeout)
		break
	case elevatorStatus.DOOR_OPEN:
		fmt.Println("State: DOOR OPEN")
		e = updateFSM_DOOR_OPEN(event, e, DelOrder, DoorTimeout)
		break
	default:
		fmt.Println("Error: No valid state in UpdateFSM")
	}
	outOfFSM <- e
}
// Nå vil vi legge ut en "ny" e på outOfFSM selv om Staten ikke nødvendigvis er oppdatert.
// Blir dette for ofte? Burde vi heller sjekke om staten er forskjellig fra forrige gang, og kun oppdatere da?


func updateFSM_IDLE(event elevatorStatus.Event, e elevatorStatus.Elevator, DelOrder chan [4]int, DoorTimeout <-chan time.Time)elevatorStatus.Elevator{
	switch(event){
	case elevatorStatus.NEXT_ORDER:
		e.Dir = GetNextDirection(e)
		driver.Set_motor_speed(e.Dir)
		if (e.Dir == driver.MDIR_UP){
			e.State = elevatorStatus.GO_UP
		} else if(e.Dir == driver.MDIR_DOWN){
			e.State = elevatorStatus.GO_DOWN
		} else{
			e.State = elevatorStatus.IDLE
		}
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		fmt.Println("UpdateFSM_IDLE: new order at current")
		DoorTimeout = time.Tick(time.Second * 3)
		orderHandling.DeleteCompletedOrders(&e, DelOrder)
		e.State = elevatorStatus.DOOR_OPEN
		break
	case elevatorStatus.TIMER_OUT:
		fmt.Println("UpdateFSM_IDLE: timer out")
	case elevatorStatus.NO_EVENT:
		fmt.Println("UpdateFSM_IDLE: no event")
		fmt.Println("Length of queue", orderHandling.LengthOfQueue(e))
		//do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_IDLE!")
	}
	return e
}



func updateFSM_GO_UP(event elevatorStatus.Event,e elevatorStatus.Elevator,DelOrder chan [4]int, DoorTimeout <-chan time.Time)elevatorStatus.Elevator{
	switch (event){
	case elevatorStatus.FLOOR_REACHED:
		if (orderHandling.ShouldStop(e) == 1){
			driver.Set_motor_speed(driver.MDIR_STOP)
			//Start timer, og legg true på doorTimeout når det har gått 3 sek.
			DoorTimeout = time.Tick(time.Second * 3)
			fmt.Println("før delete i FLOOR_REACHED")
			orderHandling.DeleteCompletedOrders(&e, DelOrder)
			fmt.Println("etter delete i FLOOR_REACHED")
			e.PreviousFloor = driver.Get_floor_sensor_signal()
			e.Dir = GetNextDirection(e)
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
	return e
}

func updateFSM_GO_DOWN(event elevatorStatus.Event, e elevatorStatus.Elevator, DelOrder chan [4]int, DoorTimeout <-chan time.Time)elevatorStatus.Elevator{
	switch (event){
	case elevatorStatus.FLOOR_REACHED:
		fmt.Println("stop: ", orderHandling.ShouldStop(e))
		if (orderHandling.ShouldStop(e) == 1){
			driver.Set_motor_speed(driver.MDIR_STOP)
			e.State = elevatorStatus.DOOR_OPEN
			//Start timer, og legg true på doorTimeout når det har gått 3 sek.
			DoorTimeout = time.Tick(time.Second * 3)
			orderHandling.DeleteCompletedOrders(&e, DelOrder)
			e.PreviousFloor = driver.Get_floor_sensor_signal()
			e.Dir = GetNextDirection(e)
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
	return e
}

func updateFSM_DOOR_OPEN(event elevatorStatus.Event, e elevatorStatus.Elevator, DelOrder chan [4]int, DoorTimeout <-chan time.Time)elevatorStatus.Elevator{
	switch(event){
	case elevatorStatus.TIMER_OUT:
		driver.Set_door_open_lamp(0) //Er dette greit? Sjekk ut notatan :) å sette lys her altså
		e.State = elevatorStatus.IDLE
		/*fmt.Println("Length of queue", orderHandling.LengthOfQueue(e))
		e.Dir = GetNextDirection(e)
		fmt.Println("Retning etter door_open: ", e.Dir)
		if (orderHandling.LengthOfQueue(e) == 0){
			e.State = elevatorStatus.IDLE
		} else if (e.Dir == driver.MDIR_UP){
			driver.Set_motor_speed(e.Dir)
			e.State = elevatorStatus.GO_UP
		} else if (e.Dir == driver.MDIR_DOWN){
			driver.Set_motor_speed(e.Dir)
			e.State = elevatorStatus.GO_DOWN
		}*/
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		DoorTimeout = time.Tick(time.Second * 3)
		e.State = elevatorStatus.DOOR_OPEN
		orderHandling.DeleteCompletedOrders(&e, DelOrder)
		break
	case elevatorStatus.NO_EVENT:
		fmt.Println("Length of queue", orderHandling.LengthOfQueue(e))
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_DOOR_OPEN!")
	}
	return e
}


func getNextEvent(e elevatorStatus.Elevator, DoorTimeout <-chan time.Time)elevatorStatus.Event{
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	var event elevatorStatus.Event

	select{
	case <-DoorTimeout:
		event = elevatorStatus.TIMER_OUT
		fmt.Println("Event: TIMER_OUT")
		return event
	default:
		fmt.Println("ingenting på channel")
	}
	//fmt.Println("NewOrderAtCurrentFloor: ", orderHandling.NewOrderAtCurrentFloor(e))
	if (e.CurrentFloor != -1) && (e.PreviousFloor != e.CurrentFloor){
		event = elevatorStatus.FLOOR_REACHED
		fmt.Println("Event: FLOOR_REACHED")
	} else if ((e.State == elevatorStatus.IDLE)||(e.State == elevatorStatus.DOOR_OPEN)) && orderHandling.NewOrderAtCurrentFloor(e) == 1{
		event = elevatorStatus.NEW_ORDER_AT_CURRENT
		fmt.Println("Event: NEW_ORDER_AT_CURRENT")
	} else if (orderHandling.LengthOfQueue(e) > 0 && e.State == elevatorStatus.IDLE){  //&& (e.State == elevatorStatus.IDLE) || event == elevatorStatus.TIMER_OUT){
		fmt.Println("Length of queue: ", orderHandling.LengthOfQueue(e))
		event = elevatorStatus.NEXT_ORDER
		fmt.Println("Event: NEXT_ORDER")
	} else{
		event = elevatorStatus.NO_EVENT
		fmt.Println("Event = NO_EVENT")
	}
	return event
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

	e.PreviousFloor = driver.Get_floor_sensor_signal()
	e.State = elevatorStatus.IDLE
	e.Dir = driver.MDIR_STOP
	return e
}
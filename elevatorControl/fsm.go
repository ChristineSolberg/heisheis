package elevatorControl

import(
	
	"./driver"
	"fmt"
	"time"
	"./orderHandling"
	"./elevatorStatus"
)

var doorTimeout = make(chan bool)
//var doorTimerReset = make(chan bool)


func UpdateFSM(e elevatorStatus.Elevator)elevatorStatus.Elevator{
	fmt.Println("Inne i updateFSM")
	event := getNextEvent(e)
	e.Dir = GetNextDirection(e)
	fmt.Println("Direction: ", e.Dir)
	switch(e.State){
	case elevatorStatus.IDLE:
		fmt.Println("State: IDLE")
		e = updateFSM_IDLE(event, e)
		break
	case elevatorStatus.GO_UP:
		fmt.Println("State: GO UP")
		e = updateFSM_GO_UP(event, e)
		break
	case elevatorStatus.GO_DOWN:
		fmt.Println("State: GO DOWN")
		e = updateFSM_GO_DOWN(event, e)
		break
	case elevatorStatus.DOOR_OPEN:
		fmt.Println("State: DOOR OPEN")
		e = updateFSM_DOOR_OPEN(event, e)
		break
	default:
		fmt.Println("Error: No valid state in UpdateFSM")
	}
	return e
}


func updateFSM_IDLE(event elevatorStatus.Event, e elevatorStatus.Elevator)elevatorStatus.Elevator{
	switch(event){
	case elevatorStatus.NEXT_ORDER:
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
		e.State = elevatorStatus.DOOR_OPEN
		time.Sleep(3*time.Second)
		orderHandling.DeleteCompletedOrders(e)
		break
	case elevatorStatus.NO_EVENT:
		fmt.Println("UpdateFSM_IDLE: no event")
		//do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_IDLE!")
	}
	return e
}



func updateFSM_GO_UP(event elevatorStatus.Event,e elevatorStatus.Elevator)elevatorStatus.Elevator{
	switch (event){
	case elevatorStatus.FLOOR_REACHED:
		if (orderHandling.ShouldStop(e) == 1){
			//e.Dir = MDIR_STOP
			driver.Set_motor_speed(driver.MDIR_STOP)
			e.State = elevatorStatus.DOOR_OPEN
			time.Sleep(3*time.Second)
			orderHandling.DeleteCompletedOrders(e)
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

func updateFSM_GO_DOWN(event elevatorStatus.Event, e elevatorStatus.Elevator)elevatorStatus.Elevator{
	switch (event){
	case elevatorStatus.FLOOR_REACHED:
		if (orderHandling.ShouldStop(e) == 1){
			//e.Dir = driver.MDIR_STOP
			driver.Set_motor_speed(driver.MDIR_STOP)
			e.State = elevatorStatus.DOOR_OPEN
			time.Sleep(3*time.Second)
			orderHandling.DeleteCompletedOrders(e)
		} else{
			e.State = elevatorStatus.GO_UP
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

func updateFSM_DOOR_OPEN(event elevatorStatus.Event, e elevatorStatus.Elevator)elevatorStatus.Elevator{
	switch(event){
	case elevatorStatus.TIMER_OUT:
		driver.Set_door_open_lamp(0) //Er dette greit? Sjekk ut notatan :) å sette lys her altså
		if (orderHandling.LengthOfQueue(e) == 0){
			e.State = elevatorStatus.IDLE
		} else if (e.Dir == 1){
			driver.Set_motor_speed(e.Dir)
			e.State = elevatorStatus.GO_UP
		} else if (e.Dir == -1){
			driver.Set_motor_speed(e.Dir)
			e.State = elevatorStatus.GO_DOWN
		}
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		//start timer her
		e.State = elevatorStatus.DOOR_OPEN
		break
	case elevatorStatus.NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_DOOR_OPEN!")
	}
	return e
}


func getNextEvent(e elevatorStatus.Elevator)elevatorStatus.Event{
	floor := driver.Get_floor_sensor_signal()
	var event elevatorStatus.Event

	fmt.Println("Over select")
	select{
	case <- doorTimeout:
		event = elevatorStatus.TIMER_OUT
		fmt.Println("Event: TIMER_OUT")
	default:
		fmt.Println("ingenting på channel")
	}
	fmt.Println("under select")
	
	if (floor != -1) && (e.PreviousFloor != floor){
		event = elevatorStatus.FLOOR_REACHED
		fmt.Println("Event: FLOOR_REACHED")
	} else if ((e.State == elevatorStatus.IDLE)||(e.State == elevatorStatus.DOOR_OPEN)) && orderHandling.NewOrderAtCurrentFloor(e) == 1{
		event = elevatorStatus.NEW_ORDER_AT_CURRENT
		fmt.Println("Event: NEW_ORDER_AT_CURRENT")
	} else if (orderHandling.LengthOfQueue(e) > 0) && (e.State == elevatorStatus.IDLE){
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
	if(e.Dir != driver.MDIR_DOWN){
		fmt.Println("Inne i getnextdirOPP")
		if(orderHandling.CheckUpOrdersAbove(e) == 1){
			e.Dir = driver.MDIR_UP
			return e.Dir
		} else {
			if(orderHandling.CheckDownOrdersAbove(e) == -1){
			e.Dir = driver.MDIR_UP
			return e.Dir
			} else if (e.Dir != driver.MDIR_STOP){
			e.Dir = driver.MDIR_STOP
			return e.Dir
			}
		}
	}

	if(e.Dir != driver.MDIR_UP){
		fmt.Println("Inne i getnextdirNED")
		if orderHandling.CheckDownOrdersBelow(e) == 1{
			e.Dir = driver.MDIR_DOWN
			return e.Dir
		} else{
			if orderHandling.CheckUpOrdersBelow(e) == 1{
				e.Dir = driver.MDIR_DOWN
				return e.Dir
			} else if e.Dir != driver.MDIR_STOP{
				e.Dir = driver.MDIR_STOP
				return e.Dir
			}
		}
	}

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
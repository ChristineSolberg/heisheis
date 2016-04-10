package elevatorControl

import(
	"fmt"
	"time"
	"./driver"
	"./orderHandling"
	"./elevatorStatus"
	"../network"
)



func UpdateFSM(newOrderToFSM chan elevatorStatus.Elevator, newStateUpdate chan bool, DelOrder chan [4]int, elevChan chan elevatorStatus.Elevator, powerChan chan bool, abortElev chan bool){
	doorTimer := time.NewTimer(time.Nanosecond)
	DoorTimeout := doorTimer.C
	powerTimer := time.NewTimer(time.Nanosecond)
	PowerTimeout := powerTimer.C
	//fmt.Println("door timeout ",DoorTimeout)
	//DoorTimeout := make(<-chan time.Time) //liten bokstav??
	//var powerTimer time.Timer

	eventChan := make(chan elevatorStatus.Event,100)
	go getNextEvent(elevChan, /*&DoorTimeout*/eventChan, powerChan, PowerTimeout, DoorTimeout)

	var turnOff bool
	for{
		select{
		case newOrder := <-newOrderToFSM:
			fmt.Println("newOrder: ", newOrder)
			orderHandling.UpdateOrderMatrix(newOrder.OrderMatrix,elevChan)
		case turnOff = <- abortElev:
			//if turnOff == true{
			fmt.Println("turnoff til turned")
			//break
		

		case event := <- eventChan:
			//fmt.Println("Direction: ", e.Dir)
			state := getElevatorState(elevChan)
			switch(state){
			case elevatorStatus.IDLE:
				fmt.Println("State: IDLE")
				updateFSM_IDLE(event, elevChan, DelOrder /*&DoorTimeout*/, powerTimer,doorTimer)
				break
			case elevatorStatus.GO_UP:
				fmt.Println("State: GO UP")
				updateFSM_GO_UP(event, elevChan, DelOrder /*&DoorTimeout*/, powerChan, powerTimer,doorTimer)
				break
			case elevatorStatus.GO_DOWN:
				fmt.Println("State: GO DOWN")
				updateFSM_GO_DOWN(event, elevChan, DelOrder/*&DoorTimeout*/, powerChan, powerTimer,doorTimer)
				break
			case elevatorStatus.DOOR_OPEN:
				fmt.Println("State: DOOR OPEN")
				updateFSM_DOOR_OPEN(event, elevChan, DelOrder /*&DoorTimeout*/,doorTimer)
				break
			default:
				fmt.Println("Error: No valid state in UpdateFSM")
			
			}
			newStateUpdate <- true
		}
		if (turnOff == true && network.GetIpAddress() != "::1"){
			fmt.Println("Slår av heis")
			break
		}
	}
}

func getElevatorState(elevChan chan elevatorStatus.Elevator)elevatorStatus.ElevState{
	e := <- elevChan
	state := e.State
	elevChan <- e
	return state
}

func updateFSM_IDLE(event elevatorStatus.Event, elevChan chan elevatorStatus.Elevator, DelOrder chan [4]int, powerTimer *time.Timer, doorTimer *time.Timer){// DoorTimeout *(<-chan time.Time)){
	e := <- elevChan
	fmt.Println("elevator object fsmUpdate idle", e)
	switch(event){
	case elevatorStatus.NEXT_ORDER:
		GetNextDirection(&e)
		driver.Set_motor_speed(e.Dir)
		if (e.Dir == driver.MDIR_UP){
			e.State = elevatorStatus.GO_UP
			powerTimer.Reset(time.Second * 20)
			fmt.Println("Ny state: GO_UP")
		} else if(e.Dir == driver.MDIR_DOWN){
			e.State = elevatorStatus.GO_DOWN
			powerTimer.Reset(time.Second * 20)
			fmt.Println("Ny state: GO_DOWN")
		} else{
			e.State = elevatorStatus.IDLE
		}
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		fmt.Println("UpdateFSM_IDLE: new order at current")
		//DoorTimeout.Reset(time.Second*3)
		driver.Set_door_open_lamp(1)
		//*DoorTimeout = time.Tick(time.Second * 3)
		doorTimer.Reset(time.Second*3)
		orderHandling.DeleteCompletedOrders(&e, DelOrder)
		e.State = elevatorStatus.DOOR_OPEN
		break
	case elevatorStatus.TIMER_OUT:
		fmt.Println("UpdateFSM_IDLE: timer out")
		driver.Set_door_open_lamp(0)
	case elevatorStatus.NO_EVENT:
		fmt.Println("UpdateFSM_IDLE: no event")
		//fmt.Println("Length of queue", orderHandling.LengthOfQueue(*e))
		//do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_IDLE!")
	}
	elevChan <- e
}



func updateFSM_GO_UP(event elevatorStatus.Event,elevChan chan elevatorStatus.Elevator,DelOrder chan [4]int, powerChan chan bool, powerTimer *time.Timer, DoorTimeout *time.Timer){// DoorTimeout *(<-chan time.Time)){
	//fmt.Println("inne i updateFSM_GO_UP")
	e := <- elevChan
	switch (event){
	case elevatorStatus.FLOOR_REACHED:
		if (orderHandling.ShouldStop(e) == 1){
			driver.Set_motor_speed(driver.MDIR_STOP)
			driver.Set_door_open_lamp(1)
			//Start timer, og legg true på doorTimeout når det har gått 3 sek.
			DoorTimeout.Reset(time.Second*3)
			//*DoorTimeout = time.Tick(time.Second * 3)
			powerTimer.Stop()
			fmt.Println("før delete i FLOOR_REACHED")
			orderHandling.DeleteCompletedOrders(&e, DelOrder)
			fmt.Println("etter delete i FLOOR_REACHED")
			e.PreviousFloor = driver.Get_floor_sensor_signal()
			GetNextDirection(&e)
			e.State = elevatorStatus.DOOR_OPEN
		} else{
			e.State = elevatorStatus.GO_UP
		}
		break
	case elevatorStatus.POWER_OFF:
		fmt.Println("Vi tror power er off :o")
		powerChan <- true
		
	case elevatorStatus.NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_GO_UP!")
	}
	elevChan <- e
}


func updateFSM_GO_DOWN(event elevatorStatus.Event, elevChan chan elevatorStatus.Elevator, DelOrder chan [4]int, powerChan chan bool, powerTimer *time.Timer, DoorTimeout *time.Timer){//DoorTimeout *(<-chan time.Time)){
	e := <- elevChan
	switch (event){
	case elevatorStatus.FLOOR_REACHED:
		fmt.Println("stop: ", orderHandling.ShouldStop(e))
		if (orderHandling.ShouldStop(e) == 1){
			driver.Set_motor_speed(driver.MDIR_STOP)
			driver.Set_door_open_lamp(1)
			e.State = elevatorStatus.DOOR_OPEN
			DoorTimeout.Reset(time.Second*3)
			powerTimer.Stop()
			//*DoorTimeout = time.Tick(time.Second * 3)
			orderHandling.DeleteCompletedOrders(&e, DelOrder)
			e.PreviousFloor = driver.Get_floor_sensor_signal()
			GetNextDirection(&e)
		} else{
			e.State = elevatorStatus.GO_DOWN
		}
		break
	case elevatorStatus.POWER_OFF:
		powerChan <- true
	case elevatorStatus.NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_GO_DOWN!")
	}
	elevChan <- e
}

func updateFSM_DOOR_OPEN(event elevatorStatus.Event, elevChan chan elevatorStatus.Elevator, DelOrder chan [4]int, DoorTimeout *time.Timer){//DoorTimeout *(<-chan time.Time)){
	e := <-elevChan
	switch(event){
	case elevatorStatus.TIMER_OUT:
		driver.Set_door_open_lamp(0)
		e.State = elevatorStatus.IDLE
		break
	case elevatorStatus.NEW_ORDER_AT_CURRENT:
		DoorTimeout.Reset(time.Second*3)
		driver.Set_door_open_lamp(1)
		//*DoorTimeout = time.Tick(time.Second * 3)
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
	elevChan <- e
}


/*DoorTimeout *(<-chan time.Time)*/
func getNextEvent(elevChan chan elevatorStatus.Elevator, eventChan chan elevatorStatus.Event, powerChan chan bool, PowerTimeout <- chan time.Time, DoorTimeout <- chan time.Time){
	
	
	// får vi første no event nå? 
	var nextEvent elevatorStatus.Event
	var prevEvent elevatorStatus.Event

	for{
		e := <-elevChan
		
		go changeElev(e,elevChan)
		//fmt.Println("checking next event, this is my elevator etter", e)

		//eCopy := *e
		//fmt.Println("door timeout ",DoorTimeout)
		currentFloor := driver.Get_floor_sensor_signal()
		//fmt.Println("my floor ", currentFloor)
		select{
		case <-DoorTimeout:
			nextEvent = elevatorStatus.TIMER_OUT
			fmt.Println("Event: TIMER_OUT")
			//elevChan <- e

			if prevEvent != nextEvent{
				fmt.Println("Event: ", nextEvent)
				eventChan <-nextEvent
				prevEvent = nextEvent
			}
		case <-PowerTimeout:
			fmt.Println("Nå får vi snart event POWER_OFF")
			nextEvent = elevatorStatus.POWER_OFF 
			eventChan <-nextEvent
			//ikke ta imot bestillinger her!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		default:
			if (currentFloor != -1) && (e.PreviousFloor != currentFloor){
				nextEvent = elevatorStatus.FLOOR_REACHED
				//fmt.Println("Event: FLOOR_REACHED")
			} else if ((e.State == elevatorStatus.IDLE)||(e.State == elevatorStatus.DOOR_OPEN)) && orderHandling.NewOrderAtCurrentFloor(e) == 1{
				nextEvent = elevatorStatus.NEW_ORDER_AT_CURRENT
				//fmt.Println("Event: NEW_ORDER_AT_CURRENT")
			} else if (orderHandling.LengthOfQueue(e) > 0 && e.State == elevatorStatus.IDLE){  //&& (e.State == elevatorStatus.IDLE) || event == elevatorStatus.TIMER_OUT){
				//fmt.Println("Length of queue: ", orderHandling.LengthOfQueue(*e))
				nextEvent = elevatorStatus.NEXT_ORDER
				//fmt.Println("Event: NEXT_ORDER")
			} else{
				nextEvent = elevatorStatus.NO_EVENT
				//fmt.Println("Event = NO_EVENT")
			}
			
			if prevEvent != nextEvent{
				fmt.Println("Event: ", nextEvent)
				if nextEvent == elevatorStatus.FLOOR_REACHED{
					driver.Set_floor_indicator(driver.Get_floor_sensor_signal())
				}
				eventChan <-nextEvent
				prevEvent = nextEvent
				
			}
		}
		
	}
}




func GetNextDirection(e *elevatorStatus.Elevator){
	e.CurrentFloor = driver.Get_floor_sensor_signal()

	if e.CurrentFloor == 0 || e.CurrentFloor == 3{
		e.Dir = driver.MDIR_STOP
	}
	fmt.Println("Retning etter første if")
	if(e.Dir != driver.MDIR_DOWN && e.CurrentFloor != 3){
		fmt.Println("Inne i getnextdirOPP")
		if(orderHandling.CheckUpOrdersAbove(*e) == 1){
			e.Dir = driver.MDIR_UP
		} else {
			if(orderHandling.CheckDownOrdersAbove(*e) == 1){
			e.Dir = driver.MDIR_UP
			} else if (e.Dir != driver.MDIR_STOP){
			e.Dir = driver.MDIR_STOP
			}
		}
	}

	if(e.Dir != driver.MDIR_UP && e.CurrentFloor != 0){
		fmt.Println("Inne i getnextdirNED")
		if orderHandling.CheckDownOrdersBelow(*e) == 1{
			e.Dir = driver.MDIR_DOWN
		} else{
			if orderHandling.CheckUpOrdersBelow(*e) == 1{
				e.Dir = driver.MDIR_DOWN
			} else if e.Dir != driver.MDIR_STOP{
				e.Dir = driver.MDIR_STOP
			}
		}
	}
	fmt.Println("Retning etter andre if")
}



func StartUp(elevChan chan elevatorStatus.Elevator){
	fmt.Println("Floor: ", driver.Get_floor_sensor_signal()) //Kan fjernes

	var e elevatorStatus.Elevator
	fmt.Println("Ready to start when network connection is established")
	for {
		if network.GetIpAddress() != "::1"{
			fmt.Println("Have established a network connection")
			break
		}
	}

	//If the elevator is in between two floors, it will go down to the nearest floor below
	for (driver.Get_floor_sensor_signal() == -1){
		driver.Set_motor_speed(driver.MDIR_DOWN)
	}
	for (driver.Get_floor_sensor_signal() != -1){
		driver.Set_floor_indicator(driver.Get_floor_sensor_signal())
		driver.Set_motor_speed(driver.MDIR_STOP)
		break
	}

	//Making sure all lamps are turned off
	for floor := 0; floor <driver.NUM_FLOORS; floor++{
		for button := 0; button < driver.NUM_BUTTONS; button++{
			if(floor == 0 && button == 1) || (floor == 3 && button == 0){
			} else{
				driver.Set_button_lamp(button,floor, 0)
			}
		}
	}

	e.Dir = driver.MDIR_STOP
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	e.PreviousFloor = driver.Get_floor_sensor_signal()
	e.State = elevatorStatus.IDLE
	e.IP = network.GetIpAddress()
	//orderHandling.WriteInternals(e.OrderMatrix)
	e.OrderMatrix = orderHandling.ReadInternals()
	fmt.Println("StartUp values: ", e)
	go changeElev(e,elevChan)
	
}

func changeElev(e elevatorStatus.Elevator, elevChan chan elevatorStatus.Elevator){
	elevChan <- e 
}

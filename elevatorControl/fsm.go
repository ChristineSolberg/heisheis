package elevatorControl

import(
	
	"./driver"
	"fmt"
)

type State int
const (
	IDLE  			State 	= 0 
	GO_UP					= 1
	GO_DOWN					= 2
	DOOR_OPEN				= 3
)

type Event int
const (
	FLOOR_REACHED 		Event 	= 0
	TIMER_OUT					= 1
	NEXT_ORDER					= 2
	NEW_ORDER_AT_CURRENT		= 3
	NO_EVENT					= 4
	
)

var previousFloor int = -1
var OrderMatrix [3][4]int // is it stupid to have global variables? hmm?


func UpdateFSM(state State, dir MotorDirection){
	event := getNextEvent(FLOOR_REACHED, IDLE, 1, 1)
	switch(state){
	case IDLE:
		fmt.Println("State: IDLE")
		updateFSM_IDLE(event, dir,state)
		break
	case GO_UP:
		fmt.Println("State: GO UP")
		updateFSM_GO_UP(event, dir,state)
		break
	case GO_DOWN:
		fmt.Println("State: GO DOWN")
		updateFSM_GO_DOWN(event, dir,state)
		break
	case DOOR_OPEN:
		fmt.Println("State: DOOR OPEN")
		updateFSM_DOOR_OPEN(event, dir, state)
		break
	default:
		fmt.Println("Error: No valid state in UpdateFSM")
	}
}


func updateFSM_IDLE(event Event, dir MotorDirection, state State){
	switch(event){
	case NEXT_ORDER:
		Set_motor_speed(dir)
		if (dir == 1){
			state = GO_UP
		} else if(dir == -1){
			state = GO_DOWN
		} else{
			state = IDLE
		}
		break
	case NEW_ORDER_AT_CURRENT:
		// Start timer function here
		state = DOOR_OPEN
		break
	case NO_EVENT:
		//do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_IDLE!")
	}
}



func updateFSM_GO_UP(event Event, dir MotorDirection, state State){
	switch (event){
	case FLOOR_REACHED:
		if (shouldStop()){
			dir = 0
			Set_motor_speed(dir)
			//start timer her
			state = DOOR_OPEN
		} else{
			state = GO_UP
		}
		break
	case NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_GO_UP!")
	}
}

func updateFSM_GO_DOWN(event Event, dir MotorDirection, state State){
	switch (event){
	case FLOOR_REACHED:
		if (shouldStop()){
			dir = 0
			Set_motor_speed(dir)
			//start timer her
			state = DOOR_OPEN
		} else{
			state = GO_UP
		}
		break
	case NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_GO_DOWN!")
	}
}

func updateFSM_DOOR_OPEN(event Event, dir MotorDirection, state State){
	switch(event){
	case TIMER_OUT:
		Set_door_open_lamp(0) //Er dette greit? Sjekk ut notatan :) å sette lys her altså
		if (lengthOfQueue() == 0){
			state = IDLE
		} else if (dir == 1){
			Set_motor_speed(dir)
			state = GO_UP
		} else if (dir == -1){
			Set_motor_speed(dir)
			state = GO_DOWN
		}
		break
	case NEW_ORDER_AT_CURRENT:
		//start timer her
		state = DOOR_OPEN
		break
	case NO_EVENT:
		//Do nothing
		break
	default:
		fmt.Println("Error: Error in updateFSM_DOOR_OPEN!")
	}
}


func getNextEvent(event Event, state State, dir MotorDirection)Event{
	floor := Get_floor_sensor_signal()

	if (floor != -1) && (previousFloor != floor){
		event = FLOOR_REACHED
		fmt.Println("Event: FLOOR_REACHED")
	} else if ((state == IDLE)||(state == DOOR_OPEN)) && newOrderAtCurrentFloor(dir){
		event = NEW_ORDER_AT_CURRENT
		fmt.Println("Event: NEW_ORDER_AT_CURRENT")
	} else if (lengthOfQueue() > 0) || (state == IDLE){
		event = NEXT_ORDER
		fmt.Println("Event: NEXT_ORDER")
	} else{
		event = NO_EVENT
		fmt.Println("Event = NO_EVENT")
	}

	return event
}

// func getNextOrder(dir MotorDirection, prevFloor int)int{
// 	if(dir != -1){
// 		if(checkOrdersAbove(prevFloor) == 1){
// 			dir = MDIR_UP
// 			return 1
// 		} else if (checkOrdersAbove(prevFloor) == -1){
// 			dir = MDIR_UP
// 			return 1
// 		} else if (dir != 0){
// 			dir = MDIR_STOP
// 			return 0
// 		}
// 	}

// 	// ikke ferdig

// }

func newOrderAtCurrentFloor(dir MotorDirection)int{
	floor := Get_floor_sensor_signal()
	if(floor==3){	
		if (OrderMatrix[floor][1] == 1){
			return 1;
		} else if (OrderMatrix[floor][2] == 1){
			return 1;
		}
	} else if(floor==0){	
		if (OrderMatrix[floor][0] == 1){
			return 1;
		} else if (OrderMatrix[floor][2] == 1){
			return 1;
		}
	} 

	if(dir!= -1){
		if (OrderMatrix[floor][0] == 1){
			return 1;
		} else if (OrderMatrix[floor][2] == 1){
			return 1;						
		}
	} else if(dir!= 1){
		if (OrderMatrix[floor][1] == 1){
			return 1;
		} else if (OrderMatrix[floor][2] == 1){
			return 1;						
		}
	}
	return 0;
}

func shouldStop(dir MotorDirection)int{
	//Lag denne senere - kanskje i queue
	floor := Get_floor_sensor_signal();
	
	if(floor>previousFloor){
		if (OrderMatrix[floor][2] == 1) || (OrderMatrix[floor][0] == 1)
			return 1;	
		if (!(checkUpOrdersAbove()||checkDownOrdersAbove()))
			return 1;				
	} else if(floor<previousFloor){
		if (OrderMatrix[floor][2] == 1) || (OrderMatrix[floor][1] == 1)
			return 1;
		if (!(checkUpOrdersBelow()||checkDownOrdersBelow()))
			return 1;
	}
	return 0;
}

// Implementer, -1 for ned knapp ovenfra, 1 for for opp knapp ovenfra
func checkUpOrdersAbove(){
	for (floor := previousFloor+1; floor < NUM_FLOORS; floor++){
		if(OrderMatrix[floor][0] == 1){
			return 1;
		}
		if(OrderMatrix[floor][2] == 1){
			return 1;
		}				
	}
	return 0
}

func checkDownOrdersAbove(){
	for(floor := previousFloor + 1; floor < NUM_FLOORS; floor++){
		if(OrderMatrix[floor][1] ==1)
			return 1;
		else
			return 0;
	}
	return 0;	
} 

func checkUpOrdersBelow(){
	for(floor := 0; floor < previousFloor; floor++){
		if(OrderMatrix[floor][0] ==1){
			return 1;
		}
		else
			return 0;
	}
	return 0;

}

func checkDownOrdersBelow(){

	for (floor := 0; floor < previousFloor; floor++){				
			if(OrderMatrix[floor][1] == 1){
				return 1;
			}
			else if(OrderMatrix[floor][2] == 1){
				return 1;
			}	
	}
	return 0;
}

func lengthOfQueue()int{
	length := 0
	for floor := 0; floor < NUM_FLOORS; floor++{
		for button := 0; button < NUM_BUTTONS; button++{
			if (button == 1 && floor == 0) || (button == 2 && floor == 3){
			}else{
				length += OrderMatrix[button][floor] //lag matrisen!
			}
		}
	} 
	return length
}

func StartUp(state State, dir MotorDirection){
	for (Get_floor_sensor_signal() == -1){
		Set_motor_speed(MDIR_DOWN)
	}
	for (Get_floor_sensor_signal() != -1){
		Set_motor_speed(MDIR_STOP)
		break
	}

	for floor := 0; floor < NUM_FLOORS; floor++{
		for button := 0; button < NUM_BUTTONS; button++{
			if(floor == 0 && button == 1) || (floor == 3 && button == 0){
			} else{
				Set_button_lamp(button,floor, 0)
			}
		}
	}

	previousFloor = Get_floor_sensor_signal()
	state = IDLE
	dir = MDIR_STOP
}
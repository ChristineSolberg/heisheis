package orderHandling

import(
	"fmt"
	"../driver"
	"../elevatorStatus"
)



// this is for one elevator, når n heiser, skal få ordre fra master og ikke les av knappetrykk

func AddOrderToQueue(e elevatorStatus.Elevator) elevatorStatus.Elevator{
	fmt.Println("inne i addorderqueue")
	for floor := 0; floor < driver.NUM_FLOORS; floor++{
		for button := driver.BUTTON_CALL_UP; button <= driver.BUTTON_COMMAND ; button++{
			if (floor == 0 && button == 1) || (floor == 3 && button == 0) || (floor < 0){
			} else if driver.Get_button_signal(button, floor) == 1{
				e.OrderMatrix[floor][button] = 1
				// sett knappelys
			}
		}
	}
	return e
}

// for n antall heiser.  ordren skal sendes til master i stedet for å legges inn i matrisen for n heiser
// func readButtons() int{
// 	result int = 0
// 	for floor := 0; floor < NUM_FLOORS; floor++{
// 		for button := 0; button < NUM_BUTTONS; button++{
// 			if (floor == 0 && button == 1) || (floor == 3 && button == 0) || (floor < 0){
// 			} else if Get_button_signal(button, floor) == 1{
// 				// sett knappelys
				// send ordre til master
// 				result = 1
// 			}
// 		}
// 	}
// 	return result
// }

func ShouldStop(e elevatorStatus.Elevator)int{
	//Lag denne senere - kanskje i queue
	floor := driver.Get_floor_sensor_signal();
	
	if(floor > e.PreviousFloor){
		if (e.OrderMatrix[floor][2] == 1) || (e.OrderMatrix[floor][0] == 1){
			return 1;	
		}
		if (!(CheckUpOrdersAbove(e) == 1 ||CheckDownOrdersAbove(e) == 1)){
			return 1;				
		}
	} else if(floor < e.PreviousFloor){
		if (e.OrderMatrix[floor][2] == 1) || (e.OrderMatrix[floor][1] == 1){
			return 1;
		}
		if (!(CheckUpOrdersBelow(e) == 1||CheckDownOrdersBelow(e) == 1)){
			return 1;
		}
	}
	return 0;
}

// Implementer, -1 for ned knapp ovenfra, 1 for for opp knapp ovenfra
func CheckUpOrdersAbove(e elevatorStatus.Elevator)int{
	for floor := (e.PreviousFloor+1); floor < driver.NUM_FLOORS; floor++{
		if(e.OrderMatrix[floor][0] == 1){
			return 1;
		}
		if(e.OrderMatrix[floor][2] == 1){
			return 1;
		}				
	}
	return 0
}

func CheckDownOrdersAbove(e elevatorStatus.Elevator)int{
	for floor := e.PreviousFloor + 1; floor < driver.NUM_FLOORS; floor++{
		if(e.OrderMatrix[floor][1] ==1){
			return 1;
		} else{
			return 0;
		}
	}
	return 0;	
} 

func CheckUpOrdersBelow(e elevatorStatus.Elevator)int{
	for floor := 0; floor < e.PreviousFloor; floor++{
		if(e.OrderMatrix[floor][0] ==1){
			return 1;
		} else{
			return 0;
		}
	}
	return 0;

}

func CheckDownOrdersBelow(e elevatorStatus.Elevator)int{

	for floor := 0; floor < e.PreviousFloor; floor++{				
			if(e.OrderMatrix[floor][1] == 1){
				return 1;
			} else if(e.OrderMatrix[floor][2] == 1){
				return 1;
			}	
	}
	return 0;
}

func LengthOfQueue(e elevatorStatus.Elevator)int{
	length := 0
	for floor := 0; floor < driver.NUM_FLOORS; floor++{
		for button := 0; button < driver.NUM_BUTTONS; button++{
			if (button == 1 && floor == 0) || (button == 2 && floor == 3){
			}else{
				length += e.OrderMatrix[floor][button] //lag matrisen!
			}
		}
	} 
	return length
}

func NewOrderAtCurrentFloor(e elevatorStatus.Elevator)int{
	floor := driver.Get_floor_sensor_signal()
	if(floor==3){	
		if (e.OrderMatrix[floor][1] == 1){
			return 1;
		} else if (e.OrderMatrix[floor][2] == 1){
			return 1;
		}
	} else if(floor==0){	
		if (e.OrderMatrix[floor][0] == 1){
			return 1;
		} else if (e.OrderMatrix[floor][2] == 1){
			return 1;
		}
	} 

	if(e.Dir != -1){
		if (e.OrderMatrix[floor][0] == 1){
			return 1;
		} else if (e.OrderMatrix[floor][2] == 1){
			return 1;						
		}
	} else if(e.Dir != 1){
		if (e.OrderMatrix[floor][1] == 1){
			return 1;
		} else if (e.OrderMatrix[floor][2] == 1){
			return 1;						
		}
	}
	return 0;
}

func DeleteCompletedOrders(e elevatorStatus.Elevator){
	floor :=  driver.Get_floor_sensor_signal()

	if floor != -1{
		if floor == 0{
			e.OrderMatrix[floor][0] = 0
			e.OrderMatrix[floor][2] = 0
		} else if floor == 3{
			e.OrderMatrix[floor][1] = 0
			e.OrderMatrix[floor][2] = 0
		}

		if e.Dir== 1{
			e.OrderMatrix[floor][0] = 0
			e.OrderMatrix[floor][2] = 0
		} else if e.Dir== -1{
			e.OrderMatrix[floor][1] = 0
			e.OrderMatrix[floor][2] = 0
		} else if e.Dir== 0{
			e.OrderMatrix[floor][0] = 0
			e.OrderMatrix[floor][1] = 0
			e.OrderMatrix[floor][2] = 0
		}
	}

}
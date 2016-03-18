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
	floor := driver.Get_floor_sensor_signal()
	fmt.Println ("Floor: ", floor)
	var result int = 0

	//SETT PREVIOUS FLOOR HER!!
	var previousFloor int = e.PreviousFloor
	e.PreviousFloor = floor
	
	if(floor > previousFloor){
		if (e.OrderMatrix[floor][2] == 1) || (e.OrderMatrix[floor][0] == 1){
			result = 1	
		}
		if floor == 3{
			result = 1
		}
		if (CheckUpOrdersAbove(e) != 1 && CheckDownOrdersAbove(e) != 1){
			result = 1				
		}
	} else if(floor < previousFloor){
		if (e.OrderMatrix[floor][2] == 1) || (e.OrderMatrix[floor][1] == 1){
			result = 1
		}
		if floor == 0{
			result = 1
		}
		if (CheckUpOrdersBelow(e) != 1 && CheckDownOrdersBelow(e) != 1){
			result = 1
		}
	}
	return result
}

func CheckUpOrdersAbove(e elevatorStatus.Elevator)int{
	var result int = 0
	for floor := e.PreviousFloor+1; floor < driver.NUM_FLOORS; floor++{
		if(e.OrderMatrix[floor][0] == 1){
			fmt.Println("Fant en OPP-bestilling over")
			result = 1
		}
		if(e.OrderMatrix[floor][2] == 1){
			result = 1
		}				
	}
	return result
}

func CheckDownOrdersAbove(e elevatorStatus.Elevator)int{
	var result int = 0
	for floor := e.PreviousFloor+1; floor < driver.NUM_FLOORS; floor++{
		if(e.OrderMatrix[floor][1] == 1){
			fmt.Println("Fant en NED-bestilling over")
			result = 1
		}
	}
	return result	
} 

func CheckUpOrdersBelow(e elevatorStatus.Elevator)int{
	var result int = 0
	for floor := 0; floor < e.PreviousFloor; floor++{
		if(e.OrderMatrix[floor][0] == 1){
			fmt.Println("Fant en OPP-bestilling under")
			result = 1
		}
	}
	return result
}

func CheckDownOrdersBelow(e elevatorStatus.Elevator)int{
	var result int = 0
	for floor := 0; floor < e.PreviousFloor; floor++{				
			if(e.OrderMatrix[floor][1] == 1){
				fmt.Println("Fant en NED-bestilling under")
				result = 1
			} else if(e.OrderMatrix[floor][2] == 1){
				result = 1
			}	
	}
	return result
}

func LengthOfQueue(e elevatorStatus.Elevator)int{
	fmt.Println(e.OrderMatrix)
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
	var result int = 0
	if(floor==3){	
		if (e.OrderMatrix[floor][1] == 1){
			result = 1
		} else if (e.OrderMatrix[floor][2] == 1){
			result = 1
		}
	} else if(floor==0){	
		if (e.OrderMatrix[floor][0] == 1){
			result = 1
		} else if (e.OrderMatrix[floor][2] == 1){
			result = 1
		}
	} 

	if(e.Dir != driver.MDIR_DOWN){
		if (e.OrderMatrix[floor][0] == 1){
			result = 1
		} else if (e.OrderMatrix[floor][2] == 1){
			result = 1						
		}
	} else if(e.Dir != driver.MDIR_UP){
		if (e.OrderMatrix[floor][1] == 1){
			result = 1
		} else if (e.OrderMatrix[floor][2] == 1){
			result = 1						
		}
	}
	return result
}

func DeleteCompletedOrders(e *elevatorStatus.Elevator){
	floor :=  driver.Get_floor_sensor_signal()
	fmt.Println("Sletter utført bestilling", e.Dir, floor)
	if floor != -1{
		if floor == 0{
			e.OrderMatrix[floor][0] = 0
			e.OrderMatrix[floor][2] = 0
		} else if floor == 3{
			e.OrderMatrix[floor][1] = 0
			e.OrderMatrix[floor][2] = 0
		}

		if e.Dir == driver.MDIR_UP{
			fmt.Println("Sletter når retn er OPP")
			e.OrderMatrix[floor][0] = 0
			e.OrderMatrix[floor][2] = 0
			if (CheckUpOrdersAbove(*e) != 1 && CheckDownOrdersAbove(*e) != 1){
				e.OrderMatrix[floor][1] = 0
			}

		} else if e.Dir == driver.MDIR_DOWN{
			fmt.Println("Sletter når retn er NED")
			e.OrderMatrix[floor][1] = 0
			e.OrderMatrix[floor][2] = 0
			if (CheckUpOrdersBelow(*e) != 1 && CheckDownOrdersBelow(*e) != 1){
				e.OrderMatrix[floor][0] = 0
			}
		} else if e.Dir == driver.MDIR_STOP{
			e.OrderMatrix[floor][0] = 0
			e.OrderMatrix[floor][1] = 0
			e.OrderMatrix[floor][2] = 0
		}
		fmt.Println(e.OrderMatrix)
	}

}
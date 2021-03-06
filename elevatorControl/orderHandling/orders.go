package orderHandling

import(
	"fmt"
	"time"
	"../driver"
	"../elevatorStatus"
)



// this is for one elevator, når n heiser, skal få ordre fra master og ikke les av knappetrykk

// For tre heiser: Dette skal kun gjøres når knappene inni heisen trykkes
func AddOrderToQueue(e elevatorStatus.Elevator, order [2]int) elevatorStatus.Elevator{
	fmt.Println("inne i addorderqueue")
	floor := order[0]
	button := order[1]
	e.OrderMatrix[floor][button] = 1




	// for floor := 0; floor < driver.NUM_FLOORS; floor++{
	// 	for button := driver.BUTTON_CALL_UP; button <= driver.BUTTON_COMMAND ; button++{
	// 		if (floor == 0 && button == 1) || (floor == 3 && button == 0) || (floor < 0){
	// 		} else if {
	// 			e.OrderMatrix[floor][button] = 1
	// 			// sett knappelys
	// 		}
	// 	}
	// }
	return e
}

//for n antall heiser.  ordren skal sendes til master i stedet for å legges inn i matrisen for n heiser
func ReadButtons(buttonChan chan [2]int){
	var order [2]int
	var prevOrder [2]int
	for{
		for floor := 0; floor < driver.NUM_FLOORS; floor++{
			for button := driver.BUTTON_CALL_UP; button < driver.NUM_BUTTONS; button++{
				if (floor == 0 && button == 1) || (floor == 3 && button == 0) || (floor < 0){
				} else if driver.Get_button_signal(button, floor) == 1{
					// sett knappelys

					order[0] = floor
					order[1] = int(button)


					if prevOrder != order{
						fmt.Println("Button pushed")
						buttonChan <-order
						prevOrder = order
						
					}
				}
			}
		}
	time.Sleep(time.Millisecond * 10)
	}
}


func UpdateOrderMatrix(update [4][3]int, old [4][3]int)[4][3]int{
	for floor := 0; floor < driver.NUM_FLOORS; floor++{
		for button := driver.BUTTON_CALL_UP; button < driver.NUM_BUTTONS; button++{
			if update[floor][button] == 1{
				old[floor][button] = 1
			}
		}
	}
	return old
}


func ShouldStop(e elevatorStatus.Elevator)int{
	//Lag denne senere - kanskje i queue
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	fmt.Println ("Floor: ", e.CurrentFloor)
	var result int = 0

	//SETT PREVIOUS FLOOR HER!!
	var previousFloor int = e.PreviousFloor
	e.PreviousFloor = e.CurrentFloor
	
	if(e.CurrentFloor > previousFloor){
		if (e.OrderMatrix[e.CurrentFloor][2] == 1) || (e.OrderMatrix[e.CurrentFloor][0] == 1){
			result = 1	
		}
		if e.CurrentFloor == 3{
			result = 1
		}
		if (CheckUpOrdersAbove(e) != 1 && CheckDownOrdersAbove(e) != 1){
			result = 1				
		}
	} else if(e.CurrentFloor < previousFloor){
		if (e.OrderMatrix[e.CurrentFloor][2] == 1) || (e.OrderMatrix[e.CurrentFloor][1] == 1){
			result = 1
		}
		if e.CurrentFloor == 0{
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
	for floor := e.CurrentFloor+1; floor < driver.NUM_FLOORS; floor++{
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
	for floor := e.CurrentFloor+1; floor < driver.NUM_FLOORS; floor++{
		if(e.OrderMatrix[floor][1] == 1){
			fmt.Println("Fant en NED-bestilling over")
			result = 1
		}
	}
	return result	
} 

func CheckUpOrdersBelow(e elevatorStatus.Elevator)int{
	var result int = 0
	for floor := 0; floor < e.CurrentFloor; floor++{
		if(e.OrderMatrix[floor][0] == 1){
			fmt.Println("Fant en OPP-bestilling under")
			result = 1
		}
	}
	return result
}

func CheckDownOrdersBelow(e elevatorStatus.Elevator)int{
	var result int = 0
	for floor := 0; floor < e.CurrentFloor; floor++{				
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
	//fmt.Println(e.OrderMatrix)
	length := 0
	for floor := 0; floor < driver.NUM_FLOORS; floor++{
		for button := 0; button < driver.NUM_BUTTONS; button++{
			if (button == 1 && floor == 0) || (button == 0 && floor == 3){
			}else{
				length += e.OrderMatrix[floor][button] 
			}
		}
	} 
	return length
}

func NewOrderAtCurrentFloor(e elevatorStatus.Elevator)int{
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	var result int = 0
	if(e.CurrentFloor==3){	
		if (e.OrderMatrix[e.CurrentFloor][1] == 1){
			result = 1
		} else if (e.OrderMatrix[e.CurrentFloor][2] == 1){
			result = 1
		}
	} else if(e.CurrentFloor==0){	
		if (e.OrderMatrix[e.CurrentFloor][0] == 1){
			result = 1
		} else if (e.OrderMatrix[e.CurrentFloor][2] == 1){
			result = 1
		}
	} 

	if(e.Dir != driver.MDIR_DOWN){
		if (e.OrderMatrix[e.CurrentFloor][0] == 1){
			result = 1
		} else if (e.OrderMatrix[e.CurrentFloor][2] == 1){
			result = 1						
		}
	} 

	if(e.Dir != driver.MDIR_UP){
		if (e.OrderMatrix[e.CurrentFloor][1] == 1){
			result = 1
		} else if (e.OrderMatrix[e.CurrentFloor][2] == 1){
			result = 1						
		}
	}
	return result
}

func DeleteCompletedOrders(e *elevatorStatus.Elevator, DelOrder chan [4]int){
	e.CurrentFloor =  driver.Get_floor_sensor_signal()
	DeleteOrder := [4]int{0, 0, 0, 0}
	DeleteOrder[3] = e.CurrentFloor
	fmt.Println("Sletter utført bestilling", e.Dir, e.CurrentFloor)
	if e.CurrentFloor != -1{
		if e.CurrentFloor == 0{
			e.OrderMatrix[e.CurrentFloor][0], e.OrderMatrix[e.CurrentFloor][2] = 0,0
			DeleteOrder[0], DeleteOrder[2] = 1,1
			
		} else if e.CurrentFloor == 3{
			e.OrderMatrix[e.CurrentFloor][1],e.OrderMatrix[e.CurrentFloor][2] = 0,0
			DeleteOrder[1], DeleteOrder[2] = 1,1
		}

		if e.Dir == driver.MDIR_UP{
			fmt.Println("Sletter når retn er OPP")
			e.OrderMatrix[e.CurrentFloor][0], e.OrderMatrix[e.CurrentFloor][2] = 0,0
			DeleteOrder[0], DeleteOrder[2] = 1,1
			if (CheckUpOrdersAbove(*e) != 1 && CheckDownOrdersAbove(*e) != 1){
				e.OrderMatrix[e.CurrentFloor][1] = 0
				DeleteOrder[1]= 1
			}

		} else if e.Dir == driver.MDIR_DOWN{
			fmt.Println("Sletter når retn er NED")
			e.OrderMatrix[e.CurrentFloor][1], e.OrderMatrix[e.CurrentFloor][2] = 0,0
			DeleteOrder[1], DeleteOrder[2] = 1,1
			if (CheckUpOrdersBelow(*e) != 1 && CheckDownOrdersBelow(*e) != 1){
				e.OrderMatrix[e.CurrentFloor][0] = 0
				DeleteOrder[0]= 1
			}
		} else if e.Dir == driver.MDIR_STOP{
			e.OrderMatrix[e.CurrentFloor][0], e.OrderMatrix[e.CurrentFloor][1], e.OrderMatrix[e.CurrentFloor][2] = 0,0,0
			DeleteOrder[0], DeleteOrder[1], DeleteOrder[2] = 1,1,1
		}
		fmt.Println(e.OrderMatrix)
	}

	DelOrder <-DeleteOrder

}
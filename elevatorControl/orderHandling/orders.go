package orderHandling

import (
	"../../network"
	"../driver"
	"../elevatorStatus"
	"fmt"
	"io/ioutil"
	"time"
)

func AddOrderToQueue(e elevatorStatus.Elevator, order [2]int) elevatorStatus.Elevator {
	floor := order[0]
	button := order[1]
	e.OrderMatrix[floor][button] = 1
	return e
}

func ReadButtons(buttonPushed chan [2]int, elevObject chan elevatorStatus.Elevator) {
	var order [2]int
	var prevOrder [2]int
	prevOrder[0] = -1
	prevOrder[1] = -1
	for {
		for floor := 0; floor < driver.NUM_FLOORS; floor++ {
			for button := 0; button < driver.NUM_BUTTONS; button++ {
				if (floor == 0 && button == 1) || (floor == 3 && button == 0) || (floor < 0) {
				} else if driver.Get_button_signal(button, floor) == 1 {
					order[0] = floor
					order[1] = button
					if prevOrder != order {
						buttonPushed <- order
						prevOrder = order
						if order[1] == driver.BUTTON_COMMAND {
							driver.Set_button_lamp(order[1], order[0], 1)
						}
					}
				} else if driver.Get_button_signal(button, floor) == 0 && prevOrder[0] != -1 {
					prevOrder[0] = -1
					prevOrder[1] = -1
				}
			}
		}
		time.Sleep(time.Millisecond * 10)
	}
}

func UpdateOrderMatrix(update [driver.NUM_FLOORS][driver.NUM_BUTTONS]int, elevObject chan elevatorStatus.Elevator) {
	e := <-elevObject
	for floor := 0; floor < driver.NUM_FLOORS; floor++ {
		for button := driver.BUTTON_CALL_UP; button < driver.NUM_BUTTONS; button++ {
			if update[floor][button] == 1 {
				e.OrderMatrix[floor][button] = 1
			}
		}
	}
	elevObject <- e
}

func WriteInternalsToFile(matrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int) {
	var internalOrders [driver.NUM_FLOORS]int
	for floor := 0; floor < driver.NUM_FLOORS; floor++ {
		internalOrders[floor] = matrix[floor][2]
	}
	buffer := make([]byte, driver.NUM_FLOORS)
	for floor := 0; floor < driver.NUM_FLOORS; floor++ {
		buffer[floor] = byte(internalOrders[floor])
	}
	ioutil.WriteFile("InternalOrders.txt", buffer, 0644)
}

func ReadInternalsToFile() [driver.NUM_FLOORS][driver.NUM_BUTTONS]int {
	buffer, err := ioutil.ReadFile("InternalOrders.txt")
	if err != nil {
		fmt.Println("Error in opening file")
	}
	var internalOrders [driver.NUM_FLOORS]int
	for floor := 0; floor < driver.NUM_FLOORS; floor++ {
		internalOrders[floor] = int(buffer[floor])
	}
	var matrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int
	for floor := 0; floor < driver.NUM_FLOORS; floor++ {
		matrix[floor][2] = internalOrders[floor]
		if matrix[floor][2] == 1 {
			driver.Set_button_lamp(2, floor, 1)
		}
	}
	return matrix
}

func ShouldStop(e elevatorStatus.Elevator) bool {
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	fmt.Println("Floor: ", e.CurrentFloor)
	var shouldStop bool = false
	var previousFloor int = e.PreviousFloor
	e.PreviousFloor = e.CurrentFloor
	if e.CurrentFloor > previousFloor {
		if (e.OrderMatrix[e.CurrentFloor][2] == 1) || (e.OrderMatrix[e.CurrentFloor][0] == 1) {
			shouldStop = true
		}
		if e.CurrentFloor == 3 {
			shouldStop = true
		}
		if CheckUpOrdersAbove(e) != true && CheckDownOrdersAbove(e) != true {
			shouldStop = true
		}
	} else if e.CurrentFloor < previousFloor {
		if (e.OrderMatrix[e.CurrentFloor][2] == 1) || (e.OrderMatrix[e.CurrentFloor][1] == 1) {
			shouldStop = true
		}
		if e.CurrentFloor == 0 {
			shouldStop = true
		}
		if CheckUpOrdersBelow(e) != true && CheckDownOrdersBelow(e) != true {
			shouldStop = true
		}
	}
	return shouldStop
}

func CheckUpOrdersAbove(e elevatorStatus.Elevator) bool {
	var found bool = false
	for floor := e.CurrentFloor + 1; floor < driver.NUM_FLOORS; floor++ {
		if e.OrderMatrix[floor][0] == 1 {
			found = true
		}
		if e.OrderMatrix[floor][2] == 1 {
			found = true
		}
	}
	return found
}

func CheckDownOrdersAbove(e elevatorStatus.Elevator) bool {
	var found bool = false
	for floor := e.CurrentFloor + 1; floor < driver.NUM_FLOORS; floor++ {
		if e.OrderMatrix[floor][1] == 1 {
			found = true
		}
	}
	return found
}

func CheckUpOrdersBelow(e elevatorStatus.Elevator) bool {
	var found bool = false
	for floor := 0; floor < e.CurrentFloor; floor++ {
		if e.OrderMatrix[floor][0] == 1 {
			found = true
		}
	}
	return found
}

func CheckDownOrdersBelow(e elevatorStatus.Elevator) bool {
	var found bool = false
	for floor := 0; floor < e.CurrentFloor; floor++ {
		if e.OrderMatrix[floor][1] == 1 {
			found = true
		} else if e.OrderMatrix[floor][2] == 1 {
			found = true
		}
	}
	return found
}

func LengthOfQueue(e elevatorStatus.Elevator) int {
	length := 0
	for floor := 0; floor < driver.NUM_FLOORS; floor++ {
		for button := 0; button < driver.NUM_BUTTONS; button++ {
			if (button == 1 && floor == 0) || (button == 0 && floor == 3) {
			} else {
				length += e.OrderMatrix[floor][button]
			}
		}
	}
	return length
}

func NewOrderAtCurrentFloor(e elevatorStatus.Elevator) bool {
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	var found bool = false
	if e.CurrentFloor == 3 {
		if e.OrderMatrix[e.CurrentFloor][1] == 1 {
			found = true
		} else if e.OrderMatrix[e.CurrentFloor][2] == 1 {
			found = true
		}
	} else if e.CurrentFloor == 0 {
		if e.OrderMatrix[e.CurrentFloor][0] == 1 {
			found = true
		} else if e.OrderMatrix[e.CurrentFloor][2] == 1 {
			found = true
		}
	}
	if e.Direction != driver.MDIR_DOWN {
		if e.OrderMatrix[e.CurrentFloor][0] == 1 {
			found = true
		} else if e.OrderMatrix[e.CurrentFloor][2] == 1 {
			found = true
		}
	}
	if e.Direction != driver.MDIR_UP {
		if e.OrderMatrix[e.CurrentFloor][1] == 1 {
			found = true
		} else if e.OrderMatrix[e.CurrentFloor][2] == 1 {
			found = true
		}
	}
	return found
}

func DeleteCompletedOrders(e *elevatorStatus.Elevator, DelOrder chan [4]int) {
	e.CurrentFloor = driver.Get_floor_sensor_signal()
	DeleteOrder := [4]int{0, 0, 0, 0} //{button up, button down, internal button, floor}
	DeleteOrder[3] = e.CurrentFloor
	if e.CurrentFloor != -1 {
		if e.CurrentFloor == 0 {
			e.OrderMatrix[e.CurrentFloor][0], e.OrderMatrix[e.CurrentFloor][2] = 0, 0
			DeleteOrder[0], DeleteOrder[2] = 1, 1
		} else if e.CurrentFloor == 3 {
			e.OrderMatrix[e.CurrentFloor][1], e.OrderMatrix[e.CurrentFloor][2] = 0, 0
			DeleteOrder[1], DeleteOrder[2] = 1, 1
		}
		if e.Direction == driver.MDIR_UP {
			e.OrderMatrix[e.CurrentFloor][0], e.OrderMatrix[e.CurrentFloor][2] = 0, 0
			DeleteOrder[0], DeleteOrder[2] = 1, 1
			if CheckUpOrdersAbove(*e) != true && CheckDownOrdersAbove(*e) != true {
				e.OrderMatrix[e.CurrentFloor][1] = 0
				DeleteOrder[1] = 1
			}
		} else if e.Direction == driver.MDIR_DOWN {
			e.OrderMatrix[e.CurrentFloor][1], e.OrderMatrix[e.CurrentFloor][2] = 0, 0
			DeleteOrder[1], DeleteOrder[2] = 1, 1
			if CheckUpOrdersBelow(*e) != true && CheckDownOrdersBelow(*e) != true {
				e.OrderMatrix[e.CurrentFloor][0] = 0
				DeleteOrder[0] = 1
			}
		} else if e.Direction == driver.MDIR_STOP {
			e.OrderMatrix[e.CurrentFloor][0], e.OrderMatrix[e.CurrentFloor][1], e.OrderMatrix[e.CurrentFloor][2] = 0, 0, 0
			DeleteOrder[0], DeleteOrder[1], DeleteOrder[2] = 1, 1, 1
		}
	}
	if DeleteOrder[2] == 1 {
		floor := DeleteOrder[3]
		driver.Set_button_lamp(2, floor, 0)
	}
	if network.GetIpAddress() == "::1" {
		for button := 0; button < driver.NUM_BUTTONS-1; button++ {
			floor := DeleteOrder[3]
			driver.Set_button_lamp(button, floor, 0)
		}
	}
	DelOrder <- DeleteOrder
}

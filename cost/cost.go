package cost

import(
	"fmt"
	"../elevatorControl/driver"
	"../message"
	"../elevatorControl/orderHandling"
	"../elevatorControl/elevatorStatus"
)

func absValue(sum int)int{
	if sum < 0{
		return -sum
	} else{
		return sum
	}
}

func CalculateCost(elevator elevatorStatus.Elevator, Order message.UpdateMessage)int{
	sum := elevator.CurrentFloor-Order.Order[0]
	var distanceCost int = absValue(sum)*5
	var directionCost int = -1
	belowOrAbove := elevator.CurrentFloor-Order.Order[1] 
	elevDir := elevator.Dir

	if elevDir == driver.MDIR_STOP{
		directionCost = 0
	} else if (elevDir == driver.MDIR_DOWN && belowOrAbove > 0) || (elevDir == driver.MDIR_UP && belowOrAbove < 0){
		directionCost = 10
	} else if (elevDir == driver.MDIR_DOWN && belowOrAbove < 0) || (elevDir == driver.MDIR_UP && belowOrAbove > 0){
		directionCost = 40
	} else if (elevDir == driver.MDIR_DOWN || elevDir == driver.MDIR_UP) && (belowOrAbove == 0){
		directionCost = 40
	}

	queueCost := 20 * orderHandling.LengthOfQueue(elevator)
	totalCost := distanceCost + directionCost + queueCost
	
	fmt.Println("distanceCost: ", distanceCost)
	fmt.Println("directionCost: ", directionCost)
	fmt.Println("totalCost: ", totalCost)

	return totalCost
	
}


func AssignOrdersToElevator(order message.UpdateMessage, elevators map[string]*elevatorStatus.Elevator)string{
	min_value := 1000 
	neword := order.Order
	fmt.Println("Ny ordre: ", neword)
	var assignedElev string //elevatorStatus.Elevator.ElevatorId
	for _, elev := range elevators {
		value := CalculateCost(*elev, order)
		if value < min_value {
			min_value = value
			assignedElev = elev.IP
		}
	}
	return assignedElev
}






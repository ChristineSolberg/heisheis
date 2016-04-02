package cost

import(
	//"../elevatorControl/driver"
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
	sum := elevator.CurrentFloor-Order.Order[1]
	var distanceCost int = absValue(sum)*5
	var directionCost int = -1
	belowOrAbove := elevator.CurrentFloor-Order.Order[1] 
	elevDir := elevator.CurrentFloor

	if elevator.CurrentFloor == message.Idle{
		directionCost = 0
	} else if (elevDir == message.Down && belowOrAbove > 0) || (elevDir == message.Up && belowOrAbove < 0){
		directionCost = 10
	} else if (elevDir == message.Down && belowOrAbove < 0) || (elevDir == message.Up && belowOrAbove > 0){
		directionCost = 40
	} else if (elevDir == message.Down || elevDir == message.Up) && (belowOrAbove == 0){
		directionCost = 40
	}

	queueCost := 20 * orderHandling.LengthOfQueue(elevator)
	totalCost := distanceCost + directionCost + queueCost
	
	return totalCost
	
}


func AssignOrdersToElevator(order message.UpdateMessage, elevators map[string]elevatorStatus.Elevator)string{
	min_value := 1000 
	var assignedElev string //elevatorStatus.Elevator.ElevatorId
	for _, elev := range elevators {
		value := CalculateCost(elev, order)
		if value < min_value {
			min_value = value
			assignedElev = order.RecieverIP
		}
	}
	return assignedElev
}






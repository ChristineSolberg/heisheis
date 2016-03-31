package cost

import(
	"../driver"
	"../message"
	"math"
)


func CalculateCost(elevator elevatorStatus.Elevator, Order message.UpdateMessage)int{
	var distanceCost int := math.Abs(elevator.CurrentFloor-Order.NewOrder[1])*5
	var directionCost int := -1
	belowOrAbove := elevator.CurrentFloor-Order.NewOrder[1] 
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

	queueCost := 20 * lengthOfQueue(elevator)
	totalCost := distanceCost + directionCost + queueCost
	
	return totalCost
	
}


func AssignOrdersToElevator(order message.UpdateMessage, elevators []elevatorStatus.Elevator, networkChan chan UpdateMessage){
	min_value := 1000 
	var assignedElev elevatorStatus.Elevator.ElevatorId
	for _, elev := range elevators {
		value := CalculateCost(elev, order)
		if value < min_value {
			min_value = value
			assignedElev = elev
		}
	}
	networkChan <- assignedElev
}



func AssignOrdersToElevator(order message.UpdateMessage, elevators []elevatorStatus.Elevator){
	min_value := 1000 
	var min_elev elevatorStatus.Elevator.ElevatorId
	for elev, cost := range elevators {
		value := CalculateCost(elev, order)
		if value < min_value {
			min_value = value
			min_elev = elev
		}
	}


func lengthOfQueue(elevator message.UpdateMessage)int{
	length := 0
	for floor := 0; floor < driver.NUM_FLOORS; floor++{
		for button := 0; button < driver.NUM_BUTTONS; button++{
			if (button == 1 && floor == 0) || (button == 2 && floor == 3){
			}else{
				length += elevator.OrderMatrix[button][floor]
			}
		}
	} 
	return length
}
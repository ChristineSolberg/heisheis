package cost

import(
	"../driver"
	"../message"
	"math"
)

var elevator message.UpdateMessage



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

	queueCost := 5 * lengthOfQueue(elevator)
	totalCost := distanceCost + directionCost + queueCost
	
	return totalCost

	
}


func AssignOrdersToElevator(elevator message.UpdateMessage, cost int){
	
	elevCost := make(map[int]int)
	elevCost[elevator.ElevatorId] = cost




	//smallestCost := 1000
	// for (går igjennom alle heiser){
	// 	cost[elev] = CalulateCost(elev)
	// 	if cost[elev] < cost
	// 		oppdater
	//}

	var min int := 1000
	for elevID, cost := range elevCost{

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
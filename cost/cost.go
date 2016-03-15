package cost

import(
	"../driver"
	"../message"
	"math"
)

var elevator message.UpdateMessage



func CalculateCost(elevator message.UpdateMessage, Order message.UpdateMessage)(map[int]int){
	distanceCost := math.Abs(elevator.CurrentState[0]-Order.NewOrder[1])*5
	directionCost := -1
	belowOrAbove := elevator.CurrentState[0]-Order.NewOrder[1]
	elevDir := elevator.CurrentState[1]

	if elevator.CurrentState[1] == message.Idle{
		directionCost = 0
	} else if (elevDir == message.Down && belowOrAbove > 0 )|| (elevDir == message.Up && belowOrAbove < 0) {
		directionCost = 10
	} else if (elevDir == message.Down && belowOrAbove < 0) || (elevDir == message.Up && belowOrAbove > 0){
		directionCost = 40
	} else if (elevDir == message.Down || elevDir == message.Up) && (belowOrAbove == 0){
		directionCost = 40

	}

	queueCost := 5 * lengthOfQueue(elevator)
	totalCost := distanceCost + directionCost + queueCost
	elevCost := make(map[int]int)
	elevCost[id] = totalCost
	return elevCost
}


func AssignOrdersToElevator(elevator message.UpdateMessage){
	//smallestCost := 1000



	// cost:= 1000
	// for (går igjennom alle heiser){
	// 	cost[elev] = CalulateCost(elev)
	// 	if cost[elev] < cost
	// 		oppdater
	//}

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
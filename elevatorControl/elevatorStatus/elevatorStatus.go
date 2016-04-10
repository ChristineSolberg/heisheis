package elevatorStatus

import "../driver"

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
	POWER_OFF					= 4
	NO_EVENT					= 5
)

type Elevator struct{
	Direction driver.MotorDirection
	CurrentFloor int
	PreviousFloor int
	State State
	IP string
	OrderMatrix [driver.NUM_FLOORS][driver.NUM_BUTTONS]int
	Master string
}

func MakeCopyOfElevator(elevChan chan Elevator)Elevator{
	e := <- elevChan
	elevChan <- e
	return e
}
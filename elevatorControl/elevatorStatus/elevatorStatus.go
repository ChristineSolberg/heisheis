package elevatorStatus

import (
	"../driver"
	
	)

type ElevState int
const (
	IDLE  			ElevState 	= 0 
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
	NO_EVENT					= 4
	
)

type Elevator struct{
	Dir driver.MotorDirection
	PreviousFloor int
	State ElevState
	OrderMatrix [4][3]int
	//TimerChan: make(chan string)
}
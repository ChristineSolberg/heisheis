package elevatorStatus

import (
	"../driver"
	//"time"
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
	POWER_OFF					= 4
	NO_EVENT					= 5
	
)

type Elevator struct{
	Dir driver.MotorDirection
	CurrentFloor int
	PreviousFloor int
	State ElevState
	IP string
	OrderMatrix [4][3]int
	Master string
	//DoorTimeout <-chan time.Time
}
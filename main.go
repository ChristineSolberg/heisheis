package main

import (
    "./elevatorControl/driver"
    "./elevatorControl"
    "./elevatorControl/orderHandling"
    "./elevatorControl/elevatorStatus"

    //"./network"
    //"./message"
    "fmt"
    //"time"
)



func main() {
	driver.Init()
	//driver. Set_motor_speed(driver.MDIR_STOP)

	//var OrderMatrix [3][4]int

	// for floor := 0; floor < 4; floor++{
	// 	for button := 0; button < 3; button++{
	// 		if(floor == 0 && button == 1) || (floor == 3 && button == 0){
	// 		}else{
	// 			OrderMatrix [button][floor] = 1
	// 		}
	// 	}
	// }

	//fmt.Println(OrderMatrix)
	
	var e elevatorStatus.Elevator

	e = elevatorControl.StartUp(e)
	for{
		e = orderHandling.AddOrderToQueue(e)
		e = elevatorControl.UpdateFSM(e)
		

	}

    // for {
    //     // Change direction when we reach top/bottom floor
    //     if (driver.Get_floor_sensor_signal() == driver.NUM_FLOORS - 1) {
    //         driver.Set_motor_speed(driver.MDIR_DOWN)
    //     } else if (driver.Get_floor_sensor_signal() == 0) {
    //         driver.Set_motor_speed(driver.MDIR_UP)
    //     }

    //     // Stop elevator and exit program if the stop button is pressed
      
        
    // }


	// fmt.Println("Start main")
 //    recvChan := make(chan message.UpdateMessage)
 //    sendChan := make(chan message.UpdateMessage)

 //    conn1 := network.ServerConnection()
 //    conn2 := network.ClientConnection()
 //    fmt.Println("main2")


 //    var Alive message.UpdateMessage
 //    ticker := time.NewTicker(time.Millisecond*500)
	// fmt.Println("ticker started")
	// Alive.MessageType = message.IAmAlive
    
 //    go message.RecvMsg(conn1,recvChan)
 //    go message.SendMsg(conn2,sendChan)
	// for{
	// 	select{
	// 		case <-ticker.C:
	// 			fmt.Println("Legger på kanalen")
	// 			sendChan<-Alive

	// 		case msg := <-recvChan:
	// 			fmt.Println(msg)
	// 	}
	// }
 


    
    fmt.Println("main3")
    


}
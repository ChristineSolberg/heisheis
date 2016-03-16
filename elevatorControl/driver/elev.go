package driver

import(
	"errors"
	"fmt"
)

var driverInitialized = false

const NUM_FLOORS = 4
const NUM_BUTTONS = 3

type MotorDirection int

const (

	MDIR_DOWN  MotorDirection = -1 // hvorfor motordirection, syntax?
	MDIR_STOP 				  = 0
	MDIR_UP 			      = 1
)

type ButtonEvent struct{
	Floor int
	Button int
}

var lampChannel = [NUM_FLOORS][NUM_BUTTONS] int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var buttonChannels = [NUM_FLOORS][NUM_BUTTONS] int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},

}

func Init() error{
	if driverInitialized{
		return errors.New("Already initialized")
	} else{
		driverInitialized = true
		if (IO_init() == 0){
			return errors.New("IO not initialized")
		} else{
			Init_button_lamps()
			Startup_floor()
			return nil
		}
	}
}	
	


func Init_button_lamps(){

	Set_stop_lamp(0)
	Set_door_open_lamp(0)
    Set_floor_indicator(0)


	for floor := 0; floor < NUM_FLOORS; floor++{
		for button := 0; button < NUM_BUTTONS; button ++{
			if ((floor == 0 && button == 1 )|| (floor == 3 && button == 0)) {
				// these buttons do not exist in our system

			} else{
				Set_button_lamp(button,floor, 0)
			}
		}
	}
}

//move to defined state
func Startup_floor(){
	if (Get_floor_sensor_signal() == -1){
		Set_motor_speed(MDIR_DOWN)
	}
}


// make a buttonPoller? time: 33.58

func Set_motor_speed(dir MotorDirection){
	switch dir{
	case MDIR_DOWN:
		IO_set_bit(MOTORDIR)
		IO_write_analog(MOTOR,2800)

	case MDIR_STOP:
		IO_write_analog(MOTOR,0)

	case MDIR_UP:
		IO_clear_bit(MOTORDIR) // check this out, why not set?
		IO_write_analog(MOTOR,2800)

	}
}

func Set_button_lamp(button int, floor int, value int){

	if (value == 1) {
        IO_set_bit(lampChannel[floor][button]);
    } else {
        IO_clear_bit(lampChannel[floor][button]);
    }
}

func Set_floor_indicator(floor int){
	
	if (floor < 0){
    	fmt.Println("Error: Floor is less than 0")

    }
      if (floor > NUM_FLOORS){
    	fmt.Println("Error: Floor is greater than NUM_FLOORS")

    }

	 if (floor & 0x02 != 0) {
        IO_set_bit(LIGHT_FLOOR_IND1);
    } else {
        IO_clear_bit(LIGHT_FLOOR_IND1);
    }    

    if (floor & 0x01 != 0) {
        IO_set_bit(LIGHT_FLOOR_IND2);
    } else {
        IO_clear_bit(LIGHT_FLOOR_IND2);
    } 

}

func Set_door_open_lamp(value int){
	  if (value == 1) {
        IO_set_bit(LIGHT_DOOR_OPEN);
    } else {
        IO_clear_bit(LIGHT_DOOR_OPEN);
    }
}



func Set_stop_lamp(value int) {
    if (value == 1) {
        IO_set_bit(LIGHT_STOP);
    } else {
        IO_clear_bit(LIGHT_STOP);
    }
}


func Get_floor_sensor_signal() int{

	if (IO_read_bit(SENSOR_FLOOR1) == 1){
        return 0;
    } else if (IO_read_bit(SENSOR_FLOOR2) == 1) {
        return 1;
    } else if (IO_read_bit(SENSOR_FLOOR3) == 1) {
        return 2;
    } else if (IO_read_bit(SENSOR_FLOOR4) == 1) {
        return 3;
    } else {
        return -1;
    }
}

// if we want stop and obstruction :):::):):):):):):):):implement here 
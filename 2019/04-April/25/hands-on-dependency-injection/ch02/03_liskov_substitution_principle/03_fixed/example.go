package fixedv2

type actions interface {
	start()
	drive()
}

type Car struct {
	poweredVehicle
}

func (c Car) start() {
	c.poweredVehicle.startEngine()
}

func (c Car) drive() {
	// TODO: implement
}

type poweredVehicle struct {
}

func (p poweredVehicle) startEngine() {
	// common engine start code
}

type Sled struct {
}

func (s Sled) start() {
	// push start
}

func (s Sled) drive() {
	// TODO: implement
}

package lsp_violation

func Go(vehicle actions) {
	if sled, ok := vehicle.(Sled); ok {
		sled.pushStart()
	} else {
		vehicle.startEngine()
	}

	vehicle.drive()
}

type actions interface {
	drive()
	startEngine()
}

type Vehicle struct {
}

func (v Vehicle) drive() {

}

func (v Vehicle) startEngine() {

}

func (v Vehicle) stopEngine() {

}

type Car struct {
	Vehicle
}

type Sled struct {
	Vehicle
}

func (s Sled) startEngine() {
	// override so that is does nothing
}

func (s Sled) stopEngine() {
	// override so that is does nothing
}

func (s Sled) pushStart() {

}

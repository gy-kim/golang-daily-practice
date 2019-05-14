package main

import "fmt"

type Originator struct {
	state string
}

func (o *Originator) GetState() string {
	return o.state
}

func (o *Originator) SetState(state string) {
	fmt.Println("Setting state to " + state)
	o.state = state
}

func (o *Originator) GetMemento() Memento {
	return Memento{o.state}
}

func (o *Originator) Restore(memento Memento) {
	// restore state
	o.state = memento.GetState()
}

type Memento struct {
	serializedState string
}

func (m *Memento) GetState() string {
	return m.serializedState
}

func Caretaker() {

	// assume that A is the original state of the Orginator
	theOriginator := Originator{"A"}
	theOriginator.SetState("A")
	fmt.Println("theOriginator state = ", theOriginator.GetState())

	// before mutating, get an memento
	theOriginator.SetState("unclean")
	fmt.Println("theOriginator state = ", theOriginator)
}

func main() {
	Caretaker()
}

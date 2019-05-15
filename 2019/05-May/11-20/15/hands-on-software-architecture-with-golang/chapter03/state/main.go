package main

import "fmt"

type State interface {
	Op1(*Context)
	Op2(*Context)
}

type Context struct {
	state State
}

func (c *Context) Op1() {
	c.state.Op1(c)
}

func (c *Context) Op2() {
	c.state.Op2(c)
}

func (c *Context) SetState(state State) {
	c.state = state
}

func NewContext() *Context {
	c := new(Context)
	c.SetState(new(StateA))
	return c
}

type StateA struct{}

func (s *StateA) Op1(c *Context) {
	fmt.Println("State A : Op1 ")
}

func (s *StateA) Op2(c *Context) {
	fmt.Println("State A : Op2 ")
	c.SetState(new(StateB)) // <-- State Change!
}

type StateB struct{}

func (s *StateB) Op1(c *Context) {
	fmt.Println("State B : Op1 ")
}

func (s *StateB) Op2(c *Context) {
	fmt.Println("State B : Op2 ")
	c.SetState(new(StateA))
}

func main() {
	context := NewContext()

	// state operations
	context.Op1()
	context.Op2()
	context.Op1()
	context.Op2()
}

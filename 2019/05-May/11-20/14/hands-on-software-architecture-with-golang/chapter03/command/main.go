package main

import (
	"errors"
	"fmt"
)

// The Command - encapsulates the action to be done
type Report interface {
	Execute() // here the action is called Execute
}

type ConcreteReportA struct {
	receiver *Receiver
}

func (c *ConcreteReportA) Execute() {
	c.receiver.Action("ReportA")
}

type ConcreteReportB struct {
	receiver *Receiver
}

func (c *ConcreteReportB) Execute() {
	c.receiver.Action("ReportB")
}

// Receiver - ancillary object passed to command execution
// this can pass useful information for
type Receiver struct{}

func (r *Receiver) Action(msg string) {
	fmt.Println(msg)
}

// Invoker - this object which knows how to execute a command, and optionally
// does bookkeeping about the command execution.
type Invoker struct {
	repository []Report
}

func (i *Invoker) Schedule(cmd Report) {
	i.repository = append(i.repository, cmd)
}

func (i *Invoker) Run() {
	for _, cmd := range i.repository {
		cmd.Execute()
	}
}

// Chain of Responsibilty
// use Command to represent requests as objects
type ChainedReceiver struct {
	canHandle string
	next      *ChainedReceiver
}

func (r *ChainedReceiver) Finish() error {
	fmt.Println(r.canHandle, " Receiver Finishing")
	return nil
}

func (r *ChainedReceiver) Handle(what string) error {
	// Check if this receiver can handle
	// this of course is a dummy check
	if what == r.canHandle {
		return r.Finish()
	} else if r.next != nil {
		return r.next.Handle(what)
	} else {
		fmt.Println("No Receiver could handle the request!")
		return errors.New("No Receiver to Handle")
	}
}

func main() {
	receiver := new(Receiver)
	ReportA := &ConcreteReportA{receiver}
	ReportB := &ConcreteReportB{receiver}
	invoker := new(Invoker)
	invoker.Schedule(ReportA)
	invoker.Run()
	invoker.Schedule(ReportB)
	invoker.Run()
}

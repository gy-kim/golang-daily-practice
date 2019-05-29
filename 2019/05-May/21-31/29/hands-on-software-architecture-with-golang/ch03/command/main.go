package main

import (
	"errors"
	"fmt"
)

type Report interface {
	Execute()
}

type ConcreteReporterA struct {
	receiver *Receiver
}

func (c *ConcreteReporterA) Execute() {
	c.receiver.Action("ReportA")
}

type ConcreteReporterB struct {
	receiver *Receiver
}

func (c *ConcreteReporterB) Execute() {
	c.receiver.Action("ReportB")
}

// Receiver - ancillary objects passed to command execution
// this can pass  useful information for
type Receiver struct{}

func (r *Receiver) Action(msg string) {
	fmt.Println(msg)
}

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

// Chain Of Responsibilty
// uses command to represent requests as object
type ChainedReceiver struct {
	canHandle string
	next      *ChainedReceiver
}

func (r *ChainedReceiver) SetNext(next *ChainedReceiver) {
	r.next = next
}

func (r *ChainedReceiver) Finish() error {
	fmt.Println(r.canHandle, "Receiver Finishing")
	return nil
}

func (r *ChainedReceiver) Handle(what string) error {
	// Check if this receiver can handle
	// this of course os a dummy check
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
	ReportA := &ConcreteReporterA{receiver}
	ReportB := &ConcreteReporterB{receiver}
	invoker := new(Invoker)
	invoker.Schedule(ReportA)
	invoker.Run()
	invoker.Schedule(ReportB)
	invoker.Run()
}

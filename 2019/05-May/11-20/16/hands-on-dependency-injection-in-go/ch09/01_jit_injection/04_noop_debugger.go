package jit_injection

import "errors"

type ObjectWithDebugger struct {
	Debugger Debugger
}

func (o *ObjectWithDebugger) DoSomethingAmazing(input int) error {
	o.getDebugger().Log("input was : %d", input)

	err := o.doSomething()

	o.getDebugger().Log("result was: %v", err)
	return err
}

func (o *ObjectWithDebugger) getDebugger() Debugger {
	if o.Debugger == nil {
		o.Debugger = &noopDebugger{}
	}

	return o.Debugger
}

func (o *ObjectWithDebugger) doSomething() error {
	return errors.New("not implemented yet")
}

type Debugger interface {
	Log(msg string, args ...interface{})
}

type noopDebugger struct{}

func (n *noopDebugger) Log(_ string, args ...interface{}) {
}

package require

import "testing"

// AssertionTesterInterface defines an interface to be used for testing assertion methods
type AssertionTesterInterface interface {
	TestMethod()
}

// AssertionTesterConformingObejct is an objct that conforms to the AssertionTestInterface interface
type AssertionTesterConformingObject struct{}

func (a *AssertionTesterConformingObject) TestMethod() {}

// AssertionTesterNonConformingObject is an object that does not conform to the AssertionTesterInterface interface
type AssertionTesterNonConformingObject struct{}

type MockT struct {
	Failed bool
}

func (t *MockT) FailNow() {
	t.Failed = true
}

func (t *MockT) Errorf(format string, args ...interface{}) {
	_, _ = format, args
}

func TestImplements(t *testing.T) {

}

package structural

import (
	"errors"
	"strings"
	"testing"
)

func TestPrintAPI1(t *testing.T) {
	api1 := PrinterAPI1{}

	err := api1.PrintMessage("Hello")
	if err != nil {
		t.Errorf("Error trying to use the API1 implementation: Message: %s\n", err.Error())
	}
}

type TestWriter struct {
	Msg string
}

func (t *TestWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 {
		t.Msg = string(p)
		return n, nil
	}

	err = errors.New("Content received on Writer was empty")
	return
}

func TestPrintAPI2(t *testing.T) {
	api2 := PrinterAPI2{}

	err := api2.PrintMessage("Hello")
	if err != nil {
		expectedErrorMessage := "You need to pass an io.Writer to PrinterAPI2"
		if !strings.Contains(err.Error(), expectedErrorMessage) {
			t.Errorf("Error message was not corret.\n Actual: %s\nExpected: %s\n", err.Error(), expectedErrorMessage)
		}
	}

	testWriter := TestWriter{}
	api2 = PrinterAPI2{
		Writer: &testWriter,
	}

	err = api2.PrintMessage("Hello")
	if err != nil {
		t.Errorf("Error trying to use the API2 implementation: Message: %s\n", err.Error())
	}

	expectedMessage := "Hello"
	if testWriter.Msg != expectedMessage {
		t.Fatalf("API2 did not write correctly on the io.Writer. \n  Actual: %s\nExpected: %s\n", testWriter.Msg, expectedMessage)
	}
}

package ocp

import (
	"io"
	"net/http"
)

func BuildOutputOCPFail(response http.ResponseWriter, format string, person Person) {
	var err error

	switch format {
	case "csv":
		err = outputCSV(response, person)
	case "JSON":
		err = outputJSON(response, person)
	}

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

func outputCSV(writer io.Writer, person Person) error {
	return nil
}

func outputJSON(writer io.Writer, person Person) error {
	return nil
}

type Person struct {
	Name  string
	Email string
}

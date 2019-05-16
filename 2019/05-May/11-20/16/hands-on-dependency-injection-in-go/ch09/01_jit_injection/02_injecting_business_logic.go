package jit_injection

import (
	"errors"
	"net/http"
)

type LoadPersonHandler struct {
	businessLogic LoadPersonLogic
}

func (h *LoadPersonHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	requestID, err := h.extractInputFromRequest(request)

	output, err := h.businessLogic.Load(requestID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.writeOutput(response, output)
}

func (h *LoadPersonHandler) extractInputFromRequest(request *http.Request) (int, error) {
	return 0, errors.New("not implemented yet")
}

func (h *LoadPersonHandler) writeOutput(write http.ResponseWriter, person Person) {
	// not implemented yet
}

type LoadPersonLogic interface {
	// Load person by supplied ID
	Load(ID int) (Person, error)
}

package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gy-kim/golang-daily-practice/2019/08-Aug/21-31/26-31/hands-on-dependency-injection-in-go/ch09/acme/internal/modules/data"
)

// RegisterModel will validate and save a registration
//go:generate mockery -name=RegisterModel -case underscore -testonly -inpkg -note @generated
type RegisterModel interface {
	Do(ctx context.Context, in *data.Person) (int, error)
}

// NewRegisterHandler is the constructor for RegisterHandler
func NewRegisterHandler(model RegisterModel) *RegisterHandler {
	return &RegisterHandler{
		registerer: model,
	}
}

// RegisterHandler is the HTTP handler for the "Register" endpoint
// In this simplified example we are assuming all possible errors are errors and returning "bad request" HTTP 400.
// There are some programmer errors possible but hopefully these will be caught in testing.
type RegisterHandler struct {
	registerer RegisterModel
}

// ServeHTTP implements http.Handler
func (h *RegisterHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// set latency budget for this API
	subCtx, cancel := context.WithTimeout(request.Context(), 1500*time.Millisecond)
	defer cancel()

	// extract payload from request
	requestPayload, err := h.extractPayload(request)
	if err != nil {
		// output error
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// call the business logic using the request data and context
	id, err := h.register(subCtx, requestPayload)
	if err != nil {
		// not need to log here as we can expect other layers to do so
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// happy path
	response.Header().Add("Location", fmt.Sprintf("/person/%d/", id))
	response.WriteHeader(http.StatusCreated)
}

func (h *RegisterHandler) extractPayload(request *http.Request) (*registerRequest, error) {
	requestPayload := &registerRequest{}

	decode := json.NewDecoder(request.Body)
	err := decode.Decode(requestPayload)
	if err != nil {
		return nil, err
	}

	return requestPayload, nil
}

// call the logic layer
func (h *RegisterHandler) register(ctx context.Context, requestPayload *registerRequest) (int, error) {
	person := &data.Person{
		FullName: requestPayload.FullName,
		Phone:    requestPayload.Phone,
		Currency: requestPayload.Currency,
	}

	return h.registerer.Do(ctx, person)
}

// register endpoint request format
type registerRequest struct {
	// FullName of the person
	FullName string `json:"fullName"`
	// Phone of the person
	Phone string `json:"phone"`
	// Currency the wish to register in
	Currency string `json:"currency"`
}

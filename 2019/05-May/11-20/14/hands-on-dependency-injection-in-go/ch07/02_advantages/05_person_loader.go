package advantages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// thrown then the supplied order does not exist in the database
	errNotFound = errors.New("order not found")
)

// Loads orders based on supplied owner and order OD
type OrderLoader interface {
	loadOrder(owner Owner, orderID int) (Order, error)
}

func NewLoadOrderHandler(loader OrderLoader) *LoadOrderHandler {
	return &LoadOrderHandler{
		loader: loader,
	}
}

type LoadOrderHandler struct {
	loader OrderLoader
}

func (l *LoadOrderHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// extract user from supplied authentication credentials.
	currentUser, err := l.authenticateUser(request)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	// extract order ID from request
	orderID, err := l.extractOrderID(request)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	// load order using the current user as a request-scoped dependency
	// (with method injection)
	order, err := l.loader.loadOrder(currentUser, orderID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// output order
	encoder := json.NewEncoder(response)
	err = encoder.Encode(order)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

// AuthenticateLoader will load orders for based on the supplied owner.
type AuthenticateLoader struct {
	// This pool is  expensive to create. We will want to create it once then reuse it.
	db *sql.DB
}

func (a *AuthenticateLoader) loadByOwner(owner Owner, orderID int) (*Order, error) {
	order, err := a.load(orderID)
	if err != nil {
		return nil, err
	}

	if order.OwnerID != owner.ID() {
		// Return not found so we do not leak information to hackers
		return nil, errNotFound
	}

	// happy path
	return order, nil
}

func (a *AuthenticateLoader) load(orderID int) (*Order, error) {
	// load order from DB
	return &Order{OwnerID: 1}, nil
}

type Owner interface {
	ID() int
}

type Order struct {
	OwnerID int
}

type User struct {
	id int
}

func (u *User) ID() int {
	return u.id
}

func (l *LoadOrderHandler) authenticateUser(request *http.Request) (*User, error) {
	return &User{id: 1}, nil
}

// Extract the order ID from the request
func (l *LoadOrderHandler) extractOrderID(request *http.Request) (int, error) {
	return 2, nil
}

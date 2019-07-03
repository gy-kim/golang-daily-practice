package advantages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	errNotFound = errors.New("order not found")
)

// Loads orders based on supploed owner and ordr ID
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
	// extract user from supplied authentication credentials
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

	order, err := l.loader.loadOrder(currentUser, orderID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ouput order
	encoder := json.NewEncoder(response)
	err = encoder.Encode(order)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}

type AuthenticatedLoader struct {
	// this pool is expecsive to create. We will want to create it once and reuse it.
	db *sql.DB
}

func (a *AuthenticatedLoader) loadByOwner(owner Owner, orderID int) (*Order, error) {
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

func (a *AuthenticatedLoader) load(orderID int) (*Order, error) {
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

func (l *LoadOrderHandler) extractOrderID(request *http.Request) (int, error) {
	return 2, nil
}

package main

// Reservation and Invoice as two generic products
type Reservation interface{}
type Invoice interface{}

// AbstractFactory is the abstract factory which will create two products - both which need to work with each other
type AbstractFactory interface {
	CreateReservation() Reservation
	CreateInvoice() Invoice
}

// HotelFactory implements AbstractFactory and creates products for the Hotel family of products
type HotelFactory struct{}

func (f HotelFactory) CreateReservation() Reservation {
	return new(HotelReservation)
}

func (f HotelFactory) CreateInvoice() Invoice {
	return new(HotelInvoice)
}

type FlightFactory struct{}

func (f FlightFactory) CreateReservation() Reservation {
	return new(FlightReservation)
}

func (f FlightFactory) CreateInvoice() Invoice {
	return new(FlightReservation)
}

func GetFactory(vertical string) AbstractFactory {
	var factory AbstractFactory
	switch vertical {
	case "flight":
		factory = FlightFactory{}
	case "hotel":
		factory = HotelFactory{}
	}

	return factory
}

type HotelReservation struct{}
type HotelInvoice struct{}
type FlightReservation struct{}
type FlightInvoice struct{}

func main() {
}

package main

import "fmt"

type Reservation interface {
	GetReservationDate() string
	CalculateCancellationFee() float64
}

type HotelReservation interface {
	Reservation
	ChangeType()
}

type FlightReservation interface {
	ReservationAddExtraLuggageAllowance(peices int)
}

type HotelReservationImpl struct {
	reservationDate string
}

func (r HotelReservationImpl) GetReservationDate() string {
	return r.reservationDate
}

func (r HotelReservationImpl) CalculateCancellationFee() float64 {
	return 1.0
}

type FlightReservationImpl struct {
	reservationDate string
	luggageAllowed  int
}

func (r FlightReservationImpl) AddExtraLuggageAllowance(peices int) {
	r.luggageAllowed = peices
}

func (r FlightReservationImpl) CalculateCancellationFee() float64 {
	return 2.0
}

func (r FlightReservationImpl) GetReservationDate() string {
	return r.reservationDate
}

type Trip struct {
	reservations []Reservation
}

func (t *Trip) CalculateCancellationFee() float64 {
	total := 0.0

	for _, r := range t.reservations {
		total += r.CalculateCancellationFee()
	}

	return total
}

func (t *Trip) AddReservation(r Reservation) {
	t.reservations = append(t.reservations, r)
}

func main() {
	var (
		h HotelReservationImpl
		f FlightReservationImpl
		t Trip
	)

	fmt.Println(f.CalculateCancellationFee())
	fmt.Println(h.CalculateCancellationFee())

	t.AddReservation(h)
	t.AddReservation(f)
	fmt.Println(t.CalculateCancellationFee())
}

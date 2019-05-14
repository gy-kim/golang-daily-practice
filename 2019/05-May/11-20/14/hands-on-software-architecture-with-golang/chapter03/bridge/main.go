// this file demonstrates the bridge design pattern
package main

import "fmt"

// Reservation is the main abstraction
// the abstraction here is a struct not an interface, since in Go you can have abstract structs/interfaces ,
// where one can store reference to the Seller implementation
type Reservation struct {
	sellerRef Seller
}

func (r Reservation) Cancel() {
	// charge $10 as cancellation feed
	r.sellerRef.CancelReservation(10)
}

type PremiumReservation struct {
	Reservation
}

func (r PremiumReservation) Cancel() {
	r.sellerRef.CancelReservation(0)
}

// This is the interface for all Sellers
type Seller interface {
	CancelReservation(charge float64)
}

type InstitutionSeller struct{}

func (s InstitutionSeller) CancelReservation(charge float64) {
	fmt.Println("InstitutionSeller CancelReservation charge =", charge)
}

type SmallScaleSeller struct{}

func (s SmallScaleSeller) CancelReservation(charge float64) {
	fmt.Println("SmallScalSeller CancelReservation charge =", charge)
}

func main() {
	res := Reservation{InstitutionSeller{}}
	res.Cancel()

	premiumRes := PremiumReservation{Reservation{SmallScaleSeller{}}}
	premiumRes.Cancel()
}

package main

import "fmt"

type HotelBoutiqueProxy struct {
	subject *HotelBoutique
}

func (p *HotelBoutiqueProxy) Book() {
	if p.subject == nil {
		p.subject = new(HotelBoutique)
	}
	fmt.Println("Proxy Delegating Booking call")

	p.subject.Book()
}

type HotelBoutique struct{}

func (s *HotelBoutique) Book() {
	fmt.Println("Booking done on external site.")
}

func main() {}

package wasted_effort

// Excessive comments
func outputOrderedPeopleA(in []*Person) {
	// This code orders people by name.
	// In cases where the name is the same, it will order by phone number.
	// The sort algorithm used is a bubble sort
	// WARNING: this sort will change the items of the inpt array
	for range in {
		// ... sort code removed ...
	}
}

// Comments replced with descriptive names
func outputOrderedPeopleB(in []*Person) {
	sortPeople(in)
	outputPeople(in)
}

func outputPeople(in []*Person) {
	// TODO: implement
}

// any spedcial instructions that MUST be decuomented relating to the sort should go here
func sortPeople(in []*Person) {
	// TODO: implement
}

// Person data object
type Person struct {
	Name  string
	Phone string
}

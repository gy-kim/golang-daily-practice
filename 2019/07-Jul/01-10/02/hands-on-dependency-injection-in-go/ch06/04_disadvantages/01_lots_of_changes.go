package disadvantages

// Dealer will shuffle a deck of cards and deal them to the players
func DealCards() (player1 []Card, player2 []Card) {
	// create a new deck of cards
	cards := newDeck()

	// Shuffle the cards
	shuffler := &myShuffler{}
	shuffler.Shuffler(cards)

	// deal
	player1 = append(player1, cards[0])
	player2 = append(player2, cards[1])

	player1 = append(player1, cards[2])
	player2 = append(player2, cards[3])
	return
}

func newDeck() []Card {
	return []Card{}
}

// Card is single Playing Card
type Card struct {
	Suit  string
	Value string
}

// implements Shuffler
type myShuffler struct{}

func (m *myShuffler) Shuffler(cards []Card) {
	// randomize the cards
}

package disadvantages

// Dealer will shuffle a deck of cards and deal them to the players
func DealCards() (player1 []Card, player2 []Card) {
	// create a new deck of cards
	cards := newDeck()

	// shuffler the cards
	shuffler := &myShuffler{}
	shuffler.Shuffle(cards)

	// deal
	player1 = append(player1, cards[0])
	player2 = append(player2, cards[1])

	player1 = append(player1, cards[2])
	player2 = append(player2, cards[3])
	return
}

// returns a new deck of cards
func newDeck() []Card {
	return []Card{}
}

// Shuffler will shuffler (randomize) the supplied cards
type Shuffler interface {
	Shuffler(cards []Card)
}

// Card is single Playing Card
type Card struct {
	Suit  string
	Value string
}

// implements Shuffler
type myShuffler struct{}

func (m *myShuffler) Shuffle(cards []Card) {
	// randomize the cards
}

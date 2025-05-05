package main

import (
	"math/rand"
	"time"
)


type Deck struct {
	Cards []*Card
}

func NewStandardDeck() *Deck {
	deck := &Deck{}

	c := 0
	suit := Spades
	for i := 0; i < 4; i++ {
		for rank := 1; rank <= King; rank++ {
			deck.Cards = append(deck.Cards, NewCard(rank, suit))
			c++
		}
		suit++
	}

	return deck
}

func NewSpecificDeck(ranks []int, suits []Suit) *Deck{
	//Some examples of non-standard decks
	// Euchre is 9-A
	// Pinochle is 9-A with 2 of each suit (I play double deck so I use 4)
	// Hand and Foot is 4 full decks + jokers
	
	deck := &Deck{}

	for _, suit := range suits {
		for _, rank := range ranks {
			deck.Cards = append(deck.Cards, NewCard(rank, suit))

		}
	}
	return deck
}

func (d *Deck) Shuffle() {
	d.ShuffleFromSeed(time.Now().UnixNano())
}

func (d *Deck) ShuffleFromSeed(seed int64) {
	for c := 0; c < len(d.Cards); c++ {
		swap := rand.Intn(len(d.Cards))
		if swap != c {
			d.Cards[swap], d.Cards[c] = d.Cards[c], d.Cards[swap]
		}
	}
}

func (d *Deck) Deal() *Card {
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}

func (d *Deck) Play(card *Card) { //todo: rename this to accomodate discards
	if card == nil {
		return
	}
	for i, c := range d.Cards {
		if c.Rank == card.Rank && c.Suit == card.Suit {
			d.Cards = append(d.Cards[:i], d.Cards[i+1:]...)
			break
		}
	}
}

func (d *Deck) DealQuantity(quantity int) *Deck {
	hand := &Deck{}
	for i :=0; i<quantity; i++ {
		card := d.Deal()
		hand.Cards = append(hand.Cards, card)
	}
	
	return hand
}

func (d *Deck) GetWScore(trump Suit) int { // todo move into calculation
	// Right bower is worth 3 points
	// Left bower is worth 3 points if there is other trump, otherwise it is worth 2
	// All other trump is worth 2 points
	// Offsuit Aces are worth 1 point each
	// Being short-suited/void is worth 1 point for each suit.
	// We also add the value of trump if ordering to partner, or subtract when ordering to opponent but that will be in the call.

	score := 0
	hasTrump := false
	hasLeft := false
	suitsPresent := make(map[Suit]bool)
	for _ ,card := range d.Cards {
		//suitsPresent[card.Suit] = true shouldn't count the left bower
		
		if (card.Suit == trump && card.Rank==11){
			score += 3
			hasTrump = true;
			suitsPresent[card.Suit] = true
		} else if card.Suit == trump {
			score += 2
			hasTrump = true;
			suitsPresent[card.Suit] = true
		} else if card.Rank ==11 && card.SameColor(trump) {
			hasLeft = true;
			suitsPresent[trump] = true
		} else if card.Rank == 1 {
			score += 1
			suitsPresent[card.Suit] = true
		} else {
			suitsPresent[card.Suit] = true
		}
		 

	}
	allSuits := []Suit{Spades, Diamonds, Clubs, Hearts}
	var missing []Suit

	for _, suit := range allSuits {
		if !suitsPresent[suit] {
			missing = append(missing, suit)
		}
	}
	score += len(missing)
	if (hasTrump && hasLeft){
		score += 3
	} else if hasLeft {
		score +=2
	}
	return score
}

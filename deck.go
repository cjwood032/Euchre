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

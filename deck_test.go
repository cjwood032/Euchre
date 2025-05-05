package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeckSize(t *testing.T) {
	d := NewStandardDeck()
	assert.Equal(t, 52, len(d.Cards))
}

func TestNewSpecificDeck(t *testing.T) {

	ranks := []int{1,9,10,11,12,13}
	suits := []Suit{Spades,Diamonds,Clubs,Hearts,Spades,Diamonds,Clubs,Hearts}	
	expected := len(ranks) * len(suits)
	d:= NewSpecificDeck(ranks, suits)
	assert.Equal(t,expected,len(d.Cards))

}

func TestShuffle(t *testing.T) {
	unshuffledDeck := NewStandardDeck()
	shuffledDeck := NewStandardDeck()
	assert.Equal(t, *unshuffledDeck.Cards[0], *shuffledDeck.Cards[0])
	shuffledDeck.Shuffle()
	assert.NotEqual(t, *unshuffledDeck.Cards[0], *shuffledDeck.Cards[0])
}

func TestDeal(t *testing.T) {
	deck := NewStandardDeck()
	topCard := deck.Cards[0]
	card := deck.Deal()
	assert.Equal(t, 51, len(deck.Cards))
	assert.Equal(t, *topCard, *card)
}

func TestDealQuantity(t *testing.T) {
	cardsToDeal := 3
	deck := NewStandardDeck()
	unchangedDeck := NewStandardDeck()
	cards := deck.DealQuantity(cardsToDeal).Cards
	assert.Equal(t, len(unchangedDeck.Cards) - cardsToDeal, len(deck.Cards))
	assert.Equal(t, cardsToDeal, len(cards))

	for i := 0; i < len(cards); i++ {
		card := cards[i]
		expectedCard := unchangedDeck.Cards[i]
		assert.Equal(t, *expectedCard, *card)
	}
}
func TestPlayCardFromDeck(t *testing.T) {
	deck := NewStandardDeck()
	assert.Equal(t, 52, len(deck.Cards))
	card := NewCard(10,Diamonds)
	deck.Play(card)
	assert.NotContains(t,deck.Cards,card)
	assert.Equal(t, 51, len(deck.Cards))
}

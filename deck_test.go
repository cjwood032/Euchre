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
// todo move this to calculations when we introduce it.

func TestGetWScore(t *testing.T) {
	trump := Diamonds
	Card1 := NewCard(11,Diamonds)
	Card2 := NewCard(1,Clubs)
	Card3 := NewCard(9,Clubs)
	Card4 := NewCard(13,Spades)
	Card5 := NewCard(12,Hearts)
	newDeck := &Deck{Cards: []*Card{Card1, Card2, Card3, Card4, Card5}}
	wScore := newDeck.GetWScore(trump);
	assert.Equal(t,4,wScore)
}
func TestGetWScoreWithLeft(t *testing.T) {
	trump := Diamonds
	Card1 := NewCard(11,Diamonds)
	Card2 := NewCard(1,Clubs)
	Card3 := NewCard(9,Hearts)
	Card4 := NewCard(13,Spades)
	Card5 := NewCard(11,Hearts)
	newDeck := &Deck{Cards: []*Card{Card1, Card2, Card3, Card4, Card5}}
	wScore := newDeck.GetWScore(trump);
	assert.Equal(t,7,wScore)
	newDeck.Cards[0] = NewCard(10, Hearts)
	wScore = newDeck.GetWScore(trump);
	assert.Equal(t,3,wScore) // - 3 for no right, -1 for left losing a point
}
func TestGetWScoreWithVoidSuits(t *testing.T) {
	trump := Diamonds
	Card1 := NewCard(11,Diamonds)
	Card2 := NewCard(1,Clubs)
	Card3 := NewCard(9,Hearts)
	Card4 := NewCard(13,Spades)
	Card5 := NewCard(10,Diamonds)
	newDeck := &Deck{Cards: []*Card{Card1, Card2, Card3, Card4, Card5}}
	wScore := newDeck.GetWScore(trump);
	assert.Equal(t,6,wScore)
	newDeck.Cards[2] = NewCard(10, Clubs)
	wScore = newDeck.GetWScore(trump);
	assert.Equal(t,7,wScore) 
}
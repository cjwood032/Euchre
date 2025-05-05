package main

import (
	"testing"
)


func AssertSize(t *testing.T,expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected deck length of %v, but got %v", expected, actual)
	}
}
func AssertCardsMatch(t *testing.T,expected Card, actual Card) {
	if expected != actual {
		t.Errorf("Expected %v of %s, but got %v of %s", expected.Rank, expected.Suit.FriendlySuit(), actual.Rank, actual.Suit.FriendlySuit())
	}
}
func AssertCardsDoNotMatch(t *testing.T,expected Card, actual Card) {
	if expected == actual {
		t.Errorf("Expected cards to be different, but got %v of %s for both", expected.Rank, expected.Suit.FriendlySuit())
	}
}
func TestNewDeckSize(t *testing.T) {
	d := NewStandardDeck()
	AssertSize(t, 52, len(d.Cards))
}

func TestNewSpecificDeck(t *testing.T) {

	ranks := []int{1,9,10,11,12,13}
	suits := []Suit{Spades,Diamonds,Clubs,Hearts,Spades,Diamonds,Clubs,Hearts}	
	expected := len(ranks) * len(suits)
	d:= NewSpecificDeck(ranks, suits)
	AssertSize(t,expected,len(d.Cards))

}

func TestShuffle(t *testing.T) {
	unshuffledDeck := NewStandardDeck()
	shuffledDeck := NewStandardDeck()
	AssertCardsMatch(t, *unshuffledDeck.Cards[0], *shuffledDeck.Cards[0])
	shuffledDeck.Shuffle()
	AssertCardsDoNotMatch(t, *unshuffledDeck.Cards[0], *shuffledDeck.Cards[0])
}

func TestDeal(t *testing.T) {
	deck := NewStandardDeck()
	topCard := deck.Cards[0]
	card := deck.Deal()
	AssertSize(t, 51, len(deck.Cards))
	AssertCardsMatch(t, *topCard, *card)
}

func TestDealQuantity(t *testing.T) {
	cardsToDeal := 3
	deck := NewStandardDeck()
	unchangedDeck := NewStandardDeck()
	cards := deck.DealQuantity(cardsToDeal).Cards
	AssertSize(t, len(unchangedDeck.Cards) - cardsToDeal, len(deck.Cards))
	AssertSize(t, cardsToDeal, len(cards))

	for i := 0; i < len(cards); i++ {
		card := cards[i]
		expectedCard := unchangedDeck.Cards[i]
		AssertCardsMatch(t, *expectedCard, *card)
	}
}
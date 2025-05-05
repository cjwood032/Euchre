package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCard(t *testing.T) {
	expectedSuit := Spades
	expectedRank := 5
	c := NewCard(5,expectedSuit)
	
	assert.Equal(t, expectedSuit, c.Suit, "Expected matching suits")
	assert.Equal(t, expectedRank, c.Rank, "Expected matching ranks")
}
func TestFaceCardRank(t *testing.T) {
	expectedSuit := Spades
	expectedRank := 13
	c := NewCard(13,expectedSuit)
	
	assert.Equal(t, expectedSuit, c.Suit, "Expected matching suits")
	assert.Equal(t, expectedRank, c.Rank, "Expected matching ranks")
}

func TestColor(t *testing.T) {
	suits := []Suit{Spades,Clubs,Hearts,Diamonds}
	for _, suit := range suits {
		expectedColor := SuitColorBlack
		if (suit == Hearts || suit == Diamonds) {
			expectedColor = SuitColorRed
		}
		c := NewCard(1,suit)
		actualColor := c.Color()
		assert.Equal(t, expectedColor, actualColor, "Expected matching card colors")
	}

}
func TestFriendlyColor(t *testing.T) {
	suits := []Suit{Spades,Clubs,Hearts,Diamonds}
	suitNames := []string{"Spades","Clubs","Hearts","Diamonds"}
	for i, suit := range suits {
		expectedSuit := suitNames[i]
		actualSuit := suit.FriendlySuit()

		assert.Equal(t, expectedSuit, actualSuit, "Expected suit to be the string")
	}

}
func TestCardFaces(t *testing.T) {
	card := NewCard(3, Clubs)
	assert.False(t,card.FaceUp, "Expected a new card to be face down")
	card.TurnFaceUp()
	assert.True(t,card.FaceUp, "Expected card that is turned face up to be face up")
	card.TurnFaceDown()
	assert.False(t,card.FaceUp, "Expected card that is turned face down to be face down")
}

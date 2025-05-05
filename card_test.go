package main

import (
	"testing"
)

func TestNewCard(t *testing.T) {
	expectedSuit := Spades
	expectedRank := 5
	c := NewCard(5,expectedSuit)
	
	if c.Suit != expectedSuit {
		t.Errorf("Suit did not match, expected %v got %v", expectedSuit, c.Suit)
	}
	if c.Rank != expectedRank {
		t.Errorf("Card value did not match, expected %v got %v", expectedRank, c.Rank)
	}
}
func TestFaceCardRank(t *testing.T) {
	expectedSuit := Spades
	expectedRank := 13
	c := NewCard(13,expectedSuit)
	
	if c.Suit != expectedSuit {
		t.Errorf("Suit did not match, expected %v got %v", expectedSuit, c.Suit)
	}
	if c.Rank != expectedRank {
		t.Errorf("Card value did not match, expected %v got %v", expectedRank, c.Rank)
	}
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
		if expectedColor != actualColor {
			t.Errorf("Expected %s to be color %s, got %s",suit.FriendlySuit(), expectedColor.FriendlySuitColor(), actualColor.FriendlySuitColor())
		}
	}

}
func TestFriendlyColor(t *testing.T) {
	suits := []Suit{Spades,Clubs,Hearts,Diamonds}
	suitNames := []string{"Spades","Clubs","Hearts","Diamonds"}
	for i, suit := range suits {
		expectedSuit := suitNames[i]
		actualSuit := suit.FriendlySuit()

		
		if expectedSuit != actualSuit {
			t.Errorf("Expected %s, but got %s", expectedSuit, actualSuit)
		}
	}

}
func TestCardFaces(t *testing.T) {
	card := NewCard(3, Clubs)
	if card.FaceUp {
		t.Error("Expected newly created cards to be face down")
	}
	card.TurnFaceUp()
	if !card.FaceUp {
		t.Error("Expected card that is turned face up to be face up")
	}
	card.TurnFaceDown()
	if card.FaceUp {
		t.Error("Expected card that is turned face down to be face down")
	}
}

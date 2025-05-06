package main

import (
	"log"
)

type Suit int

type SuitColor int

const (
	Spades Suit = iota
	Diamonds
	Clubs
	Hearts

	// The colors of the suits for the color call
	SuitColorBlack SuitColor = iota
	SuitColorRed
)

const (
	Jack = 11
	Queen = 12
	King = 13
	//The other ranks will be the card rank. Aces being high or low will be determined by the game
)

type Card struct {
	Rank int
	Suit Suit
	FaceUp bool
	
}


func NewCard(rank int, suit Suit) *Card {
	if (rank < 1 || rank > 13) {
		log.Fatal("Invalid card rank")
	}
	return &Card{Rank: rank, Suit: suit}
}

func (c *Card) Color() SuitColor {
	if c.Suit == Clubs || c.Suit == Spades {
		return SuitColorBlack
	}

	return SuitColorRed
}
func (c *Card) SameColor(trump Suit) bool {
	if ((c.Suit == Clubs && trump == Spades) || 
		(c.Suit == Spades && trump == Clubs) ||
		(c.Suit == Hearts && trump == Diamonds) ||
		(c.Suit == Diamonds && trump == Hearts)) {
		return true
	}

	return false
}
func (suit *Suit) SameColor(trump Suit) bool {
	if ((*suit == Clubs && trump == Spades) || 
		(*suit == Spades && trump == Clubs) ||
		(*suit == Hearts && trump == Diamonds) ||
		(*suit == Diamonds && trump == Hearts)) {
		return true
	}

	return false
}
func (c *Card) TurnFaceUp() {
	c.FaceUp = true
}

// TurnFaceDown sets the FaceUp field to false - so the card should be hidden
func (c *Card) TurnFaceDown() {
	c.FaceUp = false
}

func (suit Suit) FriendlySuit() string {
	switch suit {
	case Spades:
		return "Spades"
	case Diamonds:
		return "Diamonds"
	case Clubs:
		return "Clubs"
	case Hearts:
		return "Hearts"
	default:
		return "Unknown"
	}
}
func (color SuitColor) FriendlySuitColor() string {
	switch color {
	case SuitColorBlack:
		return "Black"
	case SuitColorRed:
		return "Red"
	default:
		return "Unknown"
	}
}

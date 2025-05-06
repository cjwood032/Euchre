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



func (cm *CardMap) GetWScore(trump Suit) int {
	// Right bower is worth 3 points
	// Left bower is worth 3 points if there is other trump, otherwise it is worth 2
	// All other trump is worth 2 points
	// Offsuit Aces are worth 1 point each
	// Being short-suited/void is worth 1 point for each suit.
	// We also add the value of trump if ordering to partner, or subtract when ordering to opponent but that will be in the call.
	score := 0
	hasTrump := false
	hasLeft := false

	for suit := Spades; suit <= Hearts; suit++ {
		for rank := 1; rank <= 13; rank++ {
			if !cm.Hand[suit][rank] {
				continue
			}

			// Right bower
			if suit == trump && rank == 11 {
				score += 3
				hasTrump = true
			} else if suit.SameColor(trump) && suit != trump && rank == 11 {
				// Left bower (will count as trump)
				hasLeft = true
			} else if suit == trump {
				score += 2
				hasTrump = true
			} else if rank == 1 {
				// Offsuit ace
				score += 1
			}
		}
	}

	// Add bonus for left bower based on whether we have other trump
	if hasLeft && hasTrump {
		score += 3
	} else if hasLeft {
		score += 2
	}

	// Add points for void suits
	suitCounts := cm.CountSuits(trump)
	for _, count := range suitCounts {
		if count == 0 {
			score += 1
		}
	}

	return score
}


func (cm *CardMap) BestTrumpScore(excludedSuit Suit) (bestSuit Suit, bestScore int) {
	allSuits := []Suit{Spades, Diamonds, Clubs, Hearts}
	bestScore = -1 // initialize lower than possible score

	for _, suit := range allSuits {
		if suit == excludedSuit {
			continue
		}
		score := cm.GetWScore(suit)
		
		if score > bestScore {
			bestScore = score
			bestSuit = suit
		}
	}

	return bestSuit, bestScore
}


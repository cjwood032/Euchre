package main

type CardMap struct {
	Hand [4][14]bool // [Suit][Rank] - cards in player's hand
	Seen [4][14]bool // [Suit][Rank] - cards the player has seen
}


func (cm *CardMap) AddToHand(card *Card) {
	cm.Hand[card.Suit][card.Rank] = true
}

func (cm *CardMap) RemoveFromHand(card *Card) {
	cm.Hand[card.Suit][card.Rank] = false
	cm.Seen[card.Suit][card.Rank] = true // Also mark as seen
}
func (cm *CardMap) MarkSeen(card *Card) {
	cm.Seen[card.Suit][card.Rank] = true
}

func (cm *CardMap) HasInHand(card *Card) bool {
	return cm.Hand[card.Suit][card.Rank]
}

func (cm *CardMap) HasSeen(card *Card) bool {
	return cm.Seen[card.Suit][card.Rank]
}

func (cm *CardMap) CountSuits(trump Suit) map[Suit]int {
	counts := make(map[Suit]int)
	allSuits := []Suit{Spades, Diamonds, Clubs, Hearts}

	for _, suit := range allSuits {
		for rank := 1; rank <= 13; rank++ {
			if cm.Hand[suit][rank] {
				actualSuit := suit
				if rank == 11 && suit.SameColor(trump) && suit != trump {
					// Left bower counts as trump
					actualSuit = trump
				}
				counts[actualSuit]++
			}
		}
	}

	// Ensure all suits are represented
	for _, suit := range allSuits {
		if _, ok := counts[suit]; !ok {
			counts[suit] = 0
		}
	}

	return counts
}


func (cm *CardMap) isLeftBower(cardSuit, trump Suit) bool {
	return cardSuit != trump &&
		cardSuit.SameColor(trump) &&
		cm.Hand[cardSuit][11]
}

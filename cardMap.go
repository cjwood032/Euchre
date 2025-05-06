package main

type CardMap struct {
	Hand [4][14]bool // [Suit][Rank] - cards in player's hand
	Seen [4][14]bool // [Suit][Rank] - cards the player has seen
}


func (cm *CardMap) AddToHand(card *Card) {
	cm.Hand[card.Suit][card.Rank] = true
}

func (cm *CardMap) AddCardsToHand(cards *Deck) {
	for _, card := range cards.Cards {
		cm.AddToHand(card)
	}
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


func (cm *CardMap) hasLeftBower(trump Suit) bool {

	cardSuit := trump.GetWeakColor()
	return cm.Hand[cardSuit][11]
}

func (cm CardMap) ToSlice() []*Card {
	var cards []*Card
	for suit := 0; suit < 4; suit++ {
		for rank := 0; rank < 14; rank++ {
			if cm.Hand[suit][rank] {
				cards = append(cards, &Card{Suit: Suit(suit), Rank: rank})
			}
		}
	}
	return cards
}

func (cm CardMap) CardsInSuit(suit Suit) []*Card {
	var cards []*Card
	for rank := 0; rank < 14; rank++ {
		if cm.Hand[suit][rank] {
			cards = append(cards, &Card{Suit: suit, Rank: rank})
		}
	}
	return cards
}
func (cm CardMap) CountSuit(suit Suit) int {
	count := 0
	for rank := 0; rank < 14; rank++ {
		if cm.Hand[suit][rank] {
			count++
		}
	}
	return count
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

func (cm *CardMap) Sort(suit Suit, isTrump bool) []*Card {
	
	hasLeft := cm.hasLeftBower(suit)
	hasRight := false
	hasAce :=false
	var cards []*Card
	for rank := 0; rank < 14; rank++ {
		if cm.Hand[suit][rank] {
			if rank == 1 {
				hasAce = true;
			}else if rank == 11{
				hasRight = true;
			} else {
				cards = append(cards, &Card{Suit: suit, Rank: rank})
			}
		}
	}
	if hasAce {
		cards = append(cards, &Card{Suit: suit, Rank: 1})
	}
	if hasLeft && isTrump {
		cards = append(cards, &Card{Suit: suit.GetWeakColor(), Rank: 11})
	}
	if hasRight && isTrump {
		cards = append(cards, &Card{Suit: suit, Rank: 11})
	}
	
	return cards
}

func (cm *CardMap) getStrongestOffsuit(trump Suit) *Card {
	oppositeColorSuits := trump.GetOppositeColors()

	// 1. Prefer short-suited opposite-color suits (only one card)
	for _, suit := range oppositeColorSuits {
		if cm.CountSuit(suit) == 1 {
			for rank := 13; rank >= 0; rank-- {
				if cm.Hand[suit][rank] {
					return &Card{Suit: suit, Rank: rank}
				}
			}
		}
	}

	// 2. Otherwise, pick the strongest card among opposite-color suits
	for _, suit := range oppositeColorSuits {
		for rank := 13; rank >= 0; rank-- {
			if cm.Hand[suit][rank] {
				return &Card{Suit: suit, Rank: rank}
			}
		}
	}

	// 3. Fallback: strongest non-trump card
	for suit := 0; suit < 4; suit++ {
		if Suit(suit) == trump {
			continue
		}
		for rank := 13; rank >= 0; rank-- {
			if cm.Hand[suit][rank] {
				return &Card{Suit: Suit(suit), Rank: rank}
			}
		}
	}

	return nil
}

	

package main


type Call int
const (
	Pass Call = iota
	OrderUp
	Alone
)

type Player struct {
	Name string
	Hand *Deck
	CardMap CardMap
	CardsInSuit map[Suit]int
	Score int
	Wins int
	Losses int
	ComputerPlayer bool
}


func (player *Player)PlayCard(card *Card) *Card {
	player.Hand.Play(card)
	return card
}

var minimumScore = 7
var lonerScore = 10

func (player *Player)CallOrPass(trump Suit, teamPickup bool ) Call {
	wScore := player.CardMap.GetWScore(trump)
	//todo, sit if your to the left of the dealer and you're stronger in next
	if teamPickup {
		wScore += 2
	} else {
		wScore -= 2
	}
	return DetermineCall(wScore)
}

func (player *Player)DeclareTrump(unavailableSuit Suit) (Call, Suit) {
	suit, score := player.CardMap.BestTrumpScore(unavailableSuit)
	return DetermineCall(score), suit
}

func (player *Player)PickUp(card *Card) *Card {
	
	player.Hand.Cards = append(player.Hand.Cards, card)
	return player.DiscardCard()
}


func (player *Player)DiscardCard() *Card {
	card := player.Hand.Cards[0] // todo Capture discarded card
	player.Hand.Play(card)
	return card
}

func DetermineCall(score int) Call {
	if score >= lonerScore {
		return OrderUp // todo declare a loner
	} else if score >= minimumScore {
		return OrderUp
	}
	return Pass
}

	// Determine who still has to play from currentCards
	// who is the team that called
	// who is winning the trick?
	// if player's team is winning how strong is the winning card? 
	// if player has the suit still
	// if player's team is not winning, or the winning card is weak Q or less, play strongest card that can win
	// if player's team is winning, play lowest card if
	// if player does not have the suit 
	// if player's team is not winning, play trump to win
	// if player's team is winning, play to short suit if player has only one card in another non-trump suit, otherwise throw low non-trump

	func (player *Player) BestPlay(currentTrick []Card, round Round) Card {
		if len(currentTrick) ==0 {
			return *player.PlayCard(player.Hand.Cards[0])//BestLead()
		}
		leadSuit := currentTrick[0].Suit
		winningCard, winningPlayer := getWinningCard(currentTrick, round.Trump, leadSuit)
		winningTeam := (winningPlayer == round.Caller) || (player.getPartner(round.Players) == winningPlayer)
	
		hand := player.Hand.Cards
		playable := getPlayableCards(hand, leadSuit, round.Trump)
		hasLeadSuit := len(playable.inSuit) > 0
	
		if !hasLeadSuit {
			if !winningTeam {
				if len(playable.trump) > 0 {
					betterTrump := getLowestWinningTrump(playable.trump, winningCard, round.Trump, leadSuit)
					if betterTrump != nil {
						return *betterTrump
					} 
				}
				return getLowest(playable.other, round.Trump)
			} else {
				shortSuit := findShortSuit(player.Hand, round.Trump)
				if shortSuit != -1 {
					return getCardInSuit(player.Hand, shortSuit, true)
				}
				return getLowest(playable.other, round.Trump)
			}
		} else {
			if !winningTeam || isWeak(winningCard) {
				winning := getStrongerThan(playable.inSuit, winningCard, round.Trump)
				if len(winning) > 0 {
					return getStrongest(winning, round.Trump)
				}
				return getLowest(playable.inSuit, round.Trump)
			} else {
				return getLowest(playable.inSuit, round.Trump)
			}
		}
	}
	
	func getLowestWinningTrump(cards []*Card, currentWinner Card, trump Suit, lead Suit) *Card {
		var winningTrumps []*Card
		for _, c := range cards {
			if c.Beats(currentWinner, trump, lead) {
				winningTrumps = append(winningTrumps, c)
			}
		}
		if len(winningTrumps) == 0 {
			return nil
		}
		// Return the lowest trump that still wins
		lowest := winningTrumps[0]
		for _, c := range winningTrumps[1:] {
			if !lowest.Beats(*c, trump, lead) {
				lowest = c
			}
		}
		return lowest
	}
	
	
	func getWinningCard(cards []Card, trump Suit, lead Suit) (Card, *Player) {
		winning := cards[0]
		for _, card := range cards[1:] {
			if card.Beats(winning, trump, lead) {
				winning = card
			}
		}
		return winning, nil // Placeholder
	}
	
	func getPlayableCards(hand []*Card, lead Suit, trump Suit) (result struct{ inSuit, trump, other []*Card }) {
		for _, c := range hand {
			if c.Suit == lead || (c.Rank == 11 && c.SameColor(trump) && c.Suit != trump && lead == trump) {
				result.inSuit = append(result.inSuit, c)
			} else if c.Suit == trump || (c.Rank == 11 && c.SameColor(trump) && c.Suit != trump) {
				result.trump = append(result.trump, c)
			} else {
				result.other = append(result.other, c)
			}
		}
		return
	}
	
	func getStrongest(cards []*Card, trump Suit) Card {
		strongest := cards[0]
		for _, c := range cards[1:] {
			if c.Beats(*strongest, trump, strongest.Suit) {
				strongest = c
			}
		}
		return *strongest
	}
	
	func getLowest(cards []*Card, trump Suit) Card {
		lowest := cards[0]
		for _, c := range cards[1:] {
			if !c.Beats(*lowest, trump, c.Suit) {
				lowest = c
			}
		}
		return *lowest
	}
	
	func getStrongerThan(cards []*Card, target Card, trump Suit) []*Card {
		var result []*Card
		for _, c := range cards {
			if c.Beats(target, trump, target.Suit) {
				result = append(result, c)
			}
		}
		return result
	}
	
	func isWeak(card Card) bool {
		return card.Rank <= 12 // Assume Q or lower is weak
	}
	
	func findShortSuit(hand *Deck, trump Suit) Suit {
		suitCounts := make(map[Suit]int)
		for _, c := range hand.Cards {
			if c.Suit != trump {
				suitCounts[c.Suit]++
			}
		}
		for suit, count := range suitCounts {
			if count == 1 {
				return suit
			}
		}
		return -1
	}
	
	func getCardInSuit(hand *Deck, suit Suit, lowest bool) Card {
		var candidates []Card
		for _, c := range hand.Cards {
			if c.Suit == suit {
				candidates = append(candidates, *c)
			}
		}
		if len(candidates) == 0 {
			return *hand.Cards[0] // fallback
		}
		if lowest {
			return getLowest(cardPointers(candidates), suit)
		}
		return getStrongest(cardPointers(candidates), suit)
	}
	
	func (player *Player) getPartner(players []*Player) *Player {
		for i, p := range players {
			if p == player {
				if i > 1 {

					return players[i-2]
				}
				return players[i+2]
			}
		}
		return nil
	}
	
	func cardPointers(cards []Card) []*Card {
		var ptrs []*Card
		for i := range cards {
			ptrs = append(ptrs, &cards[i])
		}
		return ptrs
	}
	


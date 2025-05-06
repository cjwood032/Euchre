package main


type Call int
const (
	Pass Call = iota
	OrderUp
	Alone
)

type Player struct {
	Name string
	CardMap CardMap
	CardsInSuit map[Suit]int
	Score int
	Wins int
	Losses int
	ComputerPlayer bool
	TricksWon int
}

var minimumScore = 7
var lonerScore = 10

func (player *Player) PlayCard(card *Card) *Card {
	player.CardMap.RemoveFromHand(card)
	return card
}

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

func (player *Player) PickUp(card *Card) *Card {
	player.CardMap.AddToHand(card)
	return player.DiscardCard(card)
}

func (player *Player) InitCardMap() {
    player.CardMap = CardMap{
        Hand: [4][14]bool{},  // Clears all cards from hand
        Seen: [4][14]bool{},  // Clears all seen cards
    }
    player.CardsInSuit = make(map[Suit]int)
    player.TricksWon = 0
}

func (player *Player) DiscardCard(card *Card) *Card {
		player.CardMap.RemoveFromHand(card)
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
func (player *Player) BestPlay(currentTrick []*Card, round Round) Card {
	if len(currentTrick) == 0 {
		//we lead
		trumpCards := player.CardMap.CardsInSuit(round.Trump)
		if player.getPartner(round.Players) == round.Caller && len(trumpCards) > 0 {
			return *player.CardMap.Sort(round.Trump, true)[0]
		}
		return *player.CardMap.getStrongestOffsuit(round.Trump)
	}
	leadSuit := currentTrick[0].Suit
	winningCard, winningPlayer := getWinningCard(currentTrick, round.Players, round.Trump, leadSuit)
	winningTeam := player.getPartner(round.Players) == winningPlayer

	hand := player.CardMap.ToSlice()
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
			shortSuit := findShortSuit(player.CardMap, round.Trump)
			if shortSuit != -1 {
				return getCardInSuit(player.CardMap, shortSuit, true)
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

func getLowestWinningTrump(cards []*Card, currentWinner *Card, trump Suit, lead Suit) *Card {
	var winningTrumps []*Card
	for _, c := range cards {
		if c.Beats(currentWinner, trump, lead) {
			winningTrumps = append(winningTrumps, c)
		}
	}
	if len(winningTrumps) == 0 {
		return nil
	}
	lowest := winningTrumps[0]
	for _, c := range winningTrumps[1:] {
		if !lowest.Beats(c, trump, lead) {
			lowest = c
		}
	}
	return lowest
}

func getWinningCard(cards []*Card, players []*Player, trump Suit, lead Suit) (*Card, *Player) {
	winning := cards[0]
	position := 0
	for i, card := range cards[1:] {
		if card.Beats(winning, trump, lead) {
			winning = card
			position = i + 1
		}
	}
	return winning, players[position]
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
		if c.Beats(strongest, trump, strongest.Suit) {
			strongest = c
		}
	}
	return *strongest
}

func getLowest(cards []*Card, trump Suit) Card {
	lowest := cards[0]
	for _, c := range cards[1:] {
		if !c.Beats(lowest, trump, c.Suit) {
			lowest = c
		}
	}
	return *lowest
}

func getStrongerThan(cards []*Card, target *Card, trump Suit) []*Card {
	var result []*Card
	for _, c := range cards {
		if c.Beats(target, trump, target.Suit) {
			result = append(result, c)
		}
	}
	return result
}

func isWeak(card *Card) bool {
	return card.Rank <= 12
}

func findShortSuit(cardMap CardMap, trump Suit) Suit {
	for suit := Suit(0); suit < 4; suit++ { // Assuming 4 suits: 0 to 3
		if suit != trump && cardMap.CountSuit(suit) == 1 {
			return suit
		}
	}
	return -1
}	

func getCardInSuit(cardMap CardMap, suit Suit, lowest bool) Card {
	cards := cardMap.CardsInSuit(suit)
	if len(cards) == 0 {
		return *cardMap.ToSlice()[0] // fallback
	}
	if lowest {
		return getLowest(cards, suit)
	}
	return getStrongest(cards, suit)
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

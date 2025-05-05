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
	Score int
	Wins int
	Losses int
}

func (player *Player)PlayCard() *Card {
	card := player.Hand.Cards[0]
	player.Hand.Play(card)
	return card
}

var minimumScore = 7
var lonerScore = 10

func (player *Player)CallOrPass(trump Suit, teamPickup bool ) Call {
	wScore := player.Hand.GetWScore(trump)
	//todo, sit if your to the left of the dealer and you're stronger in next
	if teamPickup {
		wScore += 2
	} else {
		wScore -= 2
	}
	return DetermineCall(wScore)
}

func (player *Player)DeclareTrump(unavailableSuit Suit) (Call, Suit) {
	
		suit, score := player.Hand.BestTrumpScore(unavailableSuit)
		
	
	
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
	if score > lonerScore {
		return OrderUp // todo declare a loner
	} else if score > minimumScore {
		return OrderUp
	}
	return Pass
}
//best play
//
// do I need anything else at this point?
package main

type Call int
const (
	Pass Call = iota
	Declare 
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

func (player *Player)CallOrPass(trump Suit) Call {
	return Pass
}

func (player *Player)DeclareTrump(unavailableSuit Suit) (Call, Suit) {
	return Declare, Spades
}

func (player *Player)PickUp(card *Card) *Card {
	player.Hand.Cards = append(player.Hand.Cards, card)
	return player.DiscardCard()
}
func (player *Player)DiscardCard() *Card {
	card := player.Hand.Cards[0] // todo Capture discoarded card
	player.Hand.Play(card)
	return card
}
//best play
//
// do I need anything else at this point?
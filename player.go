package main


type Player struct {
	Name string
	Hand *Deck
	Score int
	Wins int
	Losses int
}

func (player *Player)PlayCard(card *Card) {
	player.Hand.Play(card)
}
// do I need anything else at this point?
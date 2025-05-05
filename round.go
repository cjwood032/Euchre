package main

type Round struct{
	Dealer *Player
	Caller *Player
	Tricks []*Deck
	
}

type PlayerRound struct {
	Player *Player
	TricksWon int
}


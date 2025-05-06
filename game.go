package main

import (
	"math/rand"
)

type Game struct {
	Players []*Player
	Deck *Deck
	Suits []Suit
	Ranks []int
	ScoreLimit int
	Dealer int
	CardsToDeal int
	Rounds []*Round
}

func CreateEuchreGame(players []*Player) *Game {
	game := &Game{
		Players: players,
		ScoreLimit: 10,
		CardsToDeal: 5,
		Ranks : []int{1,9,10,11,12,13},
		Suits : []Suit{Spades,Diamonds,Clubs,Hearts},
	}
	game.Deck = NewSpecificDeck(game.Ranks, game.Suits)
	game.Dealer = rand.Intn(len(game.Players))
	return game
}

func (game *Game) NewGame(changeTeams bool) {
    if changeTeams {
        game.RotateSeats()
    }
    
    // Clear all player states
    for _, player := range game.Players {
        player.InitCardMap()  // This now clears hands
        player.Score = 0
        player.TricksWon = 0
    }
    
    // Reset dealer and start new round
    game.Dealer = rand.Intn(len(game.Players))
    game.NewRound()
}

func (game *Game) NewRound() {
    game.Dealer++
    if game.Dealer == len(game.Players) {
        game.Dealer = 0
    }

    // Initialize player states
    for _, player := range game.Players {
        player.InitCardMap()
        player.TricksWon = 0
    }

    // Create a new deck with the correct cards for Euchre
    round := &Round{
        Players: game.Players,
        Dealer: game.Dealer,
        Deck: NewSpecificDeck(game.Ranks, game.Suits),
    }
    round.Begin()
    game.Rounds = append(game.Rounds, round)
}

func (game *Game) EndRound() {
	if(game.SomeoneWon()){
		game.RecordResults()
		return
	}
	game.NewRound()
}

func (game *Game) ClearScores() {
	for _, player := range game.Players {
		player.Score = 0
	}
}

func (game *Game) RandomizeSeats() {
	for seat := 0; seat < len(game.Players); seat++ {
		swap := rand.Intn(len(game.Players))
		if swap != seat {
			game.Players[swap], game.Players[seat] = game.Players[seat], game.Players[swap]
		}
	}
	game.Dealer = rand.Intn(len(game.Players))
}

func (game *Game) RotateSeats() {
	// this is common in tournaments where 3 players will rotate seats to get new partners
	players := game.Players
	newSeats := make([]*Player, len(players))
	newSeats[0] = players[0]                    // First player stays in place
	newSeats[1] = players[len(players)-1]           // Last player moves to second seat
	copy(newSeats[2:], players[1:len(players)-1])
	
	game.Players = newSeats
}

func (game *Game) SomeoneWon() bool {
	for seat := 0; seat < len(game.Players); seat++ {
		if game.Players[seat].Score >= game.ScoreLimit{
			return true
		}
	}
	return false
}

func (game *Game) RecordResults() {
	for _, player := range game.Players {
		if player.Score >= game.ScoreLimit {
			player.Wins ++
		} else {
			player.Losses++
		}
	}
}

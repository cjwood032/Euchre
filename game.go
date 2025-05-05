package main

import "math/rand"

type Game struct {
	Players []*Player
	Deck *Deck
	Suits []Suit
	Ranks []int
	Kitty *Deck
	ScoreLimit int
	Dealer int
	CardsToDeal int
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

	return game
}

func (game *Game) NewGame(changeTeams bool) {
	if(changeTeams) {
		game.RotateSeats()
	}
	game.ClearScores()
	
}
func (game *Game) ClearScores() {
	for i := 0; i< len(game.Players); i++ {
		game.Players[i].Score = 0
	}
}

func (game *Game) RandomizeSeats() {
	for seat := 0; seat < len(game.Players); seat++ {
		swap := rand.Intn(len(game.Players))
		if swap != seat {
			game.Players[swap], game.Players[seat] = game.Players[seat], game.Players[swap]
		}
	}
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

func (game *Game) EndRound() {
	if(game.SomeoneWon()){
		game.AddResults()
		return
	}
	game.NewRound()
}

func (game *Game) SomeoneWon() bool {
	for seat := 0; seat < len(game.Players); seat++ {
		if game.Players[seat].Score >= game.ScoreLimit{
			return true
		}
	}
	return false
}

func (game *Game) NewRound() {
	game.Dealer++
	if game.Dealer == len(game.Players){
		game.Dealer = 0
	}
	game.Deal()
	//deal
}
func (game *Game) Deal() {
	game.Deck = NewSpecificDeck(game.Ranks, game.Suits) //Rework me to not reinstantiate
	game.Deck.Shuffle()
	// starting with the player to the left of the dealer 
	for seat := 0; seat < len(game.Players); seat++ {
		//game.Players[seat].Hand = &Deck{}
		game.Players[seat].Hand = game.Deck.DealQuantity(5) // deal the appropriate amount
	}
	
}
func (game *Game) AddResults() {
	for i := 0; i < len(game.Players); i++ {
		player := game.Players[i]
		if player.Score >= game.ScoreLimit {
			player.Wins ++
		}else {
			player.Losses++
		}
	}
}
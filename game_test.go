package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var players = []*Player{
	{Name: "Chris"},
	{Name: "Don"},
	{Name: "MaryAnn"},
	{Name: "Andy"},
}
var ranks = []int{1,9,10,11,12,13}
var suits = []Suit{Spades,Diamonds,Clubs,Hearts}
var expectedEuchreDeckSize = 24
var cardsToDeal = 5
func TestCreateEuchreGame(t *testing.T) {
	game := CreateEuchreGame(players)
	assert.Equal(t, expectedEuchreDeckSize, len(ranks) * len(suits))
	assert.Equal(t,expectedEuchreDeckSize,len(game.Deck.Cards))
}

func TestNewGameScoreStartsAtZero(t *testing.T){
	game := CreateEuchreGame(players)
	for i := 0; i < len(game.Players); i++ {
		assert.Equal(t,game.Players[i].Score, 0 , "Expected a score of zero for a new game")
	}
}
func TestNewRoundScoreStartsAtZero(t *testing.T){
	game := CreateEuchreGame(players)
	for i := range game.Players {
		game.Players[i].Score = i+1
		assert.NotEqual(t,game.Players[i].Score, 0 , "Score should increment")
	}
	game.NewGame(false)
	for i := 0; i < len(game.Players); i++ {
		assert.Equal(t,game.Players[i].Score, 0 , "Expected a score of zero for a new game")
	}
}

func TestRotateSeats(t *testing.T) {
	game := CreateEuchreGame(players)
	assert.Equal(t,game.Players[0], players[0],"Expected player in Order")
	assert.Equal(t,game.Players[1], players[1],"Expected player in Order")
	game.RotateSeats()
	assert.Equal(t,game.Players[0], players[0],"Expected first player to stay")
	assert.Equal(t,game.Players[2], players[1],"Expected other players to rotate")
}

func TestGameDealsCorrectly(t *testing.T) {
	game := CreateEuchreGame(players)
	game.Deal()
	
	assert.Equal(t,cardsToDeal,game.CardsToDeal)
	for i := range game.Players {
		assert.Equal(t,len(game.Players[i].Hand.Cards), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	assert.NotEqual(t,game.Players[0].Hand, game.Players[1].Hand)
	totalDealt := (game.CardsToDeal * len(players))
	expectedRemaining := expectedEuchreDeckSize - (cardsToDeal * len(game.Players))
	assert.Equal(t,expectedEuchreDeckSize,totalDealt + len(game.Deck.Cards))
	assert.Equal(t,expectedRemaining, len(game.Deck.Cards))
}
func TestNewRoundDealsCorrectly(t *testing.T) {
	game := CreateEuchreGame(players)
	game.Deal()
	for i := range game.Players {
		player := game.Players[i]
		for range player.Hand.Cards {
			player.Hand.Play(player.Hand.Cards[0])
		}
		assert.Equal(t, 0, len(player.Hand.Cards))
	}
	game.NewRound()
	for i := range game.Players {
		assert.Equal(t,len(game.Players[i].Hand.Cards), cardsToDeal, "Expected players to be dealt 5 cards")
	}
}
func TestEndRoundStopsGameIfScoreMet(t *testing.T) {
	game := CreateEuchreGame(players)
	assert.False(t, game.SomeoneWon())
	game.Players[0].Score = 10
	game.EndRound()
	assert.True(t, game.SomeoneWon())
}
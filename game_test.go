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
	for _, player := range game.Players {
		assert.Equal(t, player.Score, 0 , "Expected a score of zero for a new game")
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

func TestEndRoundStopsGameIfScoreMet(t *testing.T) {
	game := CreateEuchreGame(players)
	assert.False(t, game.SomeoneWon())
	game.Players[0].Score = 10
	game.EndRound()
	assert.True(t, game.SomeoneWon())
}

func TestPlayerWinsLossesIncrementCorrectly(t *testing.T) {
	players = []*Player{
		{Name: "Chris"},
		{Name: "Don"},
		{Name: "MaryAnn"},
		{Name: "Andy"},
	} // we recreate the players otherwise the round stop test adds an extra win/loss
	game := CreateEuchreGame(players)
	game.Players[1].Score = 10
	game.Players[3].Score = 10
	game.EndRound()
	for _, player := range game.Players {
		if player.Score >= 10{
			assert.Equal(t, 1, player.Wins)
			assert.Equal(t, 0, player.Losses)
		} else {
			assert.Equal(t, 0, player.Wins)
			assert.Equal(t, 1, player.Losses)
		}
	}
}

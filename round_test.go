package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstRoundDealsCorrectly(t *testing.T) {
	game := CreateEuchreGame(players)
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	assert.Equal(t,cardsToDeal,game.CardsToDeal)
	for _, player := range round.Players {
		assert.Equal(t,len(player.Hand.Cards), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	assert.NotEqual(t, round.Players[0].Hand, round.Players[1].Hand)
	totalDealt := (game.CardsToDeal * len(players))
	expectedRemaining := expectedEuchreDeckSize - (cardsToDeal * len(round.Players))
	assert.Equal(t,expectedEuchreDeckSize,totalDealt + len(round.Deck.Cards))
	assert.Equal(t,expectedRemaining, len(round.Deck.Cards))
}

func TestNewRoundDealsCorrectly(t *testing.T) {
	game := CreateEuchreGame(players)
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	for _, player := range round.Players {
		for range player.Hand.Cards {
		
			player.Hand.Play(player.Hand.Cards[0])
		}
		assert.Equal(t, 0, len(player.Hand.Cards), "expected empty hands")
	}
	game.NewRound()
	for _, player := range round.Players {
		assert.Equal(t,len(player.Hand.Cards), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	assert.True(t,round.Deck.Cards[0].FaceUp)
}
func TestNewRoundScoreStartsAtZero(t *testing.T){
	game := CreateEuchreGame(players)
	for i, player := range game.Players {
		player.Score = i+1
		assert.NotEqual(t,player.Score, 0 , "Score should increment")
	}
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	for _, player := range round.Players{
		assert.NotEqual(t, player.Score, 0 , "Expected score to persist for new rounds")
	}
}

func TestDeclareTrumpDiscardsCorrectly(t *testing.T) {
	game := CreateEuchreGame(players)
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	assert.Equal(t,cardsToDeal,game.CardsToDeal)
	for _, player := range round.Players {
		
		assert.Equal(t,len(player.Hand.Cards), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	round.DetermineTrump()
	for _, player := range round.Players {
		assert.Equal(t,len(player.Hand.Cards), cardsToDeal, "Expected players to be dealt 5 cards")
	}
}
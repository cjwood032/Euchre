package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstRoundDealsCorrectly(t *testing.T) {
	players := CreatePlayers()
	game := CreateEuchreGame(players)
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	assert.Equal(t,cardsToDeal,game.CardsToDeal)

	for _, player := range round.Players {
		hand := player.CardMap.ToSlice()
		assert.Equal(t,len(hand), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	assert.NotEqual(t, round.Players[0].CardMap, round.Players[1].CardMap)
	totalDealt := (game.CardsToDeal * len(players))
	expectedRemaining := expectedEuchreDeckSize - (cardsToDeal * len(round.Players))
	assert.Equal(t,expectedEuchreDeckSize,totalDealt + len(round.Deck.Cards))
	assert.Equal(t,expectedRemaining, len(round.Deck.Cards))
}

func TestNewRoundDealsCorrectly(t *testing.T) {
	game := CreateEuchreGame(CreatePlayers())
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	for _, player := range round.Players {
		hand := player.CardMap.ToSlice()
		for _, card := range hand {
			player.PlayCard(card)
		}
		hand = player.CardMap.ToSlice()
		assert.Equal(t, 0, len(hand), "expected empty hands")
	}
	game.NewRound()
	for _, player := range round.Players {
		hand := player.CardMap.ToSlice()
		assert.Equal(t,len(hand), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	assert.True(t,round.Deck.Cards[0].FaceUp)
}
func TestNewRoundScoreStartsAtZero(t *testing.T){
	game := CreateEuchreGame(CreatePlayers())
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
	game := CreateEuchreGame(CreatePlayers())
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	assert.Equal(t,cardsToDeal,game.CardsToDeal)
	for _, player := range round.Players {
		hand := player.CardMap.ToSlice()
		assert.Equal(t,len(hand), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	round.DetermineTrump()
	for _, player := range round.Players {
		hand := player.CardMap.ToSlice()
		assert.Equal(t,cardsToDeal,len(hand), "Expected players to still have 5 cards after trump declared")
	}
}
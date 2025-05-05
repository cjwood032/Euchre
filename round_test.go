package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoundDealsCorrectly(t *testing.T) {
	game := CreateEuchreGame(players)
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	for i := range round.Players {
		player := round.Players[i]
		for range player.Hand.Cards {
			player.Hand.Play(player.Hand.Cards[0])
		}
		assert.Equal(t, 0, len(player.Hand.Cards), "expected empty hands")
	}
	game.NewRound()
	fmt.Println(len(round.Deck.Cards))
	for i := range round.Players {
		fmt.Println(len(round.Players[i].Hand.Cards))
		assert.Equal(t,len(round.Players[i].Hand.Cards), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	assert.True(t,round.Deck.Cards[0].FaceUp)
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

func TestGameDealsCorrectly(t *testing.T) {
	game := CreateEuchreGame(players)
	game.NewRound()
	round := game.Rounds[len(game.Rounds)-1]
	assert.Equal(t,cardsToDeal,game.CardsToDeal)
	for i := range round.Players {
		assert.Equal(t,len(round.Players[i].Hand.Cards), cardsToDeal, "Expected players to be dealt 5 cards")
	}
	assert.NotEqual(t, round.Players[0].Hand, round.Players[1].Hand)
	totalDealt := (game.CardsToDeal * len(players))
	expectedRemaining := expectedEuchreDeckSize - (cardsToDeal * len(round.Players))
	assert.Equal(t,expectedEuchreDeckSize,totalDealt + len(round.Deck.Cards))
	assert.Equal(t,expectedRemaining, len(round.Deck.Cards))
}
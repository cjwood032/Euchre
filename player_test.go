package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func CreatePlayers() []*Player {
	var players = []*Player{
		{Name: "Chris"},
		{Name: "Don"},
		{Name: "MaryAnn"},
		{Name: "Andy"},
	}
	return players
}

func TestPlayCard(t *testing.T) {
	deck := NewStandardDeck()
	hand := deck.DealQuantity(5)
	p := Player{}
	p.CardMap.AddCardsToHand(hand)
	p.PlayCard(hand.Cards[0])
	assert.Equal(t, 4,len(p.CardMap.ToSlice()))
}

func TestPassBadHand(t *testing.T){
	card1 := NewCard(11,Diamonds)
	card2 := NewCard(1,Clubs)
	card3 := NewCard(9,Clubs)
	card4 := NewCard(13,Spades)
	card5 := NewCard(12,Hearts)

	newMap:= &CardMap{}
	newMap.AddToHand(card1)
	newMap.AddToHand(card2)
	newMap.AddToHand(card3)
	newMap.AddToHand(card4)
	newMap.AddToHand(card5)
	p := Player{CardMap: *newMap}
	actual := p.CallOrPass(Spades, true)
	assert.Equal(t, Pass, actual)
}
func TestOrderGoodHand(t *testing.T){
	card1 := NewCard(11,Clubs)
	card2 := NewCard(1,Clubs)
	card3 := NewCard(9,Clubs)
	card4 := NewCard(13,Spades)
	card5 := NewCard(1,Hearts)

	newMap:= &CardMap{}
	newMap.AddToHand(card1)
	newMap.AddToHand(card2)
	newMap.AddToHand(card3)
	newMap.AddToHand(card4)
	newMap.AddToHand(card5)
	p := Player{CardMap: *newMap}
	actual := p.CallOrPass(Clubs, true)
	assert.Equal(t, OrderUp, actual)
}

func TestOrderDependsOnPickup(t *testing.T){
	card1 := NewCard(11,Clubs)
	card2 := NewCard(1,Clubs)
	card3 := NewCard(9,Diamonds)
	card4 := NewCard(13,Spades)
	card5 := NewCard(1,Hearts)

	newMap:= &CardMap{}
	newMap.AddToHand(card1)
	newMap.AddToHand(card2)
	newMap.AddToHand(card3)
	newMap.AddToHand(card4)
	newMap.AddToHand(card5)
	p := Player{CardMap: *newMap}
	actual := p.CallOrPass(Clubs, true)
	assert.Equal(t, OrderUp, actual)
	actual = p.CallOrPass(Clubs, false)
	assert.Equal(t, Pass, actual)
}

func TestBestPlay_PlayerCanWinWithTrump(t *testing.T) {
	player := CreateTestPlayer("Active",
		&Deck{Cards: []*Card{
			NewCard(11, Spades),
			NewCard(9, Hearts),
		}})
	opponent1 := CreateTestPlayer("Opponent1",
		&Deck{Cards: []*Card{
			NewCard(10, Spades),
			NewCard(9, Diamonds),
		}})
	partner := CreateTestPlayer( "Partner",
		&Deck{Cards: []*Card{
			NewCard(12, Spades),
			NewCard(11, Hearts),
		}})
	opponent2 := CreateTestPlayer( "Opponent2",
		&Deck{Cards: []*Card{
			NewCard(13, Clubs),
			NewCard(9, Clubs),
		}})

	currentTrick := []*Card{
		NewCard(13, Clubs),
		NewCard(12, Clubs), 
		NewCard(1, Clubs),
	}

	round := Round{
		Trump:  Spades,
		Caller: player,
		Players: []*Player{player, opponent1, partner, opponent2},
	}

	best := player.BestPlay(currentTrick, round)
	expected:= *NewCard(11, Spades) 
	assert.Equal(t,expected,best)
}
func CreateTestPlayer(name string, cards *Deck) *Player{
	player := &Player{
		Name: name,
	}
	player.CardMap.AddCardsToHand(cards)
	return player
}
func TestBestPlay_PlayerHasLeadSuit(t *testing.T) {
	player := CreateTestPlayer("Tester",
	&Deck{Cards: []*Card{
		NewCard(10, Clubs), 
		NewCard(9, Hearts),
	}})
	opponent1 := CreateTestPlayer("Opponent1",
	&Deck{Cards: []*Card{
			NewCard(10, Spades),
			NewCard(9, Diamonds),
		}})
	
	partner := CreateTestPlayer("Partner",
		&Deck{Cards: []*Card{
			NewCard(12, Spades),
			NewCard(11, Hearts),
		}})
	opponent2 := CreateTestPlayer("Opponent2",
		&Deck{Cards: []*Card{
			NewCard(13, Clubs),
			NewCard(9, Clubs),
		}})

	currentTrick := []*Card{
		NewCard(13, Clubs),
		NewCard(12, Clubs),
	}


	round := Round{
		Trump:  Spades,
		Caller: player,
		Players: []*Player{player, opponent1, partner, opponent2},
	}

	best := player.BestPlay(currentTrick, round)
	if best != *NewCard(10, Clubs) {
		t.Errorf("Expected to follow suit with 10 of Clubs, got %+v", best)
	}
}

func TestBestPlay_PlayerIsShortSuited(t *testing.T) {
	player := CreateTestPlayer("Tester",
		&Deck{Cards: []*Card{
			NewCard(9, Spades),  
			NewCard(1, Hearts),  
			NewCard(10, Diamonds),
			NewCard(8, Diamonds),
		}})
	
	opponent1 := CreateTestPlayer("Opponent1",
		&Deck{Cards: []*Card{
			NewCard(10, Spades),
			NewCard(9, Diamonds),
		}})
	partner := CreateTestPlayer("Partner",
		&Deck{Cards: []*Card{
			NewCard(12, Spades),
			NewCard(11, Hearts),
		}})
	opponent2 := CreateTestPlayer("Opponent2",
		&Deck{Cards: []*Card{
			NewCard(13, Clubs),
			NewCard(9, Clubs),
		}})

	currentTrick := []*Card{
		NewCard(13, Clubs),
		NewCard(12, Clubs),
		NewCard(1, Clubs),
	}

	round := Round{
		Trump:  Spades,
		Caller: partner,
		Players: []*Player{player, opponent1, partner, opponent2},
	}

	best := player.BestPlay(currentTrick, round)

	if best.Suit != Spades {
		t.Errorf("Expected to play trump to try to win, got %+v", best)
	}
}

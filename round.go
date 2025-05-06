package main

type Round struct {
    Players []*Player
    Dealer int
    Caller *Player
    TricksWon int
    Deck *Deck
    Trump Suit
    Turn int
    Lead int
    Alone bool
}

func (round *Round) Begin() {
    // Ensure we have a valid deck with cards
    if round.Deck == nil || len(round.Deck.Cards) == 0 {
        round.Deck = NewSpecificDeck([]int{9, 10, 11, 12, 13, 1}, []Suit{Spades, Diamonds, Clubs, Hearts})
    }
    round.Deck.Shuffle()
    round.Deal()
}

func (round *Round) Deal() {
    // Ensure we have enough cards to deal (5 cards to each of 4 players = 20 cards)
    if len(round.Deck.Cards) < 20 {
        round.Deck = NewSpecificDeck([]int{9, 10, 11, 12, 13, 1}, []Suit{Spades, Diamonds, Clubs, Hearts})
        round.Deck.Shuffle()
    }

    // Deal 5 cards to each player
    for _, player := range round.Players {
        cards := round.Deck.DealQuantity(5)
        if len(cards.Cards) < 5 {
            panic("Not enough cards in deck to deal")
        }
        player.CardMap.AddCardsToHand(cards)
    }
    
    // Set the top card face up
    if len(round.Deck.Cards) > 0 {
        round.Deck.Cards[0].TurnFaceUp()
    }
}

func (round *Round)DetermineTrump() {
	for i := range round.Players	{
		// 0, 1, 2, 3
		// dealer is 2 order goes 3,0,1,2
		playerPosition := round.Dealer + i;
		if playerPosition >= len(round.Players) {
			playerPosition -= 4
		}
		player := round.Players[playerPosition]
		suit := round.Deck.Cards[0].Suit
		call := player.CallOrPass(suit, round.Dealer % 2 == playerPosition % 2) //todo use the player round
		if ( call != Pass) {
			round.BeginPlay(call, suit)
			return
		}
	}
	
	passedSuit := round.Deck.Cards[0].Suit
	for i := range round.Players{
		playerPosition := round.Dealer + i;
		if playerPosition >= len(round.Players) {
			playerPosition -= 4
		}
		player := round.Players[playerPosition]
		call, trump := player.DeclareTrump(passedSuit)
		if ( call != Pass) {
			round.BeginPlay(call, trump)
			return
		}
	}
	//dealer turned it up

	//each player can call or pass

	//if called dealer picks up
	//if all pass, dealer turns down
	//call trump from available suits
}

func (round *Round) BeginPlay(call Call, trump Suit) {
	
	
	if(round.Deck.Cards[0].FaceUp){
		round.Players[round.Dealer].PickUp(round.Deck.Cards[0])
	}
	//round.Players[round.LeftOfDealer()].PlayCard()
}
func (round *Round) LeftOfDealer() int{
	if round.Dealer == len(round.Players) -1 {
		return 0
	}
	return round.Dealer +1
}
func (r *Round) DetermineTrickWinner(trick []*Card, lead int) int {
	winningIndex := lead
	winningCard := trick[lead]

	for i := 1; i < 4; i++ {
		pos := (lead + i) % 4
		card := trick[pos]
		if card.Beats(winningCard, r.Trump, trick[lead].Suit) {
			winningCard = card
			winningIndex = pos
		}
	}
	return winningIndex
}

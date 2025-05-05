package main

type Round struct{
	Players []*Player
	Dealer int
	Caller *Player
	Tricks *Deck
	Deck *Deck
	Trump Suit
}

type PlayerRound struct {
	Player *Player
	TricksWon int
}

func (round *Round) Begin() {
	round.Deal()
	//round.DetermineTrump()
}

func (round *Round) Deal() {

	
	round.Deck.Shuffle()
	// starting with the player to the left of the dealer 
	for seat := 0; seat < len(round.Players); seat++ {
		round.Players[seat].Hand = round.Deck.DealQuantity(5) // deal the appropriate amount
	}
	round.Deck.Cards[0].TurnFaceUp()
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
		call := player.CallOrPass(suit, round.Dealer % 2 == playerPosition % 2)
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
	//todo: loner logic
	
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
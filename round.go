package main

import "fmt"

type Round struct {
	Players        []*Player
	Dealer         int
	Caller         *Player
	TricksWon      int
	Deck           *Deck
	Trump          Suit
	Turn           int
	Lead           int
	Alone          bool
	SelectingTrump bool
	ActivePlayer   int
}

func (round *Round) Begin() {
	// Ensure we have a valid deck with cards
	if round.Deck == nil || len(round.Deck.Cards) == 0 {
		round.Deck = NewSpecificDeck([]int{9, 10, 11, 12, 13, 1}, []Suit{Spades, Diamonds, Clubs, Hearts})
	}
	round.SelectingTrump = true
	round.Deck.Shuffle()
	round.ActivePlayer = (round.Dealer + 1) % 4
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

func (round *Round) DetermineTrump() {
	round.SelectingTrump = true

	// First round - ordering up the top card
	for i := range round.Players {
		playerPosition := (round.Dealer + i + 1) % len(round.Players) // Start with player left of dealer
		player := round.Players[playerPosition]
		round.ActivePlayer = playerPosition

		if player.ComputerPlayer {
			// Computer player makes automatic decision
			suit := round.Deck.Cards[0].Suit
			call := player.CallOrPass(suit, round.Dealer%2 == playerPosition%2)
			if call != Pass {
				round.BeginPlay(call, suit)
				round.SelectingTrump = false
				return
			}
		} else {
			// Human player - we'll handle this in the UI
			return // Wait for UI to make selection
		}
	}

	// If we get here, all passed first round - turn down the top card
	if len(round.Deck.Cards) > 0 {
		round.Deck.Cards[0].TurnFaceDown()
	}

	// Second round - calling any suit except the turned-down one
	passedSuit := Suit(-1) // Initialize with invalid suit
	if len(round.Deck.Cards) > 0 {
		passedSuit = round.Deck.Cards[0].Suit
	}

	for i := range round.Players {
		playerPosition := (round.Dealer + i + 1) % len(round.Players)
		player := round.Players[playerPosition]
		round.ActivePlayer = playerPosition

		if player.ComputerPlayer {
			// Computer player makes automatic decision
			call, trump := player.DeclareTrump(passedSuit)
			if call != Pass {
				round.BeginPlay(call, trump)
				round.SelectingTrump = false
				return
			}
		} else {
			// Human player - we'll handle this in the UI
			return // Wait for UI to make selection
		}
	}

	// If all passed both rounds, dealer must pick
	round.ActivePlayer = round.Dealer
	player := round.Players[round.Dealer]
	call, trump := player.DeclareTrump(passedSuit)
	if call != Pass {
		round.BeginPlay(call, trump)
	} else {
		// Shouldn't happen - dealer must pick something
		// Default to first available suit
		for _, s := range []Suit{Spades, Diamonds, Clubs, Hearts} {
			if s != passedSuit {
				round.BeginPlay(OrderUp, s)
				break
			}
		}
	}
	round.SelectingTrump = false
}

func (round *Round) HumanTrumpSelection(call Call, trump Suit) {
	if !round.SelectingTrump || round.ActivePlayer < 0 || round.ActivePlayer >= len(round.Players) {
		fmt.Println("Unexpected Trump selection, exiting")
		return
	}

	player := round.Players[round.ActivePlayer]
	if player.ComputerPlayer {
		fmt.Println("Unexpected Trump selection, computer player got in hereexiting")
		return // Not a human player
	}

	if call == Pass {
		fmt.Println("Human passes")
		if round.ActivePlayer == round.Dealer { // All players have passed
			if len(round.Deck.Cards) > 0 && round.Deck.Cards[0].FaceUp {
				// First round passed - flip card and go to second round
				round.Deck.Cards[0].FaceUp = false
			} else {
				// Second round passed - redeal
				round.SelectingTrump = false // ... existing pass logic ...
			}
		}
		round.ActivePlayer = (round.ActivePlayer + 1) % 4
	} else {
		fmt.Printf("Player calls %s as trump\n", trump.FriendlySuit())
		round.Trump = trump
		round.Caller = player
		if call == Alone {
			round.Alone = true
		}

		// If this is first round and card is face up
		if len(round.Deck.Cards) > 0 && round.Deck.Cards[0].FaceUp {
			if round.ActivePlayer == round.Dealer {
				// Dealer is ordering - pick up card and discard one
				pickedCard := round.Deck.Cards[0]
				player.PickUp(pickedCard)
			} else {
				// Non-dealer ordering - dealer picks up card
				round.Players[round.Dealer].PickUp(round.Deck.Cards[0])
			}
			round.Deck.Cards = round.Deck.Cards[1:] // Remove from kitty
		}

		round.SelectingTrump = false
		round.BeginPlay(call, trump) // This will hide the kitty
	}
	
}

func (r *Round) ComputerTrumpSelection(decision Call, suit Suit) {
	switch decision {
	case OrderUp, Alone:
		r.Trump = suit
		r.Caller = r.Players[r.ActivePlayer]
		r.Alone = (decision == Alone)
		r.SelectingTrump = false

		if len(r.Deck.Cards) > 0 && r.Deck.Cards[0].FaceUp {
			// Only dealer picks up the card
			if r.ActivePlayer == r.Dealer {
				pickedCard := r.Deck.Cards[0]
				r.Players[r.Dealer].PickUp(pickedCard)
				r.ComputerDealerDiscard()
			}
			r.Deck.Cards = r.Deck.Cards[1:]
		}

		r.BeginPlay(decision, suit)

	case Pass:
		r.ActivePlayer = (r.ActivePlayer + 1) % 4
		if r.ActivePlayer == r.Dealer {
			if len(r.Deck.Cards) > 0 && r.Deck.Cards[0].FaceUp {
				r.Deck.Cards[0].FaceUp = false
			} else {
				r.SelectingTrump = false
			}
		}
	}
}

func (round *Round) ComputerDealerDiscard() {
	dealer := round.Players[round.Dealer]
	if len(dealer.CardMap.ToSlice()) <= 5 {
		return // No need to discard
	}

	// Simple AI - discard weakest non-trump card
	var discard *Card
	for _, card := range dealer.CardMap.ToSlice() {
		if card.Suit != round.Trump {
			if discard == nil || card.Rank < discard.Rank {
				discard = card
			}
		}
	}

	// If all cards are trump, discard lowest trump
	if discard == nil {
		for _, card := range dealer.CardMap.ToSlice() {
			if discard == nil || card.Rank < discard.Rank {
				discard = card
			}
		}
	}

	if discard != nil {
		dealer.CardMap.RemoveFromHand(*discard)
	}
}

func (round *Round) BeginPlay(call Call, trump Suit) {
	
	round.SelectingTrump = false
	round.Lead = (round.Dealer + 1) % len(round.Players) // Left of dealer leads first trick
	round.ActivePlayer = round.Lead
	round.Trump = trump
	fmt.Printf("Beginning play, trump is %s, first lead is %v\n", trump.FriendlySuit(), round.Lead)
	for _, p := range round.Players {
		p.IsPlaying = true
	}

	// Handle "going alone"
	if call == Alone {
		partnerPos := (round.Caller.Position + 2) % 4
		round.Players[partnerPos].IsPlaying = false
	}

	if len(round.Deck.Cards) > 0 {
		round.Deck.Cards[0].TurnFaceDown()
	}
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

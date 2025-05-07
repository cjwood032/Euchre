package main

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
		return
	}

	player := round.Players[round.ActivePlayer]
	if player.ComputerPlayer {
		return // Not a human player
	}

	if call == Pass {
		// Continue to next player
		round.DetermineTrump()
	} else {
		round.Trump = trump
		round.Caller = player
		if call == Alone {
			round.Alone = true
		}

		// If this is first round and card is face up, dealer picks it up
		if len(round.Deck.Cards) > 0 && round.Deck.Cards[0].FaceUp {
			if round.ActivePlayer == round.Dealer {
				// Dealer is ordering - pick up card and discard one
				pickedCard := round.Deck.Cards[0]
				player.PickUp(pickedCard)

				// Show discard UI (handled in gameUI)
				round.SelectingTrump = false
				return
			} else {
				// Non-dealer ordering - dealer picks up card
				round.Players[round.Dealer].PickUp(round.Deck.Cards[0])
			}
			round.Deck.Cards = round.Deck.Cards[1:] // Remove from kitty
		}

		round.SelectingTrump = false
	}
}

func (r *Round) ComputerTrumpSelection(decision Call, suit Suit) {
	switch decision {
	case OrderUp, Alone:
		r.Trump = suit
		r.Caller = r.Players[r.ActivePlayer]
		if decision == Alone {
			r.Alone = true
		}

		// If this is first round and card is face up
		if len(r.Deck.Cards) > 0 && r.Deck.Cards[0].FaceUp {
			if r.ActivePlayer == r.Dealer {
				// Dealer is ordering - pick up card and discard one
				pickedCard := r.Deck.Cards[0]
				r.Players[r.Dealer].PickUp(pickedCard)
				r.ComputerDealerDiscard()
			} else {
				// Non-dealer ordering - dealer picks up card
				r.Players[r.Dealer].PickUp(r.Deck.Cards[0])
			}
			r.Deck.Cards = r.Deck.Cards[1:] // Remove from kitty
		}

		r.SelectingTrump = false
	case Pass:
		r.ActivePlayer = (r.ActivePlayer + 1) % 4
		if r.ActivePlayer == r.Dealer { // All players have passed
			if len(r.Deck.Cards) > 0 && r.Deck.Cards[0].FaceUp {
				// First round passed - flip card and go to second round
				r.Deck.Cards[0].FaceUp = false
				r.ActivePlayer = (r.Dealer + 1) % 4
			} else {
				// Second round passed - redeal
				r.SelectingTrump = false // ... existing pass logic ...
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
	round.Trump = trump
	round.Caller = round.Players[round.ActivePlayer]
	round.Alone = (call == Alone)

	// Hide the kitty by removing the face-up card if it exists
	if len(round.Deck.Cards) > 0 && round.Deck.Cards[0].FaceUp {
		if round.ActivePlayer == round.Dealer {
			// Dealer picks up the card
			round.Players[round.Dealer].PickUp(round.Deck.Cards[0])
		}
		round.Deck.Cards = round.Deck.Cards[1:] // Remove from kitty
	}

	// Set the first player to left of dealer
	round.Lead = (round.Dealer + 1) % len(round.Players)
	round.ActivePlayer = round.Lead
	round.SelectingTrump = false
}

func (round *Round) LeftOfDealer() int {
	if round.Dealer == len(round.Players)-1 {
		return 0
	}
	return round.Dealer + 1
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

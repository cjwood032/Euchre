package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type GameUI struct {
	Window      fyne.Window
	MainContent fyne.CanvasObject
	Players     []*Player
	Round       *Round
	Game        *Game

	// UI components
	HandBox        *fyne.Container
	KittyContainer *fyne.Container
	CenterNorth    *fyne.Container
	CenterEast     *fyne.Container
	CenterSouth    *fyne.Container
	CenterWest     *fyne.Container
	NorthScore     *widget.Label
	EastScore      *widget.Label
	SouthScore     *widget.Label
	WestScore      *widget.Label

	Trick                []*Card
	NewGameBtn           *widget.Button
	SouthHandBox         *fyne.Container
	BottomArea           *fyne.Container
	NorthDealerIndicator *widget.Label
	EastDealerIndicator  *widget.Label
	SouthDealerIndicator *widget.Label
	WestDealerIndicator  *widget.Label
	discardDialog        *widget.PopUp
}

func (ui *GameUI) RefreshUI() {
	// Update all static elements
	if ui.discardDialog != nil {
		ui.discardDialog.Hide()
	}
	ui.updateDealerIndicators()
	ui.updateTrickDisplay(make([]*Card, 4))

	// Refresh kitty
	kitty := createStackedKitty(ui.Round, fyne.NewSize(70, 110))
	ui.KittyContainer.Objects = []fyne.CanvasObject{kitty}

	// Update scores
	ui.NorthScore.SetText(fmt.Sprintf("Score: %d", ui.Players[0].Score))
	ui.EastScore.SetText(fmt.Sprintf("Score: %d", ui.Players[1].Score))
	ui.SouthScore.SetText(fmt.Sprintf("Score: %d", ui.Players[2].Score))
	ui.WestScore.SetText(fmt.Sprintf("Score: %d", ui.Players[3].Score))

	// Always update hand (maintains single instance)
	ui.updateHumanHand()

	ui.updateDealerIndicators()
	ui.updateTrickDisplay(make([]*Card, 4))
	ui.updateHumanHand()

	if ui.Round.SelectingTrump {
		if ui.Round.ActivePlayer == 2 { // Human's turn to order
			ui.showTrumpSelection()
		} else {
			// Computer decides automatically
			ui.processComputerTrumpSelection()
		}
	} else {
		// Card play phase
		ui.Window.SetContent(ui.MainContent)
		if ui.Round.ActivePlayer == 2 {
			// Human's turn to play card
		} else {
			// Computer plays automatically
			ui.playComputerTurn()
		}
	}
}

func (ui *GameUI) updateHumanHand() {
	// Get or create the hand container
	var handContainer *fyne.Container
	if len(ui.HandBox.Objects) > 0 {
		handContainer = ui.HandBox.Objects[0].(*fyne.Container)
		handContainer.Objects = nil // Clear existing cards
	} else {
		handContainer = container.NewHBox()
		ui.HandBox.Objects = []fyne.CanvasObject{handContainer}
	}

	player := ui.Players[2] // South is human player
	cardSize := fyne.NewSize(80, 120)

	// Only show play buttons if we're in the playing phase (not trump selection)
	showPlayButtons := !ui.Round.SelectingTrump && ui.Round.Trump != Suit(-1)

	for _, card := range player.CardMap.ToSlice() {
		currentCard := card
		cardUI := container.NewVBox(
			renderCardImage(currentCard, cardSize),
		)

		if showPlayButtons {
			playBtn := widget.NewButton("Play", func() {
				// ... existing play logic ...
			})
			cardUI.Add(playBtn)
		}
		handContainer.Add(cardUI)
	}

	ui.HandBox.Refresh()
}

func (ui *GameUI) showTrumpSelection() {
	if !ui.Round.SelectingTrump || ui.Round.ActivePlayer != 2 {
		return
	}

	// First update the hand (without play buttons)
	ui.updateHumanHand()

	firstRound := len(ui.Round.Deck.Cards) > 0 && ui.Round.Deck.Cards[0].FaceUp

	// Create a container for the trump selection UI
	trumpSelectionContainer := container.NewHBox()

	if firstRound {
		topCard := ui.Round.Deck.Cards[0]
		trumpSelectionContainer.Add(widget.NewLabel(fmt.Sprintf("Top card is %s of %s", topCard.FriendlyRank(), topCard.Suit.FriendlySuit())))
		trumpSelectionContainer.Add(widget.NewLabel("Do you want to:"))

		orderUpBtn := widget.NewButton("Order Up", func() {
			ui.Round.HumanTrumpSelection(OrderUp, topCard.Suit)
			if ui.Round.Dealer == 2 { // Human is dealer
				ui.showDiscardSelection()
			}
			ui.RefreshUI()
		})
		orderUpBtn.Importance = widget.HighImportance

		goAloneBtn := widget.NewButton("Go Alone", func() {
			ui.Round.HumanTrumpSelection(Alone, topCard.Suit)
			if ui.Round.Dealer == 2 { // Human is dealer
				ui.showDiscardSelection()
			}
			ui.RefreshUI()
		})

		passBtn := widget.NewButton("Pass", func() {
			ui.Round.HumanTrumpSelection(Pass, topCard.Suit)
			ui.processComputerTrumpSelection() // Continue with next players
		})

		trumpSelectionContainer.Add(orderUpBtn)
		trumpSelectionContainer.Add(goAloneBtn)
		trumpSelectionContainer.Add(passBtn)
	} else {
		passedSuit := Suit(-1)
		if len(ui.Round.Deck.Cards) > 0 {
			passedSuit = ui.Round.Deck.Cards[0].Suit
		}

		suitButtons := container.NewHBox()
		for _, suit := range []Suit{Spades, Diamonds, Clubs, Hearts} {
			if suit != passedSuit {
				currentSuit := suit
				btn := widget.NewButton(suit.FriendlySuit(), func() {
					ui.Round.HumanTrumpSelection(OrderUp, currentSuit)
					ui.RefreshUI()
				})
				btn.Importance = widget.MediumImportance
				suitButtons.Add(btn)
			}
		}

		trumpSelectionContainer.Add(widget.NewLabel("Choose a trump suit:"))
		trumpSelectionContainer.Add(suitButtons)

		passBtn := widget.NewButton("Pass", func() {
			ui.Round.HumanTrumpSelection(Pass, Suit(-1))
			ui.processComputerTrumpSelection() // Continue with next players
		})
		trumpSelectionContainer.Add(passBtn)
	}

	// Create the complete content
	content := container.NewBorder(
		ui.MainContent.(*fyne.Container).Objects[0], // Top controls
		container.NewHBox( // Bottom section
			container.NewCenter(trumpSelectionContainer),
		),
		ui.MainContent.(*fyne.Container).Objects[2], // West
		ui.MainContent.(*fyne.Container).Objects[3], // East
		ui.MainContent.(*fyne.Container).Objects[4], // Center
	)

	ui.Window.SetContent(content)
}

func (ui *GameUI) updateTrickDisplay(trick []*Card) {
	// Clear only if we're passing an empty trick
	if len(trick) == 0 {
		ui.CenterNorth.Objects = nil
		ui.CenterEast.Objects = nil
		ui.CenterSouth.Objects = nil
		ui.CenterWest.Objects = nil
		return
	}

	// Otherwise, update each position with current card
	cardSize := fyne.NewSize(80, 120)
	for i, card := range trick {
		if card == nil {
			continue
		}

		switch i {
		case 0: // North
			ui.CenterNorth.Objects = []fyne.CanvasObject{renderCardImage(card, cardSize)}
		case 1: // East
			ui.CenterEast.Objects = []fyne.CanvasObject{renderCardImage(card, cardSize)}
		case 2: // South
			ui.CenterSouth.Objects = []fyne.CanvasObject{renderCardImage(card, cardSize)}
		case 3: // West
			ui.CenterWest.Objects = []fyne.CanvasObject{renderCardImage(card, cardSize)}
		}
	}

	// Refresh all containers
	ui.CenterNorth.Refresh()
	ui.CenterEast.Refresh()
	ui.CenterSouth.Refresh()
	ui.CenterWest.Refresh()
}

func (ui *GameUI) showComputerDecision(player *Player, decision string, suit Suit) {
	var position *fyne.Container
	var label *widget.Label

	// Determine which position to show the decision
	switch player.Name {
	case "NORTH":
		position = ui.CenterNorth
		label = ui.NorthScore
	case "EAST":
		position = ui.CenterEast
		label = ui.EastScore
	case "WEST":
		position = ui.CenterWest
		label = ui.WestScore
	default:
		return
	}

	// Clear any previous decision
	position.Objects = nil

	// Create decision text
	decisionText := decision
	if suit != Suit(-1) {
		decisionText += " " + suit.FriendlySuit()
	}
	decisionLabel := widget.NewLabel(decisionText)
	decisionLabel.Alignment = fyne.TextAlignCenter
	decisionLabel.TextStyle.Bold = true
	position.Add(decisionLabel)

	// Temporarily update the score label
	originalText := label.Text
	label.SetText("Thinking...")
	ui.Window.Content().Refresh()

	// Pause for visibility
	time.Sleep(1 * time.Second)

	// Restore original label and clear decision
	label.SetText(originalText)
	position.Objects = nil
	ui.Window.Content().Refresh()
}

func (ui *GameUI) processComputerTrumpSelection() {
	if !ui.Round.SelectingTrump || ui.Round.ActivePlayer == 2 { // Human player's position
		return
	}

	currentPlayer := ui.Round.Players[ui.Round.ActivePlayer]
	if !currentPlayer.ComputerPlayer {
		return
	}

	var decision Call
	var suit Suit

	if len(ui.Round.Deck.Cards) > 0 && ui.Round.Deck.Cards[0].FaceUp {
		// First round decision
		suit = ui.Round.Deck.Cards[0].Suit
		decision = currentPlayer.CallOrPass(suit, ui.Round.Dealer%2 == ui.Round.ActivePlayer%2)
	} else {
		// Second round decision
		passedSuit := Suit(-1)
		if len(ui.Round.Deck.Cards) > 0 {
			passedSuit = ui.Round.Deck.Cards[0].Suit
		}
		decision, suit = currentPlayer.DeclareTrump(passedSuit)
	}

	// Show the computer's decision
	ui.showComputerDecision(currentPlayer, decision.FriendlyCall(), suit)

	// Process the decision
	if decision == Pass {
		ui.Round.ComputerTrumpSelection(decision, suit)
		if ui.Round.SelectingTrump && ui.Round.ActivePlayer == 2 {
			ui.showTrumpSelection()
		}
	} else {
		ui.Round.ComputerTrumpSelection(decision, suit)
		// Force immediate UI refresh after computer selects trump
		ui.RefreshUI()

		// If human is dealer and computer ordered up, show discard UI
		if ui.Round.Dealer == 2 && len(ui.Round.Deck.Cards) > 0 &&
			ui.Round.Deck.Cards[0].FaceUp && decision == OrderUp {
			ui.showDiscardSelection()
		}
	}
}

func (ui *GameUI) showCardPickup() {
	if len(ui.Round.Deck.Cards) == 0 {
		return
	}

	// Get dealer position
	var dealerPos *fyne.Container
	switch ui.Round.Dealer {
	case 0: // North
		dealerPos = ui.CenterNorth
	case 1: // East
		dealerPos = ui.CenterEast
	case 2: // South
		dealerPos = ui.CenterSouth
	case 3: // West
		dealerPos = ui.CenterWest
	}

	// Create animation of card moving to dealer's position
	card := ui.Round.Deck.Cards[0]
	cardImg := renderCardImage(card, fyne.NewSize(80, 120))

	// Start position (kitty)
	startPos := ui.KittyContainer.Position()
	cardImg.Move(fyne.NewPos(startPos.X, startPos.Y))

	// Add to overlay
	overlay := container.NewWithoutLayout(cardImg)
	ui.Window.Canvas().SetContent(container.NewStack(ui.Window.Content(), overlay))

	// Animate movement
	endPos := dealerPos.Position()
	anim := fyne.NewAnimation(time.Second, func(f float32) {
		x := startPos.X + f*(endPos.X-startPos.X)
		y := startPos.Y + f*(endPos.Y-startPos.Y)
		cardImg.Move(fyne.NewPos(x, y))
		overlay.Refresh()
	})

	anim.Start()
	time.Sleep(time.Second) // Let animation finish

	// Remove animation and refresh
	ui.Window.SetContent(ui.Window.Content())

	// Hide the kitty by refreshing UI
	ui.RefreshUI()
}

func (ui *GameUI) showDealerDiscard() {
	if ui.Round.Dealer != 2 { // Human is not dealer
		return
	}

	dealer := ui.Round.Players[ui.Round.Dealer]
	if len(dealer.CardMap.ToSlice()) <= 5 {
		return // No need to discard
	}

	// Create discard selection UI
	discardUI := container.NewVBox(
		widget.NewLabel("Choose a card to discard:"),
	)

	cardSize := fyne.NewSize(80, 120)
	for _, card := range dealer.CardMap.ToSlice() {
		currentCard := card
		cardUI := container.NewVBox(
			renderCardImage(currentCard, cardSize),
			widget.NewButton("Discard", func() {
				dealer.CardMap.RemoveFromHand(*currentCard)
				ui.RefreshUI() // Refresh to show updated hand
			}),
		)
		discardUI.Add(cardUI)
	}

	// Show modal dialog
	dialog := widget.NewModalPopUp(
		discardUI,
		ui.Window.Canvas(),
	)
	dialog.Show()
}

func (ui *GameUI) playComputerTurn() {
	if ui.Round.ActivePlayer == 2 || ui.Round.SelectingTrump {
		return // Not computer's turn
	}

	// Get current trick state
	var partialTrick []*Card
	for _, card := range ui.Trick {
		if card != nil {
			partialTrick = append(partialTrick, card)
		}
	}

	// Computer makes play
	computer := ui.Round.Players[ui.Round.ActivePlayer]
	play := computer.BestPlay(partialTrick, *ui.Round)
	playedCard := computer.PlayCard(&play)
	ui.Trick[ui.Round.ActivePlayer] = playedCard // Store the played card

	// Update the display with all cards in the trick
	ui.updateTrickDisplay(ui.Trick)

	// Move to next player
	ui.Round.ActivePlayer = (ui.Round.ActivePlayer + 1) % 4

	// If trick is complete, determine winner
	if len(partialTrick) == 3 { // All 4 players have played
		winner := resolveTrick(ui.Trick, ui.Round)
		ui.Round.Lead = winner
		ui.Round.ActivePlayer = winner

		// Keep cards visible for a moment before clearing
		time.Sleep(1 * time.Second)
		ui.Trick = make([]*Card, 4)     // Clear trick for next round
		ui.updateTrickDisplay(ui.Trick) // Update display to show cleared trick
	}

	ui.RefreshUI()
}

func (ui *GameUI) createDealerIndicator(position int) *widget.Label {
	indicator := widget.NewLabel("")
	indicator.Hide() // Start hidden

	// Store reference to update later
	switch position {
	case 0:
		ui.NorthDealerIndicator = indicator
	case 1:
		ui.EastDealerIndicator = indicator
	case 2:
		ui.SouthDealerIndicator = indicator
	case 3:
		ui.WestDealerIndicator = indicator
	}

	return indicator
}

func (ui *GameUI) updateDealerIndicators() {
	// Hide all indicators first
	if ui.NorthDealerIndicator != nil {
		ui.NorthDealerIndicator.Hide()
	}
	if ui.EastDealerIndicator != nil {
		ui.EastDealerIndicator.Hide()
	}
	if ui.SouthDealerIndicator != nil {
		ui.SouthDealerIndicator.Hide()
	}
	if ui.WestDealerIndicator != nil {
		ui.WestDealerIndicator.Hide()
	}

	// Show only for current dealer
	if ui.Round == nil {
		return
	}

	switch ui.Round.Dealer {
	case 0:
		if ui.NorthDealerIndicator != nil {
			ui.NorthDealerIndicator.SetText("(Dealer)")
			ui.NorthDealerIndicator.Show()
		}
	case 1:
		if ui.EastDealerIndicator != nil {
			ui.EastDealerIndicator.SetText("(Dealer)")
			ui.EastDealerIndicator.Show()
		}
	case 2:
		if ui.SouthDealerIndicator != nil {
			ui.SouthDealerIndicator.SetText("(Dealer)")
			ui.SouthDealerIndicator.Show()
		}
	case 3:
		if ui.WestDealerIndicator != nil {
			ui.WestDealerIndicator.SetText("(Dealer)")
			ui.WestDealerIndicator.Show()
		}
	}
}

func (ui *GameUI) clearTrumpSelection() {
	// Reset to main content immediately
	ui.Window.SetContent(ui.MainContent)
	ui.RefreshUI()
}

func (ui *GameUI) showDiscardSelection() {
	if ui.Round.Dealer != 2 || len(ui.Players[2].CardMap.ToSlice()) <= 5 {
		return // Not human dealer or no card to discard
	}

	player := ui.Players[2]
	cardSize := fyne.NewSize(80, 120)

	// Create the discard selection content
	content := container.NewVBox(
		widget.NewLabelWithStyle("Select a card to discard:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
	)

	handContainer := container.NewHBox()
	content.Add(handContainer)

	for _, card := range player.CardMap.ToSlice() {
		currentCard := card
		cardUI := container.NewVBox(
			renderCardImage(currentCard, cardSize),
			widget.NewButton("Discard", func() {
				player.CardMap.RemoveFromHand(*currentCard)
				ui.discardDialog.Hide()
				ui.RefreshUI() // Return to normal play
			}),
		)
		handContainer.Add(cardUI)
	}

	// Add cancel button
	content.Add(widget.NewButton("Cancel", func() {
		ui.discardDialog.Hide()
		ui.RefreshUI()
	}))

	// Create and show the modal dialog
	ui.discardDialog = widget.NewModalPopUp(
		container.NewBorder(
			nil, nil, nil, nil,
			content,
		),
		ui.Window.Canvas(),
	)
	ui.discardDialog.Show()
}

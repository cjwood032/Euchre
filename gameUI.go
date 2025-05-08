package main

import (
	"fmt"
	"runtime"
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

	Trick                [4]*Card
	NewGameBtn           *widget.Button
	SouthHandBox         *fyne.Container
	BottomArea           *fyne.Container
	NorthDealerIndicator *widget.Label
	EastDealerIndicator  *widget.Label
	SouthDealerIndicator *widget.Label
	WestDealerIndicator  *widget.Label
	discardDialog        *widget.PopUp
	CallerIndicator      *widget.Label
}

func (ui *GameUI) RefreshUI() {
	// Update all static elements
	if ui.discardDialog != nil {
		ui.discardDialog.Hide()
	}
	ui.updateDealerIndicators()
	ui.updateCallerIndicator()
	ui.updateTrickDisplay([4]*Card(make([]*Card, 4)))

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
	ui.updateTrickDisplay([4]*Card(make([]*Card, 4)))
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
				// Play card
				playedCard := ui.Players[2].PlayCard(currentCard)
				ui.Trick[2] = playedCard
				ui.updateTrickDisplay(ui.Trick)

				// If trick is complete
				if ui.isTrickComplete() {
					winner := resolveTrick(ui.Trick[:], ui.Round)
					ui.Round.Lead = winner
					ui.Round.ActivePlayer = winner

					time.Sleep(1 * time.Second)
					ui.clearTrickDisplay()
					ui.Trick = [4]*Card{}
				}
				// Move to next player
				ui.Round.ActivePlayer = (ui.Round.ActivePlayer + 1) % 4
				ui.updateHumanHand()

				// Process computer turns if needed
				if ui.Round.ActivePlayer != 2 {
					ui.playComputerTurn()
				}
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
	ui.updateCallerIndicator()
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

func (ui *GameUI) updateTrickDisplay(trick [4]*Card) {
	// Clear previous cards first
	ui.clearTrickDisplay()

	cardSize := fyne.NewSize(80, 120)

	// Update each position with current card
	for i, card := range trick {
		if card == nil {
			continue
		}

		img := renderCardImage(card, cardSize)
		switch i {
		case 0:
			ui.CenterNorth.Add(img)
		case 1:
			ui.CenterEast.Add(img)
		case 2:
			ui.CenterSouth.Add(img)
		case 3:
			ui.CenterWest.Add(img)
		}
	}

	// Refresh all containers
	ui.CenterNorth.Refresh()
	ui.CenterEast.Refresh()
	ui.CenterSouth.Refresh()
	ui.CenterWest.Refresh()
}

func (ui *GameUI) showComputerDecision(player *Player, text string, suit Suit) {
	ui.Window.Canvas().SetContent(ui.MainContent) // Ensure main content stays visible

	var pos *fyne.Container
	switch player.Name {
	case "NORTH":
		pos = ui.CenterNorth
	case "EAST":
		pos = ui.CenterEast
	case "WEST":
		pos = ui.CenterWest
	default:
		return
	}

	fullText := text
	if suit != Suit(-1) {
		fullText += " " + suit.FriendlySuit()
	}

	pos.Objects = []fyne.CanvasObject{
		widget.NewLabelWithStyle(fullText,
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true}),
	}
	pos.Refresh()
}

func (ui *GameUI) processComputerTrumpSelection() {
	if !ui.Round.SelectingTrump {
		return
	}

	// Only proceed if it's a computer player's turn
	for ui.Round.SelectingTrump && ui.Round.ActivePlayer != 2 { // 2 is human position
		player := ui.Round.Players[ui.Round.ActivePlayer]

		ui.showComputerDecision(player, "Thinking...", Suit(-1))
		time.Sleep(800 * time.Millisecond)

		var decision Call
		var suit Suit

		// Decision logic...
		if len(ui.Round.Deck.Cards) > 0 && ui.Round.Deck.Cards[0].FaceUp {
			suit = ui.Round.Deck.Cards[0].Suit
			decision = player.CallOrPass(suit, ui.Round.Dealer%2 == ui.Round.ActivePlayer%2)
		} else {
			passedSuit := Suit(-1)
			if len(ui.Round.Deck.Cards) > 0 {
				passedSuit = ui.Round.Deck.Cards[0].Suit
			}
			decision, suit = player.DeclareTrump(passedSuit)
		}

		ui.showComputerDecision(player, decision.FriendlyCall(), suit)
		time.Sleep(1 * time.Second)

		ui.Round.ComputerTrumpSelection(decision, suit)
		ui.RefreshUI()

		// Break if it's now human's turn
		if ui.Round.ActivePlayer == 2 {
			ui.showTrumpSelection()
			return
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

func (ui *GameUI) showComputerThinking(player *Player) {
	var pos *fyne.Container

	switch player.Name {
	case "NORTH":
		pos = ui.CenterNorth
	case "EAST":
		pos = ui.CenterEast
	case "WEST":
		pos = ui.CenterWest
	default:
		return
	}

	pos.Objects = []fyne.CanvasObject{
		widget.NewLabelWithStyle("Thinking...",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true}),
	}
	pos.Refresh()
}

func (ui *GameUI) playComputerTurn() {
	computer := ui.Round.Players[ui.Round.ActivePlayer]
	if !computer.IsPlaying {
		ui.Round.ActivePlayer = (ui.Round.ActivePlayer + 1) % 4
		ui.playComputerTurn() // Move to next player
		return
	}
	// Show thinking indicator
	ui.showComputerDecision(computer, "Playing...", Suit(-1))
	time.Sleep(500 * time.Millisecond)

	play := computer.BestPlay(ui.getCurrentTrick(), *ui.Round)
	playedCard := computer.PlayCard(&play)

	// Update trick display
	ui.Trick[ui.Round.ActivePlayer] = playedCard
	ui.updateTrickDisplay(ui.Trick)

	// Brief pause to see the play
	time.Sleep(1 * time.Second)

	// Move to next player
	ui.Round.ActivePlayer = (ui.Round.ActivePlayer + 1) % 4

	// If trick is complete, resolve it
	if ui.isTrickComplete() {
		winner := resolveTrick(ui.Trick[:], ui.Round)
		ui.Round.Lead = winner
		ui.Round.ActivePlayer = winner

		// Show completed trick for 1 second
		time.Sleep(1 * time.Second)

		// Clear the trick display
		ui.clearTrickDisplay()

		// Reset trick array
		ui.Trick = [4]*Card{}
	}

	// Continue with next player if computer
	if ui.Round.ActivePlayer != 2 {
		ui.playComputerTurn()
	}
}

func (ui *GameUI) getCurrentTrick() []*Card {
	var trick []*Card
	for _, c := range ui.Trick {
		if c != nil {
			trick = append(trick, c)
		}
	}
	return trick
}

func (ui *GameUI) isTrickComplete() bool {
	count := 0
	for _, c := range ui.Trick {
		if c != nil {
			count++
		}
	}
	return count == 4
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

func (ui *GameUI) runOnMainThread(f func()) {
	if ui.Window == nil || ui.Window.Canvas() == nil {
		return
	}
	ui.Window.Canvas().SetContent(container.NewWithoutLayout()) // Dummy refresh
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to ensure UI responsiveness
		runtime.LockOSThread()
		f()
		runtime.UnlockOSThread()
	}()
}

func (ui *GameUI) clearTrickDisplay() {
	ui.CenterNorth.Objects = nil
	ui.CenterEast.Objects = nil
	ui.CenterSouth.Objects = nil
	ui.CenterWest.Objects = nil

	// Refresh all containers
	ui.CenterNorth.Refresh()
	ui.CenterEast.Refresh()
	ui.CenterSouth.Refresh()
	ui.CenterWest.Refresh()

	// Force full UI update
	ui.Window.Content().Refresh()
}

func (ui *GameUI) updateCallerIndicator() {
	if ui.Round.Caller == nil || ui.Round.Trump == Suit(-1) {
		ui.CallerIndicator.SetText("")
		return
	}

	text := fmt.Sprintf("%s called %s",
		ui.Round.Caller.Name,
		ui.Round.Trump.FriendlySuit())

	if ui.Round.Alone {
		text += " (Alone!)"
	}

	ui.CallerIndicator.SetText(text)
	ui.CallerIndicator.Refresh()
}

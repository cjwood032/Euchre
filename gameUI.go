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

	Trick        []*Card
	NewGameBtn   *widget.Button
	SouthHandBox *fyne.Container
	BottomArea   *fyne.Container
}

func (ui *GameUI) RefreshUI() {
	// Refresh kitty and clear trick display
	kitty := createStackedKitty(ui.Round, fyne.NewSize(70, 110))
	ui.KittyContainer.Objects = []fyne.CanvasObject{kitty}
	ui.updateTrickDisplay(make([]*Card, 4))

	// Update scores
	ui.NorthScore.SetText(fmt.Sprintf("Score: %d", ui.Players[0].Score))
	ui.EastScore.SetText(fmt.Sprintf("Score: %d", ui.Players[1].Score))
	ui.SouthScore.SetText(fmt.Sprintf("Score: %d", ui.Players[2].Score))
	ui.WestScore.SetText(fmt.Sprintf("Score: %d", ui.Players[3].Score))

	// Always reset to main content first
	ui.Window.SetContent(ui.MainContent)

	if ui.Round.SelectingTrump {
		if ui.Round.ActivePlayer == 2 { // Human player's turn
			ui.showTrumpSelection()
		} else {
			// Process computer players' turns
			go func() {
				ui.processComputerTrumpSelection()
				// Refresh again after computers have made their decisions
				ui.RefreshUI()
			}()
		}
	} else {
		// Show normal game UI
		ui.updateHumanHand()
	}
}

func (ui *GameUI) updateHumanHand() {
	ui.HandBox.Objects = nil

	cardSize := fyne.NewSize(80, 120)
	player := ui.Players[2] // South is human player

	for _, card := range player.CardMap.ToSlice() {
		currentCard := card
		cardUI := container.NewVBox(
			renderCardImage(currentCard, cardSize),
		)

		// Only add Play button if not selecting trump
		if !ui.Round.SelectingTrump {
			cardUI.Add(widget.NewButton("Play", func() {
				ui.Trick[2] = player.PlayCard(currentCard)
				ui.updateTrickDisplay(ui.Trick)

				// Process AI turns
				for i := 1; i < 4; i++ {
					turn := (2 + i) % 4
					ai := ui.Round.Players[turn]
					if ai.ComputerPlayer {
						var partialTrick []*Card
						for _, t := range ui.Trick {
							if t != nil {
								partialTrick = append(partialTrick, t)
							}
						}
						play := ai.BestPlay(partialTrick, *ui.Round)
						ui.Trick[turn] = ai.PlayCard(&play)
						ui.updateTrickDisplay(ui.Trick)
					}
				}

				// Determine winner
				winner := resolveTrick(ui.Trick, ui.Round)
				ui.Round.Lead = winner

				time.Sleep(time.Second)
				ui.Trick = make([]*Card, 4)
				ui.updateTrickDisplay(ui.Trick)
				ui.updateHumanHand()
			}))
		}

		ui.HandBox.Add(cardUI)
	}
}

func (ui *GameUI) showTrumpSelection() {
	if !ui.Round.SelectingTrump || ui.Round.ActivePlayer != 2 {
		return
	}

	// Update hand to show cards without play buttons
	ui.updateHumanHand()

	firstRound := len(ui.Round.Deck.Cards) > 0 && ui.Round.Deck.Cards[0].FaceUp

	var trumpSelectionUI *fyne.Container
	if firstRound {
		topCard := ui.Round.Deck.Cards[0]
		trumpSelectionUI = container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Top card is %s of %s", topCard.FriendlyRank(), topCard.Suit.FriendlySuit())),
			widget.NewLabel("Do you want to:"),
			widget.NewButton("Order Up", func() {
				ui.Round.HumanTrumpSelection(OrderUp, topCard.Suit)
				ui.RefreshUI()
			}),
			widget.NewButton("Go Alone", func() {
				ui.Round.HumanTrumpSelection(Alone, topCard.Suit)
				ui.RefreshUI()
			}),
			widget.NewButton("Pass", func() {
				ui.Round.HumanTrumpSelection(Pass, topCard.Suit)
				ui.RefreshUI()
			}),
		)
	} else {
		passedSuit := Suit(-1)
		if len(ui.Round.Deck.Cards) > 0 {
			passedSuit = ui.Round.Deck.Cards[0].Suit
		}

		suitButtons := container.NewHBox()
		for _, suit := range []Suit{Spades, Diamonds, Clubs, Hearts} {
			if suit != passedSuit {
				currentSuit := suit
				suitButtons.Add(widget.NewButton(suit.FriendlySuit(), func() {
					ui.Round.HumanTrumpSelection(OrderUp, currentSuit)
					ui.RefreshUI()
				}))
			}
		}

		trumpSelectionUI = container.NewVBox(
			widget.NewLabel("Choose a trump suit:"),
			suitButtons,
			widget.NewButton("Pass", func() {
				ui.Round.HumanTrumpSelection(Pass, Suit(-1))
				ui.RefreshUI()
			}),
		)
	}

	// Create a new bottom area with the trump selection above the hand
	bottomWithTrump := container.NewVBox(
		container.NewCenter(trumpSelectionUI),
		ui.SouthHandBox,
		container.NewCenter(ui.NewGameBtn),
	)

	// Create a temporary main content with the trump selection
	tempContent := container.NewBorder(
		ui.MainContent.(*fyne.Container).Objects[0], // North
		bottomWithTrump, // Modified bottom
		ui.MainContent.(*fyne.Container).Objects[2], // West
		ui.MainContent.(*fyne.Container).Objects[3], // East
		ui.MainContent.(*fyne.Container).Objects[4], // Center
	)

	ui.Window.SetContent(tempContent)
}

func (ui *GameUI) updateTrickDisplay(trick []*Card) {
	ui.CenterNorth.Objects = nil
	ui.CenterEast.Objects = nil
	ui.CenterSouth.Objects = nil
	ui.CenterWest.Objects = nil

	for i, card := range trick {
		if card == nil {
			continue
		}

		cardSize := fyne.NewSize(80, 120)
		switch i {
		case 0: // North
			ui.CenterNorth.Add(renderCardImage(card, cardSize))
		case 1: // East
			ui.CenterEast.Add(renderCardImage(card, cardSize))
		case 2: // South
			ui.CenterSouth.Add(renderCardImage(card, cardSize))
		case 3: // West
			ui.CenterWest.Add(renderCardImage(card, cardSize))
		}
	}
}

func (ui *GameUI) showComputerDecision(player *Player, decision Call, suit Suit) {
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

	// Create and show the decision text
	var decisionText = string(decision.FriendlyCall())
	if suit != Suit(-1) {
		decisionText += " " + suit.FriendlySuit()
	}
	decisionLabel := widget.NewLabel(decisionText)
	decisionLabel.Alignment = fyne.TextAlignCenter
	position.Add(decisionLabel)

	// Temporarily update the score label to show thinking
	originalText := label.Text
	label.SetText("Thinking...")
	ui.Window.Content().Refresh()

	// Pause for visibility
	time.Sleep(1 * time.Second)

	// Restore original label
	label.SetText(originalText)
	position.Objects = nil
	ui.Window.Content().Refresh()
}

func (ui *GameUI) processComputerTrumpSelection() {
	for ui.Round.SelectingTrump {
		currentPlayer := ui.Round.Players[ui.Round.ActivePlayer]

		if currentPlayer.ComputerPlayer {
			suit := ui.Round.Deck.Cards[0].Suit
			decision := currentPlayer.CallOrPass(suit, ui.Round.Dealer%2 == ui.Round.ActivePlayer%2)

			// Show the computer's decision
			ui.showComputerDecision(currentPlayer, decision, suit)

			// Process the decision

			ui.Round.ComputerTrumpSelection(decision, suit)

			// Pause between turns
			time.Sleep(1 * time.Second)
		} else {
			// Human player's turn - break and let UI handle it
			break
		}
	}
}

package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Euchre")
	myWindow.SetPadded(true)

	// Initialize players
	players := []*Player{
		{Name: "NORTH", ComputerPlayer: true},
		{Name: "EAST", ComputerPlayer: true},
		{Name: "SOUTH"},
		{Name: "WEST", ComputerPlayer: true},
	}

	// Create game and initial round
	game := CreateEuchreGame(players)
	game.NewGame(false)
	round := game.Rounds[len(game.Rounds)-1]

	// UI elements
	trick := make([]*Card, 4)
	handBox := container.NewHBox()

	// Center card positions
	centerNorth := container.NewCenter()
	centerEast := container.NewCenter()
	centerSouth := container.NewCenter()
	centerWest := container.NewCenter()

	kittyContainer := container.NewCenter()
centerArea := container.NewGridWithColumns(3,
    container.NewGridWithRows(3,
        container.NewCenter(),  // NW (empty)
        centerWest,            // W
        container.NewCenter(), // SW (empty)
    ),
    container.NewGridWithRows(3,
        centerNorth,           // N
        kittyContainer,        // Center (kitty will go here)
        centerSouth,           // S
    ),
    container.NewGridWithRows(3,
        container.NewCenter(),  // NE (empty)
        centerEast,            // E
        container.NewCenter(), // SE (empty)
    ),
)

	// Player areas (just names and scores now)
	north := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("NORTH"), 
			widget.NewLabel(fmt.Sprintf("Score: %d", players[0].Score)),
		),
	)

	south := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("SOUTH"), 
			widget.NewLabel(fmt.Sprintf("Score: %d", players[2].Score)),
		),
	)

	east := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("EAST"),
			widget.NewLabel(fmt.Sprintf("Score: %d", players[1].Score)),
		),
	)

	west := container.NewVBox(
		container.NewHBox(
			widget.NewLabel("WEST"),
			widget.NewLabel(fmt.Sprintf("Score: %d", players[3].Score)),
		),
	)

	// Function to update trick display
	updateTrickDisplay := func(trick []*Card) {
		cardSize := fyne.NewSize(80, 120)
		centerNorth.Objects = nil
		centerEast.Objects = nil
		centerSouth.Objects = nil
		centerWest.Objects = nil

		for i, card := range trick {
			if card == nil {
				continue
			}

			switch i {
			case 0: // North
				centerNorth.Add(renderCardImage(card, cardSize))
			case 1: // East
				centerEast.Add(renderCardImage(card, cardSize))
			case 2: // South
				centerSouth.Add(renderCardImage(card, cardSize))
			case 3: // West
				centerWest.Add(renderCardImage(card, cardSize))
			}
		}
	}

	// Single updateHumanHand function
	var updateHumanHand func()

	updateHumanHand = func() {
	handBox.Objects = nil
	cardSize := fyne.NewSize(80, 120)
	player := players[2] // South is human player

	for _, card := range player.CardMap.ToSlice() {
		currentCard := card
		cardUI := container.NewVBox(
			renderCardImage(currentCard, cardSize),
			widget.NewButton("Play", func() {
				// Play the card (South is position 2)
				trick[2] = player.PlayCard(currentCard)
				updateTrickDisplay(trick)

				// Process AI turns
				for i := 1; i < 4; i++ {
					turn := (2 + i) % 4 // Rotate through positions
					ai := round.Players[turn]
					if ai.ComputerPlayer {
						var partialTrick []*Card
						for _, t := range trick {
							if t != nil {
								partialTrick = append(partialTrick, t)
							}
						}
						play := ai.BestPlay(partialTrick, *round)
						trick[turn] = ai.PlayCard(&play)
						updateTrickDisplay(trick)
					}
				}

				// Determine winner
				winner := resolveTrick(trick, round)
				round.Lead = winner

				time.Sleep(time.Second)
				for i := range trick {
					trick[i] = nil
				}
				updateTrickDisplay(trick)
				updateHumanHand() 
			}),
		)
		handBox.Add(cardUI)
	}
}


	// Function to refresh the entire game UI
	// Function to refresh the entire game UI
refreshGameUI := func() {
    // Refresh kitty - update the existing container
    kitty := createStackedKitty(round, fyne.NewSize(70, 110))
    kittyContainer.Objects = []fyne.CanvasObject{kitty}

    // Clear trick display
    updateTrickDisplay(make([]*Card, 4))

    // Update scores
    north.Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("Score: %d", players[0].Score))
    east.Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("Score: %d", players[1].Score))
    south.Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("Score: %d", players[2].Score))
    west.Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("Score: %d", players[3].Score))

    // Refresh human hand
    updateHumanHand()
}

	// New Game button
	newGameBtn := widget.NewButton("New Game", func() {
		game.NewGame(false)
		round = game.Rounds[len(game.Rounds)-1]
		refreshGameUI()
	})

	// Bottom area
	bottomArea := container.NewVBox(
		south,
		handBox,
		container.NewCenter(newGameBtn),
	)

	// Main content
	mainContent := container.NewBorder(
		container.NewCenter(north),  // Top
		bottomArea,                  // Bottom
		container.NewCenter(west),   // Left
		container.NewCenter(east),   // Right
		centerArea,                 // Center
	)

	// Initial UI setup
	refreshGameUI()

	myWindow.SetContent(mainContent)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}


func updateTrickDisplay(trick []*Card, trickBoxes [4]*fyne.Container, cardSize fyne.Size) {
	for i := 0; i < 4; i++ {
		trickBoxes[i].Objects = nil
		if i < len(trick) && trick[i] != nil {
			trickBoxes[i].Add(renderCardImage(trick[i], cardSize))
		}
	}
}

func resolveTrick(trick []*Card, round *Round) int {
	if len(trick) == 0 || trick[0] == nil {
		return 0
	}

	lead := round.Lead
	winningIndex := lead
	trump := round.Trump
	leadSuit := trick[0].Suit

	for i := 1; i < 4; i++ {
		pos := (lead + i) % 4
		if pos < len(trick) && trick[pos] != nil && 
		   trick[pos].Beats(trick[winningIndex], leadSuit, trump) {
			winningIndex = pos
		}
	}

	if winningIndex < len(round.Players) {
		round.Players[winningIndex].TricksWon++
	}
	return winningIndex
}

func renderCardImage(card *Card, size fyne.Size) *canvas.Image {
	suit := strings.ToLower(card.Suit.FriendlySuit())
	rank := fmt.Sprintf("%d", card.Rank)
	if card.Rank == 1 {
		rank = "ace"
	} else if card.Rank == 11 {
		rank = "jack"
	} else if card.Rank == 12 {
		rank = "queen"
	} else if card.Rank == 13 {
		rank = "king"
	}

	img := canvas.NewImageFromFile(fmt.Sprintf("cardimages/%s_of_%s.png", rank, suit))
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(size)
	return img
}

func renderCardBack(size fyne.Size) *canvas.Image {
	img := canvas.NewImageFromFile("cardimages/back.png")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(size)
	return img
}

func createStackedKitty(round *Round, size fyne.Size) *fyne.Container {
	stack := container.NewWithoutLayout()
	offset := float32(15)

	for i := 0; i < 3; i++ {
		back := renderCardBack(size)
		back.Resize(size)
		back.Move(fyne.NewPos(0, float32(i)*offset))
		stack.Add(back)
	}

	if len(round.Deck.Cards) > 0 && round.Deck.Cards[0].FaceUp {
		topCard := renderCardImage(round.Deck.Cards[0], size)
		topCard.Resize(size)
		topCard.Move(fyne.NewPos(0, float32(3)*offset))
		stack.Add(topCard)
	}

	return container.NewPadded(stack)
}

func updateKittyPosition(dealerIndex int, kitty *fyne.Container) {
	switch dealerIndex {
	case 0: kitty.Move(fyne.NewPos(350, 20))  // North
	case 1: kitty.Move(fyne.NewPos(650, 250)) // East
	case 2: kitty.Move(fyne.NewPos(350, 480)) // South
	case 3: kitty.Move(fyne.NewPos(50, 250))  // West
	}
}


func updateGameUI(round *Round, window fyne.Window) {
    // Rebuild the main game UI based on current game state
    if round.SelectingTrump {
        updateUIForTrumpSelection(round, window)
    } else {
        // Show normal game UI
        // ... (your existing game UI code)
    }
}


func updateUIForTrumpSelection(round *Round, window fyne.Window) {
    if !round.SelectingTrump || round.ActivePlayer != 2 { // Assuming human is always position 2
        return
    }

    // Check if we're in first round (top card face up) or second round
    firstRound := len(round.Deck.Cards) > 0 && round.Deck.Cards[0].FaceUp
    
    var content *fyne.Container
    if firstRound {
        topCard := round.Deck.Cards[0]
        content = container.NewVBox(
            widget.NewLabel(fmt.Sprintf("Top card is %s of %s", topCard.FriendlyRank(), topCard.Suit.FriendlySuit())),
            widget.NewLabel("Do you want to:"),
            widget.NewButton("Order Up", func() {
                round.HumanTrumpSelection(OrderUp, topCard.Suit)
                updateGameUI(round, window)
            }),
            widget.NewButton("Go Alone", func() {
                round.HumanTrumpSelection(Alone, topCard.Suit)
                updateGameUI(round, window)
            }),
            widget.NewButton("Pass", func() {
                round.HumanTrumpSelection(Pass, topCard.Suit)
                updateGameUI(round, window)
            }),
        )
    } else {
        // Second round - choose any suit except the turned-down one
        passedSuit := Suit(-1)
        if len(round.Deck.Cards) > 0 {
            passedSuit = round.Deck.Cards[0].Suit
        }
        
        suitButtons := container.NewHBox()
        for _, suit := range []Suit{Spades, Diamonds, Clubs, Hearts} {
            if suit != passedSuit {
                currentSuit := suit
                suitButtons.Add(widget.NewButton(suit.FriendlySuit(), func() {
                    round.HumanTrumpSelection(OrderUp, currentSuit)
                    updateGameUI(round, window)
                }))
            }
        }
        
        content = container.NewVBox(
            widget.NewLabel("Choose a trump suit:"),
            suitButtons,
            widget.NewButton("Pass", func() {
                round.HumanTrumpSelection(Pass, Suit(-1))
                updateGameUI(round, window)
            }),
        )
    }
    
    window.SetContent(content)
}
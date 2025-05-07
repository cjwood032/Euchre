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

func updateHumanHand(player *Player, trick []*Card, handBox *fyne.Container,
	trickBoxes [4]*fyne.Container, round *Round, window fyne.Window) {

	handBox.Objects = nil
	cardSize := fyne.NewSize(80, 120)

	for _, card := range player.CardMap.ToSlice() {
		currentCard := card
		cardUI := container.NewVBox(
			renderCardImage(currentCard, cardSize),
			widget.NewButton("Play", func() {
				trick[round.Lead] = player.PlayCard(currentCard)
				updateTrickDisplay(trick, trickBoxes, cardSize)

				// Process AI turns
				for i := 1; i < 4; i++ {
					turn := (round.Lead + i) % 4
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
						updateTrickDisplay(trick, trickBoxes, cardSize)
					}
				}

				winner := resolveTrick(trick, round)
				round.Lead = winner

				time.Sleep(time.Second)
				for i := range trick {
					trick[i] = nil
				}
				updateTrickDisplay(trick, trickBoxes, cardSize)
				updateHumanHand(player, trick, handBox, trickBoxes, round, window)
			}),
		)
		handBox.Add(cardUI)
	}
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Euchre")
	myWindow.SetPadded(true)

	// Initialize players
	players := []*Player{
		{Name: "North", ComputerPlayer: true},
		{Name: "East", ComputerPlayer: true},
		{Name: "South"},
		{Name: "West", ComputerPlayer: true},
	}

	// Create game and initial round
	game := CreateEuchreGame(players)
	game.NewGame(false)
	round := game.Rounds[len(game.Rounds)-1]

	// UI elements
	trick := make([]*Card, 4)
	handBox := container.NewHBox()
	trickBoxes := [4]*fyne.Container{
		container.NewHBox(), container.NewHBox(),
		container.NewHBox(), container.NewHBox(),
	}

	// Create stacked kitty
	kitty := createStackedKitty(round, fyne.NewSize(70, 110))
	updateKittyPosition(round.Dealer, kitty)

	// Player areas
	north := container.NewVBox(
		trickBoxes[0],
		container.NewHBox(
			widget.NewLabel("North"), 
			widget.NewLabel(fmt.Sprintf("Score: %d", players[0].Score)),
		),
	)

	south := container.NewVBox(
		trickBoxes[2],
		container.NewHBox(
			widget.NewLabel("South"), 
			widget.NewLabel(fmt.Sprintf("Score: %d", players[2].Score)),
		),
	)

	east := container.NewVBox(
		widget.NewLabel("East"),
		widget.NewLabel(fmt.Sprintf("Score: %d", players[1].Score)),
		trickBoxes[1],
	)

	west := container.NewVBox(
		widget.NewLabel("West"),
		widget.NewLabel(fmt.Sprintf("Score: %d", players[3].Score)),
		trickBoxes[3],
	)

	// New Game button
	newGameBtn := widget.NewButton("New Game", func() {
		game.NewGame(false)
		round = game.Rounds[len(game.Rounds)-1]
		updateKittyPosition(round.Dealer, kitty)
		handBox.Objects = nil
		for i := range trickBoxes {
			trickBoxes[i].Objects = nil
		}
		updateHumanHand(players[2], trick, handBox, trickBoxes, round, myWindow)
	})

	// Bottom area with hand and New Game button
	bottomArea := container.NewVBox(
		south,
		handBox,
		container.NewCenter(newGameBtn),
	)

	// Main layout
	content := container.NewBorder(
		container.NewCenter(north), // Top
		bottomArea,                // Bottom (now includes New Game button)
		container.NewCenter(west), // Left
		container.NewCenter(east), // Right
		kitty,                    // Center
	)

	// Initial hand update
	updateHumanHand(players[2], trick, handBox, trickBoxes, round, myWindow)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
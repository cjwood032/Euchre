package main

import (
	"fmt"
	"strings"

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
	handBox := container.NewHBox()
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
	currentRound := game.Rounds[len(game.Rounds)-1]

	// Initialize UI state
	ui := &GameUI{
		Window:  myWindow,
		Players: players,
		Round:   currentRound,
		Game:    game,
		Trick:   make([]*Card, 4),
		HandBox: container.NewHBox(),
	}

	// Initialize center positions
	ui.CenterNorth = container.NewCenter()
	ui.CenterEast = container.NewCenter()
	ui.CenterSouth = container.NewCenter()
	ui.CenterWest = container.NewCenter()
	ui.KittyContainer = container.NewCenter()

	// Create center area layout
	centerArea := container.NewGridWithColumns(3,
		container.NewGridWithRows(3,
			container.NewCenter(), // NW (empty)
			ui.CenterWest,         // W
			container.NewCenter(), // SW (empty)
		),
		container.NewGridWithRows(3,
			ui.CenterNorth,    // N
			ui.KittyContainer, // Center (kitty)
			ui.CenterSouth,    // S
		),
		container.NewGridWithRows(3,
			container.NewCenter(), // NE (empty)
			ui.CenterEast,         // E
			container.NewCenter(), // SE (empty)
		),
	)

	// Initialize score labels
	ui.NorthScore = widget.NewLabel(fmt.Sprintf("Score: %d", players[0].Score))
	ui.EastScore = widget.NewLabel(fmt.Sprintf("Score: %d", players[1].Score))
	ui.SouthScore = widget.NewLabel(fmt.Sprintf("Score: %d", players[2].Score))
	ui.WestScore = widget.NewLabel(fmt.Sprintf("Score: %d", players[3].Score))

	// Create player areas
	north := container.NewHBox(widget.NewLabel("NORTH"), ui.NorthScore)
	east := container.NewHBox(widget.NewLabel("EAST"), ui.EastScore)
	south := container.NewHBox(widget.NewLabel("SOUTH"), ui.SouthScore)
	west := container.NewHBox(widget.NewLabel("WEST"), ui.WestScore)

	// New Game button
	// New Game button
	newGameBtn := widget.NewButton("New Game", func() {
		ui.Game.NewGame(false)
		ui.Round = ui.Game.Rounds[len(ui.Game.Rounds)-1]
		ui.Trick = make([]*Card, 4)

		// Clear any existing UI state
		ui.CenterNorth.Objects = nil
		ui.CenterEast.Objects = nil
		ui.CenterSouth.Objects = nil
		ui.CenterWest.Objects = nil

		// Force refresh
		ui.RefreshUI()
	})

	// Store references in the UI struct
	ui.NewGameBtn = newGameBtn
	ui.SouthHandBox = handBox
	ui.BottomArea = container.NewVBox(
		container.NewCenter(south),
		handBox,
		container.NewCenter(newGameBtn),
	)

	// Bottom area
	bottomArea := container.NewVBox(
		container.NewCenter(south),
		ui.HandBox,
		container.NewCenter(newGameBtn),
	)

	// Main content
	ui.MainContent = container.NewBorder(
		container.NewCenter(north), // Top
		bottomArea,                 // Bottom
		container.NewCenter(west),  // Left
		container.NewCenter(east),  // Right
		centerArea,                 // Center
	)

	// Initial UI setup
	ui.RefreshUI()

	myWindow.SetContent(ui.MainContent)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
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
	offset := float32(8)

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

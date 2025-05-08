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
		{Name: "NORTH", ComputerPlayer: true, Position: 0, IsPlaying: true},
		{Name: "EAST", ComputerPlayer: true, Position: 1, IsPlaying: true},
		{Name: "SOUTH", Position: 2, IsPlaying: true},
		{Name: "WEST", ComputerPlayer: true, Position: 3, IsPlaying: true},
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
		Trick:   [4]*Card(make([]*Card, 4)),
		HandBox: container.NewHBox(),
	}
	callerIndicator := widget.NewLabel("")
	callerIndicator.Alignment = fyne.TextAlignCenter
	callerIndicator.TextStyle = fyne.TextStyle{Bold: true}

	ui.CallerIndicator = callerIndicator

	// Add it to your layout (modify your container as needed)

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

	// New Game button
	newGameBtn := widget.NewButton("New Game", func() {
		ui.Game.NewGame(false)
		ui.Round = ui.Game.Rounds[len(ui.Game.Rounds)-1]
		ui.Trick = [4]*Card(make([]*Card, 4))

		// Complete UI reset
		ui.HandBox.Objects = []fyne.CanvasObject{container.NewHBox()}
		ui.CenterNorth.Objects = nil
		ui.CenterEast.Objects = nil
		ui.CenterSouth.Objects = nil
		ui.CenterWest.Objects = nil
		ui.KittyContainer.Objects = nil

		// Force complete rebuild
		ui.Window.SetContent(ui.MainContent)
		ui.RefreshUI()
	})
	// Create player areas with dealer indicator
	north := container.NewHBox(
		widget.NewLabel("NORTH"),
		ui.createDealerIndicator(0),
		ui.NorthScore,
		callerIndicator,
	)
	east := container.NewHBox(
		widget.NewLabel("EAST"),
		ui.createDealerIndicator(1),
		ui.EastScore,
		callerIndicator,
	)
	south := container.NewHBox(
		widget.NewLabel("SOUTH"),
		ui.createDealerIndicator(2),
		ui.SouthScore,
		callerIndicator,
	)
	west := container.NewHBox(
		widget.NewLabel("WEST"),
		ui.createDealerIndicator(3),
		ui.WestScore,
		callerIndicator,
	)
	// Store references in the UI struct
	ui.NewGameBtn = newGameBtn
	ui.SouthHandBox = handBox
	// Create a centered container for the player's hand
	handContainer := container.NewHBox()
	ui.HandBox = container.NewCenter(handContainer) // Now centered

	// Bottom area with centered hand
	// Create controls section (topmost)
	// Create controls section at the very top
	controls := container.NewCenter(newGameBtn)

	// Main content with clear hierarchy
	ui.MainContent = container.NewBorder(
		container.NewVBox( // Top section
			controls, // New Game button at very top
			container.NewCenter(north),
			//callerIndicator,
		),
		container.NewVBox( // Bottom section
			container.NewCenter(south),
			ui.HandBox,
			//callerIndicator, // Centered hand (only one instance)
		),
		container.NewCenter(west), // Left
		container.NewCenter(east), // Right
		centerArea,                // Center
	)

	// Initial UI setup
	ui.RefreshUI()

	myWindow.SetContent(ui.MainContent)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func resolveTrick(trick []*Card, round *Round) int {
	winningIndex := round.Lead
	winningCard := trick[winningIndex]
	leadSuit := trick[round.Lead].Suit

	for i := 1; i < 4; i++ {
		pos := (round.Lead + i) % 4
		if trick[pos].Beats(winningCard, leadSuit, round.Trump) {
			winningIndex = pos
			winningCard = trick[pos]
		}
	}

	// Update tricks won for active players only
	if round.Players[winningIndex].IsPlaying {
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
	// Hide kitty when trump has been selected and play has begun
	if !round.SelectingTrump {
		return container.NewPadded() // Empty container when play begins
	}

	// Only show kitty during trump selection phase
	stack := container.NewWithoutLayout()
	offset := float32(8)

	// Always show 3 face-down cards
	for i := 0; i < 3; i++ {
		back := renderCardBack(size)
		back.Resize(size)
		back.Move(fyne.NewPos(0, float32(i)*offset))
		stack.Add(back)
	}

	// Show top card only if it exists and is face up
	if len(round.Deck.Cards) > 0 && round.Deck.Cards[0].FaceUp {
		topCard := renderCardImage(round.Deck.Cards[0], size)
		topCard.Resize(size)
		topCard.Move(fyne.NewPos(0, float32(3)*offset))
		stack.Add(topCard)
	} else if len(round.Deck.Cards) > 0 {
		// Show face-down card if exists
		back := renderCardBack(size)
		back.Resize(size)
		back.Move(fyne.NewPos(0, float32(3)*offset))
		stack.Add(back)
	}

	return container.NewPadded(stack)
}

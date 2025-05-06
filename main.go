package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func cardImagePath(card *Card) string {
	suit := strings.ToLower(card.Suit.FriendlySuit())
	suit = strings.ReplaceAll(suit, " ", "_")
	var rank string
	switch card.Rank {
	case 1:
		rank = "ace"
	case 11:
		rank = "jack"
	case 12:
		rank = "queen"
	case 13:
		rank = "king"
	default:
		rank = fmt.Sprintf("%d", card.Rank)
	}
	return fmt.Sprintf("cardimages/%s_of_%s.png", rank, suit)
}

func renderCardImage(card *Card, size fyne.Size) fyne.CanvasObject {
	path := cardImagePath(card)
	img := canvas.NewImageFromFile(path)
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(size)
	return img
}

func renderCardBack(size fyne.Size) fyne.CanvasObject {
	img := canvas.NewImageFromFile("cardimages/back.png")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(size)
	return img
}

func renderCardImageWithButton(card *Card, size fyne.Size, onClick func()) fyne.CanvasObject {
	img := renderCardImage(card, size)
	btn := widget.NewButton("Play", onClick)
	return container.NewVBox(img, btn)
}

func updateTrickDisplay(trick []*Card, trickBoxes [4]*fyne.Container, cardSize fyne.Size) {
	for i, card := range trick {
		trickBoxes[i].Objects = nil
		if card != nil {
			trickBoxes[i].Add(renderCardImage(card, cardSize))
		}
		//trickBoxes[i].Refresh()
	}
}

func resolveTrick(trick []*Card, round *Round) int {
	lead := round.Lead
	winner := round.Players[lead]
	winningIndex := lead
	for i := 1; i < 4; i++ {
		pos := (lead + i) % 4
		if trick[pos] != nil && trick[pos].Beats(trick[winningIndex], trick[lead].Suit, round.Trump) {
			winningIndex = pos
			winner = round.Players[pos]
		}
	}
	winner.TricksWon++
	return winningIndex
}

func updateHumanHand(player *Player, trick []*Card, handBox *fyne.Container, trickBoxes [4]*fyne.Container, round *Round, window fyne.Window) {
	handBox.Objects = nil
	hand := player.CardMap.ToSlice()
	cardSize := fyne.NewSize(64, 96)

	for _, card := range hand {
		c := card
		cardUI := renderCardImageWithButton(c, cardSize, func() {
			trick[round.Lead] = player.PlayCard(c)
			updateTrickDisplay(trick, trickBoxes, cardSize)

			for i := 1; i < 4; i++ {
				turn := (round.Lead + i) % 4
				ai := round.Players[turn]
				if ai.ComputerPlayer {
					var partialTrick []*Card
					for _, t := range trick {
						if t != nil {
							card := *t
							partialTrick = append(partialTrick, &card)
						}
					}
					play := ai.BestPlay(partialTrick, *round)
					trick[turn] = ai.PlayCard(&play)
					updateTrickDisplay(trick, trickBoxes, cardSize)
				}
			}

			winner := resolveTrick(trick, round)
			round.Lead = winner

			go func() {
				time.Sleep(1 * time.Second)
				for i := range trick {
					trick[i] = nil
				}
				updateTrickDisplay(trick, trickBoxes, cardSize)
				updateHumanHand(player, trick, handBox, trickBoxes, round, window)
			}()
		})
		handBox.Add(cardUI)
	}

	//handBox.Refresh()
}

func playerLabel(player *Player, isDealer bool) fyne.CanvasObject {
	label := fmt.Sprintf("%s (Score: %d, Tricks: %d)", player.Name, player.Score, player.TricksWon)
	if isDealer {
		label += " (Dealer)"
	}
	return widget.NewLabel(label)
}

func main() {

	var controlBox *fyne.Container
	var callButtons *fyne.Container

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Working directory:", wd)

	myApp := app.New()
	myWindow := myApp.NewWindow("Euchre")

	players := []*Player{
		{Name: "North"},
		{Name: "East", ComputerPlayer: true},
		{Name: "South"},
		{Name: "West", ComputerPlayer: true},
	}
	game := CreateEuchreGame(players)
	game.NewGame(false)
	round := game.Rounds[len(game.Rounds)-1]
	dealerIndex := round.Dealer

	trick := make([]*Card, 4)
	handBox := container.NewHBox()
	trickBoxes := [4]*fyne.Container{
		container.NewHBox(),
		container.NewHBox(),
		container.NewHBox(),
		container.NewHBox(),
	}
	cardSize := fyne.NewSize(64, 96)

	human := players[2]
	updateHumanHand(human, trick, handBox, trickBoxes, round, myWindow)

	kittyBox := container.NewHBox()
	if round.Deck.Cards[0].FaceUp {
		topCard := round.Deck.Cards[len(round.Deck.Cards)-1]
		kittyBox.Add(renderCardImage(topCard, cardSize))
	}

	top := container.NewVBox(
		container.NewCenter(playerLabel(players[0], dealerIndex == 0)),
		container.NewCenter(trickBoxes[0]),
		container.NewCenter(kittyBox),
	)

	middle := container.NewHBox(
		container.NewVBox(
			playerLabel(players[1], dealerIndex == 1),
			trickBoxes[1],
		),
		layout.NewSpacer(),
		container.NewVBox(
			playerLabel(players[3], dealerIndex == 3),
			trickBoxes[3],
		),
	)

	bottom := container.NewVBox(
		container.NewCenter(trickBoxes[2]),
		container.NewCenter(playerLabel(players[2], dealerIndex == 2)),
		handBox,
	)

	mainContent := container.NewVBox(
		top,
		layout.NewSpacer(),
		middle,
		layout.NewSpacer(),
		bottom,
	)

	controlBox = container.NewHBox()
	newGameBtn := widget.NewButton("New Game", func() {
		
		game = CreateEuchreGame(players)
		
		game.NewGame(false)  
		round = game.Rounds[len(game.Rounds)-1]
		dealerIndex = round.Dealer
		
		handBox.Objects = nil
		for i := range trickBoxes {
			trickBoxes[i].Objects = nil
		}
		kittyBox.Objects = nil
		
		// Show the top card if face up
		if len(round.Deck.Cards) > 0 && round.Deck.Cards[0].FaceUp {
			topCard := round.Deck.Cards[0]  // Fixed index (was using len-1)
			kittyBox.Add(renderCardImage(topCard, cardSize))
		}
	
		// Update UI with human player's hand
		updateHumanHand(human, trick, handBox, trickBoxes, round, myWindow)
		
		// Start the bidding process
		passedSuit := round.Deck.Cards[0].Suit  // Fixed index
		passed := false
		activePlayer := (round.Dealer + 1) % 4  // Start with player left of dealer
		
		for {
			player := round.Players[activePlayer]
			if player.ComputerPlayer {
				if !passed {
					action := player.CallOrPass(passedSuit, activePlayer%2 == round.Dealer%2)
					if action != Pass {
						round.Trump = passedSuit
						round.Caller = player
						round.Alone = action == Alone
						break
					}
				} else {
					call, suit := player.DeclareTrump(passedSuit)
					if call != Pass {
						round.Trump = suit
						round.Caller = player
						round.Alone = call == Alone
						break
					}
				}
				activePlayer = (activePlayer + 1) % 4
				if activePlayer == round.Dealer && !passed {
					passed = true
					round.Deck.Cards[0].FaceUp = false
				}
			} else {
				// Human turn to decide
				callButtons.Show()
				break
			}
		}
	})
	controlBox.Add(newGameBtn)

	callButtons = container.NewHBox(
		widget.NewButton("Pass", func() {
			fmt.Println("Player passed")
		}),
		widget.NewButton("Order", func() {
			fmt.Println("Player ordered trump")
		}),
		widget.NewButton("Alone", func() {
			fmt.Println("Player ordered alone")
		}),
	)
	callButtons.Hide()
	controlBox.Add(callButtons)

	ui := container.NewVBox(
		mainContent,
		controlBox,
	)

	myWindow.SetContent(ui)
	myWindow.Resize(fyne.NewSize(1000, 700))
	myWindow.ShowAndRun()
}

package main

import (
	"context"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
	"math/rand"
	"os"
	"time"
)

func main() {
	term, err := tcell.New()
	if err != nil {
		panic(err)
	}
	defer term.Close()

	ctx, cancel := context.WithCancel(context.Background())

	gameWidget, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}
	resourceWidget, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		panic(err)
	}

	cont, err := container.New(term,
		container.Border(linestyle.Light),
		container.BorderTitle("Hamurabi -- Press Q to Quit"),
		container.BorderTitleAlignCenter(),
		container.SplitHorizontal(
			container.Top(
				container.Focused(),
				container.Border(linestyle.Light),
				container.BorderTitle("Gameplay"),
				container.PlaceWidget(gameWidget),
			),
			container.Bottom(
				container.Border(linestyle.Light),
				container.BorderTitle("Resources"),
				container.PlaceWidget(resourceWidget),
			),
			container.SplitPercent(85)),
	)
	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, term, cont, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	} else {
		gameLoop()
	}
}

func gameLoop() {
	rand.Seed(time.Now().UnixNano())
	var gameTurns int
	var test bool

	if os.Args[0] == "-test" {
		gameTurns = 100
		test = true
	} else {
		// gameTurns = playerInput("How many turns would you like to play?", 10, math.MaxInt, "")
		gameTurns = 10
		test = false
	}
	sumer := initCityState(gameTurns)

	if !test {
		orientation()
		for t := 0; t <= sumer.turns; t++ {
			printYearResults(sumer)
			getAcres(sumer)
			technology(sumer)
			feedPeople(sumer)
			agriculture(sumer)
		}
		endOfReign(sumer)
		os.Exit(0)
	}
	for t := 0; t <= sumer.turns; t++ {
		printYearResults(sumer)
		getAcres(sumer)
		technology(sumer)
		feedPeople(sumer)
		agriculture(sumer)
	}
	endOfReign(sumer)
}

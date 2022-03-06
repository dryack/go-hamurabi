package main

import (
	"math/rand"
	"time"
)

func main() {
	gameLoop()
}

func gameLoop() {
	rand.Seed(time.Now().UnixNano())
	term := initTerminal()
	turns := getGameTurns()
	sumer := newGameSession(turns)

	orientation()
	for t := 0; t <= sumer.turns; t++ {
		sumer.printYearResults(term)
		sumer.getAcres(term)
		sumer.construction(term)
		sumer.technology(term)
		sumer.feedPeople(term)
		sumer.agriculture(term)
	}
	sumer.endOfReign()
}

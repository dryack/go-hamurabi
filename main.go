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
	sumer := newGameSession()

	orientation()
	for t := 0; t <= sumer.turns; t++ {
		sumer.printYearResults(term)
		sumer.getAcres()
		sumer.construction()
		sumer.technology()
		sumer.feedPeople()
		sumer.agriculture()
	}
	sumer.endOfReign()
}

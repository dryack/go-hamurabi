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
	sumer := newGameSession()

	orientation()
	for t := 0; t <= sumer.turns; t++ {
		sumer.printYearResults()
		sumer.getAcres()
		sumer.construction()
		sumer.technology()
		sumer.feedPeople()
		sumer.agriculture()
	}
	sumer.endOfReign()
}

package main

import (
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	gameTurns := playerInput("How many turns would you like to play?", 10, math.MaxInt, "")
	sumer := initCityState(gameTurns)
	orientation()

	for t := 0; t <= sumer.turns; t++ {
		printYearResults(sumer)
		getAcres(sumer)
		technology(sumer)
		feedPeople(sumer)
		agriculture(sumer)
	}
	endOfReign(sumer)
}

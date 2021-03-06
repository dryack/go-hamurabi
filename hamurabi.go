package main

import (
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	gameLoop()
}

func gameLoop() {
	rand.Seed(time.Now().UnixNano())
	var gameTurns int
	var test bool

	if os.Args[0] == "-test" {
		gameTurns = 100
		test = true
	} else {
		gameTurns = playerInput("How many turns would you like to play?", 10, math.MaxInt, "")

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

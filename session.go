package main

import (
	"github.com/muesli/termenv"
	"math"
	"os"
)

type gameSession struct {
	turns           int
	state           cityState
	avgBushelsAvail int
	avgPestEaten    int
	avgStarved      int
	totalDead       int
	totAcresWasted  int
	otherCityStates []string
	points          int // score for evaluation
	palaceBuilding  int // what type of palace is under construction
	p               termenv.Profile
}

func newGameSession() (*gameSession, bool) {
	var gameTurns int
	var test bool

	if os.Args[0] == "-test" {
		gameTurns = 100
		test = true
	} else {
		gameTurns = playerInput("How many turns would you like to play?", 10, math.MaxInt, "", "chose")
		test = false
	}
	state := initCityState()

	return &gameSession{
		turns:           gameTurns,
		state:           *state,
		avgBushelsAvail: 0,
		avgPestEaten:    0,
		avgStarved:      0,
		totalDead:       0,
		totAcresWasted:  0,
		otherCityStates: []string{"Dūr-Katlimmu", "Aššur", "Uruk", "Akshak", "Ur", "Nippur", "Lagash", "Larak"},
		points:          0,
		p:               termenv.EnvColorProfile(),
	}, test
}

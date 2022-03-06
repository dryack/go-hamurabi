package main

import (
	"math"
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
	palaceBuilding  int // t type of palace is under construction
}

func getGameTurns() int {
	return playerInput("How many turns would you like to play?", 10, math.MaxInt, "", "chose")
}

func newGameSession(turns int) *gameSession {
	state := initCityState()

	return &gameSession{
		turns:           turns,
		state:           *state,
		avgBushelsAvail: 0,
		avgPestEaten:    0,
		avgStarved:      0,
		totalDead:       0,
		totAcresWasted:  0,
		otherCityStates: []string{"Dūr-Katlimmu", "Aššur", "Uruk", "Akshak", "Ur", "Nippur", "Lagash", "Larak"},
		points:          0,
	}
}

package main

import "math/rand"

type cityState struct {
	turns           int
	year            int
	population      int
	starved         int
	migrated        int
	bushels         int
	acres           int
	bYield          int
	pests           int
	tradeVal        int
	avgStarved      int
	totalDead       int
	popFed          int
	planted         int
	died            int
	born            int
	granary         int
	plows           int
	cows            int
	cowMultiplier   int
	nonFarmer       int
	tradeGoods      int
	acresWastage    int
	avgBushelsAvail int
	avgPestEaten    int
}

func initCityState(gameTurns int) *cityState {
	res := cityState{
		turns:         gameTurns,
		year:          0,
		population:    100,
		starved:       0,
		migrated:      0,
		bushels:       2800,
		acres:         1000,
		bYield:        3,
		pests:         200,
		tradeVal:      17 + rand.Intn(10),
		avgStarved:    0,
		totalDead:     0,
		popFed:        100,
		planted:       0,
		died:          0,
		granary:       0,
		plows:         0,
		cows:          0,
		cowMultiplier: 15,
	}

	return &res
}

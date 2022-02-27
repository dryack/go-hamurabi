package main

import "math/rand"

type cityState struct {
	year             int
	population       int
	starved          int
	migrated         int
	bushels          int
	acres            int
	bYield           int
	pests            int
	tradeVal         int
	popFed           int
	planted          int
	died             int
	born             int
	granary          int
	plows            int
	cows             int
	cowMultiplier    int
	nonFarmer        int
	tradeGoods       int
	forceSlaughtered int // cows that were forcibly slaughtered due to lack of land
	cowsFed          int // how many citizens fed by cows
	acresWastage     int
	stelae           int
	buildingPalace   int
	palace1          bool
	palace2          bool
	palace3          bool
}

func initCityState() *cityState {
	res := cityState{
		year:           0,
		population:     100,
		starved:        0,
		migrated:       0,
		bushels:        2800,
		acres:          1000,
		bYield:         3,
		pests:          200,
		tradeVal:       17 + rand.Intn(10),
		popFed:         100,
		planted:        0,
		died:           0,
		granary:        0,
		plows:          0,
		cows:           0,
		cowMultiplier:  15,
		stelae:         0,
		buildingPalace: -1,
		palace1:        false,
		palace2:        false,
		palace3:        false,
	}

	return &res
}

package main

import "math/rand"

type cityState struct {
	year             int
	starved          int
	migrated         int
	bYield           int
	pests            int
	tradeVal         int
	popFed           int
	planted          int
	died             int
	born             int
	nonFarmer        int
	tradeGoods       int
	forceSlaughtered int // cows that were forcibly slaughtered due to lack of land
	cowsFed          int // how many citizens fed by cows
	acresWastage     int
	buildingPalace   int
	resources        stateResources
	structures       stateStructures
	technology       technology
}

type stateResources struct {
	bushels    int
	acres      int
	population int
	granary    int
	plows      int
	cows       int
}

type stateStructures struct {
	palace1 bool
	palace2 bool
	palace3 bool
	stelae  int
}

type technology struct {
	cowMultiplier int
	stela         bool
	silver        bool
	ziggurat      bool
}

func initCityState() *cityState {
	tech := technology{
		cowMultiplier: 15,
		stela:         false,
		silver:        false,
		ziggurat:      false,
	}

	resources := stateResources{
		bushels:    2800,
		acres:      1000,
		population: 100,
		granary:    0,
		plows:      0,
		cows:       0,
	}

	structures := stateStructures{
		palace1: false,
		palace2: false,
		palace3: false,
		stelae:  0,
	}

	res := cityState{
		year:           0,
		starved:        0,
		migrated:       0,
		bYield:         3,
		pests:          200,
		tradeVal:       17 + rand.Intn(10),
		popFed:         100,
		planted:        0,
		died:           0,
		buildingPalace: -1,
		resources:      resources,
		technology:     tech,
		structures:     structures,
	}

	return &res
}

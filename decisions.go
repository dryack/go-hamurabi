package main

import (
	"sort"
	"strconv"
)

func (s *gameSession) getAcres() {
	bushelChange := 0
	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.bushels) + " to use for purchase!"

	res := playerInput("How many acres do you wish to buy?", 0, s.state.bushels/s.state.tradeVal, failMsg)
	if res == 0 {
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.acres) + " acres to sell!"
		res = playerInput("How many acres do you wish to sell?", 0, s.state.acres, failMsg) * -1
		if res == 0 {
			return
		}
		bushelChange = res * -s.state.tradeVal // forcing a positive value
	} else {
		bushelChange = res * -s.state.tradeVal
	}
	s.state.acres += res
	s.state.bushels += bushelChange
	s.grainRemaining(res)
}

func (s *gameSession) feedPeople() {
	cowFeedMultiplier := 30

	var ary = []int{s.state.bushels, s.state.population * 20}
	sort.Ints(ary)
	reqBushels := ary[0] // accept the lowest between bushels and bushels needed for the population

	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.bushels) + " available!"

	res := playerInput("How many bushels do you wish to release to your people?", reqBushels, s.state.bushels, failMsg)
	s.state.bushels -= res
	s.state.popFed = int(float64(res / 20))
	s.grainRemaining(res)

	failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.cows) + " cows to slaughter!"
	res = playerInput("How many cows would you like to slaughter in order to feed 35 people?", 0, s.state.cows, failMsg)
	s.state.cows -= res
	s.state.popFed += cowFeedMultiplier * res
}

func (s *gameSession) agriculture() {
	cowCost := 1000
	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.bushels) + " bushels available!"
	maxCows := s.state.bushels / cowCost
	res := playerInput("How many cows will you purchase, at "+strconv.Itoa(cowCost)+" bushels per cow?", 0, maxCows, failMsg)
	s.state.cows += res
	s.state.bushels -= res * cowCost
	s.grainRemaining(res)

	// only 1/3rd of the population can benefit from the plows
	var plowAry = []int{s.state.plows, s.state.population / 3}
	sort.Ints(plowAry)
	effectivePlows := plowAry[0]
	ableToPlant := ((s.state.population - effectivePlows) * 10) + (15 * effectivePlows)

	// fmt.Printf("\tpopReqForPlanting: %d\n", ableToPlant) // DEBUG

	var ary = []int{s.state.bushels, ableToPlant, s.state.acres - (s.state.cows * 3)}
	sort.Ints(ary)
	maxPlantable := ary[0]
	// avoid cows forcing us to negative numbers
	if maxPlantable < 0 {
		maxPlantable = 0
	}

	switch maxPlantable {
	case s.state.bushels:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.bushels) + " available!"
	case s.state.popFed * ableToPlant:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.population) + " people to plant the fields!"
	case s.state.acres:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.acres) + " acres to plant!"
	}
	res = playerInput("How many fields will you plant?", maxPlantable, maxPlantable, failMsg)
	if ableToPlant > res {
		s.state.nonFarmer = s.state.population - (res-(effectivePlows*15))/10
	}
	// fmt.Printf("\tNon-farmers: %d\n", s.state.nonFarmer) // DEBUG

	s.state.bushels -= res
	s.state.planted = res
	s.grainRemaining(res)
}

func (s *gameSession) technology() {
	costPlow := 100

	maxPlows := s.state.bushels / costPlow
	failMsg := "Think again Hamurabi, you only have enough to purchase " + strconv.Itoa(maxPlows) + " plows!"
	res := playerInput("Do you wish to order the purchase of plows for 100 bushels, these will make it easier "+
		"for your people to plant the fields?", 0, maxPlows, failMsg)
	s.state.plows += res
	s.state.bushels -= res * costPlow
	s.grainRemaining(res)
}

func (s *gameSession) construction() {
	costGranary := 1000

	if yn("My lord, do you wish to consider construction projects this year") {
		maxGranaries := s.state.bushels / costGranary
		failMsg := "Think again Hamurabi, you only have enough to purchase " + strconv.Itoa(maxGranaries) + " granaries!"
		res := playerInput("Do you wish to order the construction of city granaries for 1000 bushels, each are able "+
			"to protect a large amount of precious barley?", 0, maxGranaries, failMsg)
		s.state.granary += res
		s.state.bushels -= res * costGranary
		s.grainRemaining(res)
	}
}

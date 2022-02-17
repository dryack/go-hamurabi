package main

import (
	"sort"
	"strconv"
)

func getAcres(state *cityState) {
	bushelChange := 0
	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(state.bushels) + " to use for purchase!"

	res := playerInput("How many acres do you wish to buy?", 0, state.bushels/state.tradeVal, failMsg)
	if res == 0 {
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(state.acres) + " acres to sell!"
		res = playerInput("How many acres do you wish to sell?", 0, state.bushels, failMsg) * -1
		if res == 0 {
			return
		}
		bushelChange = res * -state.tradeVal // forcing a positive value
	} else {
		bushelChange = res * -state.tradeVal
	}
	state.acres += res
	state.bushels += bushelChange
	grainRemaining(res, state)
}

func feedPeople(state *cityState) {
	var ary = []int{state.bushels, state.population * 20}
	sort.Ints(ary)
	reqBushels := ary[0] // accept the lowest between bushels and bushels needed for the population

	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(state.bushels) + " available!"

	res := playerInput("How many bushels do you wish to release to your people?", reqBushels, state.bushels, failMsg)
	state.bushels -= res
	state.popFed = int(float64(res / 20))
	grainRemaining(res, state)

	failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(state.cows) + " cows to slaughter!"
	res = playerInput("Would you like to slaughter cows to feed 35 people?", 0, state.cows, failMsg)
	state.cows -= res
	state.popFed += 30 * res
}

func agriculture(state *cityState) {
	cowCost := 300
	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(state.bushels) + " bushels available!"
	maxCows := state.bushels / cowCost
	res := playerInput("How many cows will you purchase, at 300 bushels per cow?", 0, maxCows, failMsg)
	state.cows += res
	state.bushels -= res * cowCost
	grainRemaining(res, state)

	// only 1/3rd of the population can benefit from the plows
	var plowAry = []int{state.plows, state.population / 3}
	sort.Ints(plowAry)
	effectivePlows := plowAry[0]
	ableToPlant := ((state.population - effectivePlows) * 10) + (15 * effectivePlows)

	// fmt.Printf("\tpopReqForPlanting: %d\n", ableToPlant) // DEBUG

	var ary = []int{state.bushels, state.population * ableToPlant, state.acres - (state.cows * 3)}
	sort.Ints(ary)
	maxPlantable := ary[0]

	switch maxPlantable {
	case state.bushels:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(state.bushels) + " available!"
	case state.popFed * ableToPlant:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(state.population) + " people to plant the fields!"
	case state.acres:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(state.acres) + " acres to plant!"
	}
	res = playerInput("How many fields will you plant?", maxPlantable, maxPlantable, failMsg)

	state.bushels -= res
	state.planted = res
	grainRemaining(res, state)
}

func technology(state *cityState) {
	costGranary := 1000
	maxGranaries := state.bushels / costGranary
	failMsg := "Think again Hamurabi, you only have enough to purchase " + strconv.Itoa(maxGranaries) + " granaries!"
	res := playerInput("Do you wish to order the construction of city granaries for 1000 bushels, each are able "+
		"to protect a large amount of precious barley?", 0, maxGranaries, failMsg)
	state.granary += res
	state.bushels -= res * costGranary
	grainRemaining(res, state)

	costPlow := 100
	maxPlows := state.bushels / costPlow
	failMsg = "Think again Hamurabi, you only have enough to purchase " + strconv.Itoa(maxPlows) + " plows!"
	res = playerInput("Do you wish to order the purchase of plows for 100 bushels, these will make it easier "+
		"for your people to plant the fields?", 0, maxPlows, failMsg)
	state.plows += res
	state.bushels -= res * costPlow
	grainRemaining(res, state)
}

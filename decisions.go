package main

import (
	"fmt"
	"sort"
	"strconv"
)

func (s *gameSession) getAcres(term *terminal) {
	bushelChange := 0
	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.bushels) + " to use for purchase!"

	res := playerInput("How many acres do you wish to buy?", 0, s.state.bushels/s.state.tradeVal, failMsg, "bought")
	if res == 0 {
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.acres) + " acres to sell!"
		res = playerInput("How many acres do you wish to sell?", 0, s.state.acres, failMsg, "sold") * -1
		if res == 0 {
			return
		}
		bushelChange = res * -s.state.tradeVal // forcing a positive value
	} else {
		bushelChange = res * -s.state.tradeVal
	}
	s.state.acres += res
	s.state.bushels += bushelChange
	term.grainRemaining(s.state.bushels, res)
}

func (s *gameSession) feedPeople(term *terminal) {
	cowFeedMultiplier := 45

	var ary = []int{s.state.bushels, s.state.population * 20}
	sort.Ints(ary)
	reqBushels := ary[0] // accept the lowest between bushels and bushels needed for the population

	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.bushels) + " available!"

	res := playerInput("How many bushels do you wish to release to your people?", reqBushels, s.state.bushels, failMsg, "released")
	s.state.bushels -= res
	s.state.popFed = int(float64(res / 20))
	term.grainRemaining(s.state.bushels, res)

	failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.cows) + " cows to slaughter!"
	slaughterMsg := fmt.Sprintf("How many cows would you like to slaughter in order to feed %d people?", cowFeedMultiplier)
	res = playerInput(slaughterMsg, 0, s.state.cows, failMsg, "slaughtered")
	s.state.cows -= res
	s.state.popFed += cowFeedMultiplier * res
}

func (s *gameSession) agriculture(term *terminal) {
	cowCost := 1000
	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.bushels) + " bushels available!"
	maxCows := s.state.bushels / cowCost
	res := playerInput("How many cows will you purchase, at "+strconv.Itoa(cowCost)+" bushels per cow?", 0, maxCows, failMsg, "purchased")
	s.state.cows += res
	s.state.bushels -= res * cowCost
	term.grainRemaining(s.state.bushels, res)

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
	res = playerInput("How many fields will you plant?", maxPlantable, maxPlantable, failMsg, "planted")
	if ableToPlant > res {
		s.state.nonFarmer = s.state.population - (res-(effectivePlows*15))/10
	}
	// fmt.Printf("\tNon-farmers: %d\n", s.state.nonFarmer) // DEBUG

	s.state.bushels -= res
	s.state.planted = res
	term.grainRemaining(s.state.bushels, res)
}

func (s *gameSession) technology(term *terminal) {
	costPlow := 100

	maxPlows := s.state.bushels / costPlow
	failMsg := "Think again Hamurabi, you only have enough to purchase " + strconv.Itoa(maxPlows) + " plows!"
	res := playerInput("Do you wish to order the purchase of plows for 100 bushels, these will make it easier "+
		"for your people to plant the fields?", 0, maxPlows, failMsg, "purchased")
	s.state.plows += res
	s.state.bushels -= res * costPlow
	term.grainRemaining(s.state.bushels, res)
}

func (s *gameSession) construction(term *terminal) {
	baseCostGranary := 500
	costGranary := (s.state.granary + 1) * baseCostGranary
	palaceCost1 := 10000
	palaceCost2 := 50000
	palaceCost3 := 250000

	if yn("My lord, do you wish to consider construction projects this year") {

		// granaries
		maxGranaries := s.state.bushels / costGranary
		failMsg := "Think again Hamurabi, you only have enough to purchase " + strconv.Itoa(maxGranaries) + " granaries!"

		inputString := fmt.Sprintf("Do you wish to order the construction of city granaries for %d bushels, each "+
			"are able to protect a large amount of precious barley?", costGranary)
		res := playerInput(inputString, 0, maxGranaries, failMsg, "built")
		s.state.granary += res
		s.state.bushels -= res * costGranary
		term.grainRemaining(s.state.bushels, res)

		// palace
		var pres bool
		var buildCost int
		var typePalace int
		if s.state.buildingPalace == -1 { // if we're already building, don't ask to build more
			switch {
			case s.state.palace2 && s.state.bushels >= palaceCost3:
				typePalace = 3
				buildCost = palaceCost3
				prompt := fmt.Sprintf("Lord shall we begin expansion of your palace at a cost of %d", palaceCost3)
				pres = yn(prompt)
			case s.state.palace1 && s.state.bushels >= palaceCost2 && !s.state.palace2:
				typePalace = 2
				buildCost = palaceCost2
				prompt := fmt.Sprintf("Lord shall we begin expansion of your palace at a cost of %d", palaceCost2)
				pres = yn(prompt)
			case !s.state.palace1 && s.state.bushels >= palaceCost1:
				typePalace = 1
				buildCost = palaceCost1
				prompt := fmt.Sprintf("Lord shall we begin construction on a palace at a cost of %d", palaceCost1)
				pres = yn(prompt)
			}
			if pres {
				s.palaceBuilding = typePalace
				s.state.buildingPalace = 0
				s.state.bushels -= buildCost
				term.grainRemaining(s.state.bushels, res)
			}
		}
	}
}

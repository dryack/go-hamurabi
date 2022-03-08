package main

import (
	"fmt"
	"sort"
	"strconv"
)

func (s *gameSession) getAcres(term *terminal) {
	bushelChange := 0
	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.bushels) + " to use for purchase!"

	res := playerInput("How many acres do you wish to buy?", 0, s.state.resources.bushels/s.state.tradeVal, failMsg, "bought")
	if res == 0 {
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.acres) + " acres to sell!"
		res = playerInput("How many acres do you wish to sell?", 0, s.state.resources.acres, failMsg, "sold") * -1
		if res == 0 {
			return
		}
		bushelChange = res * -s.state.tradeVal // forcing a positive value
	} else {
		bushelChange = res * -s.state.tradeVal
	}
	s.state.resources.acres += res
	s.state.resources.bushels += bushelChange
	term.grainRemaining(s.state.resources.bushels, res)
}

func (s *gameSession) feedPeople(term *terminal) {
	cowFeedMultiplier := 45

	var ary = []int{s.state.resources.bushels, s.state.resources.population * 20}
	sort.Ints(ary)
	reqBushels := ary[0] // accept the lowest between bushels and bushels needed for the population

	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.bushels) + " available!"

	res := playerInput("How many bushels do you wish to release to your people?", reqBushels, s.state.resources.bushels, failMsg, "released")
	s.state.resources.bushels -= res
	s.state.popFed = int(float64(res / 20))
	term.grainRemaining(s.state.resources.bushels, res)

	failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.cows) + " cows to slaughter!"
	slaughterMsg := fmt.Sprintf("How many cows would you like to slaughter in order to feed %d people?", cowFeedMultiplier)
	res = playerInput(slaughterMsg, 0, s.state.resources.cows, failMsg, "slaughtered")
	s.state.resources.cows -= res
	s.state.popFed += cowFeedMultiplier * res
}

func (s *gameSession) agriculture(term *terminal) {
	cowCost := 1000
	failMsg := "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.bushels) + " bushels available!"
	maxCows := s.state.resources.bushels / cowCost
	res := playerInput("How many cows will you purchase, at "+strconv.Itoa(cowCost)+" bushels per cow?", 0, maxCows, failMsg, "purchased")
	s.state.resources.cows += res
	s.state.resources.bushels -= res * cowCost
	term.grainRemaining(s.state.resources.bushels, res)

	if s.state.technology.orchard {
		orchardCost := 40
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.acres) + " acres available!"
		maxOrchards := s.state.resources.acres / orchardCost
		res = playerInput("How many orchards shall we create, at "+strconv.Itoa(orchardCost)+"acres per orchard?", 0, maxOrchards, failMsg, "converted")
		s.state.resources.acres -= res * orchardCost
		s.state.resources.orchards += res
	}

	// only 1/3rd of the population can benefit from the plows
	var plowAry = []int{s.state.resources.plows, s.state.resources.population / 3}
	sort.Ints(plowAry)
	effectivePlows := plowAry[0]
	ableToPlant := ((s.state.resources.population - effectivePlows) * 10) + (15 * effectivePlows)

	// fmt.Printf("\tpopReqForPlanting: %d\n", ableToPlant) // DEBUG

	var ary = []int{s.state.resources.bushels, ableToPlant, s.state.resources.acres - (s.state.resources.cows * 3)}
	sort.Ints(ary)
	maxPlantable := ary[0]
	// avoid cows forcing us to negative numbers
	if maxPlantable < 0 {
		maxPlantable = 0
	}

	switch maxPlantable {
	case s.state.resources.bushels:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.bushels) + " available!"
	case s.state.popFed * ableToPlant:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.population) + " people to plant the fields!"
	case s.state.resources.acres:
		failMsg = "Think again Hamurabi, you only have " + strconv.Itoa(s.state.resources.acres) + " acres to plant!"
	}
	res = playerInput("How many fields will you plant?", maxPlantable, maxPlantable, failMsg, "planted")
	if ableToPlant > res {
		s.state.nonFarmer = s.state.resources.population - (res-(effectivePlows*15))/10
	}
	// fmt.Printf("\tNon-farmers: %d\n", s.state.nonFarmer) // DEBUG

	s.state.resources.bushels -= res
	s.state.planted = res
	term.grainRemaining(s.state.resources.bushels, res)
}

func (s *gameSession) technology(term *terminal) {
	costPlow := 100

	maxPlows := s.state.resources.bushels / costPlow
	failMsg := "Think again Hamurabi, you only have enough to purchase " + strconv.Itoa(maxPlows) + " plows!"
	res := playerInput("Do you wish to order the purchase of plows for 100 bushels, these will make it easier "+
		"for your people to plant the fields?", 0, maxPlows, failMsg, "purchased")
	s.state.resources.plows += res
	s.state.resources.bushels -= res * costPlow
	term.grainRemaining(s.state.resources.bushels, res)
}

func (s *gameSession) construction(term *terminal) {
	baseCostGranary := 500
	costGranary := (s.state.resources.granary + 1) * baseCostGranary
	palaceCost1 := 10000
	palaceCost2 := 50000
	palaceCost3 := 250000

	if yn("My lord, do you wish to consider construction projects this year") {

		// granaries
		maxGranaries := s.state.resources.bushels / costGranary
		failMsg := "Think again Hamurabi, you only have enough to purchase " + strconv.Itoa(maxGranaries) + " granaries!"

		inputString := fmt.Sprintf("Do you wish to order the construction of city granaries for %d bushels, each "+
			"are able to protect a large amount of precious barley?", costGranary)
		res := playerInput(inputString, 0, maxGranaries, failMsg, "built")
		s.state.resources.granary += res
		s.state.resources.bushels -= res * costGranary
		term.grainRemaining(s.state.resources.bushels, res)

		// palace
		var pres bool
		var buildCost int
		var typePalace int

		// if we're already building, or our palace is maxed out, don't ask to build more
		if s.state.buildingPalace == -1 && !s.state.structures.palace3 {
			switch {
			case s.state.structures.palace2 && s.state.resources.bushels >= palaceCost3:
				typePalace = 3
				buildCost = palaceCost3
				prompt := fmt.Sprintf("Lord shall we begin expansion of your palace at a cost of %d", palaceCost3)
				pres = yn(prompt)
			case s.state.structures.palace1 && s.state.resources.bushels >= palaceCost2 && !s.state.structures.palace2:
				typePalace = 2
				buildCost = palaceCost2
				prompt := fmt.Sprintf("Lord shall we begin expansion of your palace at a cost of %d", palaceCost2)
				pres = yn(prompt)
			case !s.state.structures.palace1 && s.state.resources.bushels >= palaceCost1:
				typePalace = 1
				buildCost = palaceCost1
				prompt := fmt.Sprintf("Lord shall we begin construction on a palace at a cost of %d", palaceCost1)
				pres = yn(prompt)
			}
			if pres {
				s.palaceBuilding = typePalace
				s.state.buildingPalace = 0
				s.state.resources.bushels -= buildCost
				term.grainRemaining(s.state.resources.bushels, res)
			}
		}
	}
}

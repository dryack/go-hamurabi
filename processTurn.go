package main

import (
	"fmt"
	"math/rand"
	"os"
)

func checkForPlague(state *cityState) bool {
	if state.year > 0 && rand.Intn(15) == 0 {
		fmt.Println("\nA horrible plague has struck! Many have died!")
		state.died = state.population / (rand.Intn(4) + 2)
		state.cows = state.cows / 4
		state.population -= state.died
		state.totalDead += state.died
		return true
	}
	return false
}

func printYearResults(state *cityState) {
	var otherCityStates = []string{"Dūr-Katlimmu", "Aššur", "Uruk", "Akshak", "Ur", "Nippur", "Lagash", "Larak"}

	if state.year > 0 {
		doNumbers(state)
	}
	state.tradeVal = 17 + rand.Intn(10)
	fmt.Printf("\nMy lord, in the year %d, I beg to report to you that %d people starved, %d were born, and %d "+
		"came to the city.\n", state.year, state.starved, state.born, state.migrated)
	fmt.Printf("Population is now %d.\n", state.population)
	fmt.Printf("The city owns %d acres of land, and has %d granaries.\n", state.acres, state.granary)

	// we can't support the cows - so they are killed
	if state.acres < state.cows*3 {
		slaughtered := 0
		if state.acres <= 2 {
			slaughtered = state.cows
		} else {
			slaughtered = (state.acres / 3) % state.cows
		}
		state.cows -= slaughtered
		fmt.Printf("As we lacked the land to support them, %d cows were slaughtered!\n", slaughtered)
	}

	var cowsFed int
	if state.cows*state.cowMultiplier > state.population {
		cowsFed = state.population
	} else {
		cowsFed = state.cows * state.cowMultiplier
	}
	fmt.Printf("The city keeps %d cows whose product fed %d people this year.\n", state.cows, cowsFed)
	if state.acres < 1 || state.planted == 0 {
		fmt.Printf("Traders report that %s harvested %d bushels per acrs.\n", otherCityStates[rand.Intn(len(otherCityStates)-1)], state.bYield)
	}
	fmt.Printf("We have harvested %d bushels per acre.\n", state.bYield)
	if state.nonFarmer > 0 && state.tradeGoods > 0 {
		fmt.Printf("Thanks to having %d citizens not required to farm, trade goods and vegatables brought in %d "+
			"bushels of grain.\n", state.nonFarmer, state.tradeGoods)
	}
	fmt.Printf("Rats ate %d bushels of grain.\n", state.pests)
	fmt.Printf("We now have %d bushels in store.\n", state.bushels)
	fmt.Printf("We have distributed a total of %d hand plows amongst the people.\n", state.plows)
	fmt.Printf("Land is trading at %d bushels per acre.\n", state.tradeVal)
	state.year += 1
}

func doNumbers(state *cityState) {
	plague := checkForPlague(state)

	state.bYield = rand.Intn(9) + 1

	// starvation & population
	state.starved = state.population - (state.popFed + state.cows*state.cowMultiplier)
	if state.starved < 0 {
		state.starved = 0
	}
	state.avgStarved = int(float64(state.starved) / float64(state.population) * 100)
	state.born = int(float64(state.population) / float64(rand.Intn(8)+2))
	if plague {
		state.born /= 2 // children die from the plague as well
	}
	checkForOverthrow(state)

	state.population += state.born
	state.avgStarved = int(float64(state.starved) / float64(state.population) * 100)
	state.population -= state.starved // children die too

	// migration
	var cowMigrantAttraction int
	if state.cows > 3 {
		cowMigrantAttraction = state.cows * 5
	} else {
		cowMigrantAttraction = 0
	}
	if plague {
		// people don't come to a place with a plague
		state.migrated = (int(0.1*float64(rand.Intn(state.population)+1)) + cowMigrantAttraction) / 5
	} else {
		state.migrated = int(0.1*float64(rand.Intn(state.population)+1)) + cowMigrantAttraction
	}
	state.population += state.migrated

	// pests
	granaryProtectMultiplier := 3000
	unprotectedGrain := state.bushels - state.granary*granaryProtectMultiplier
	if unprotectedGrain < 0 {
		unprotectedGrain = 0
	}
	state.pests = int(float64(unprotectedGrain) / float64(rand.Intn(4)+3))

	// agricultural results
	state.bushels += (state.planted - state.cows*3) * state.bYield
	state.bushels -= state.pests
	if state.bushels < 0 {
		state.bushels = 0
	}

	state.tradeGoods = state.nonFarmer * (rand.Intn(49) + 1)
	state.bushels += state.tradeGoods

	state.totalDead += state.starved
	state.avgPestEaten += state.pests
	state.avgBushelsAvail += state.bushels
}

func checkForOverthrow(state *cityState) {
	if state.starved > int(0.45*float64(state.population)) {
		fmt.Printf("\nYou starved %d out of your population of only %d, this has caused you to be deposed by force!\n",
			state.starved, state.population)
		state.totalDead += state.starved
		endOfReign(state)
	}
}

func endOfReign(state *cityState) {
	fmt.Printf("In your %d year reign %d percent of the population starved per year on average. A total of %d "+
		"people died during your reign.\n", state.year, state.avgStarved/state.year, state.totalDead)
	fmt.Printf("The city began with 100 citizens and ended with %d.\n", state.population)
	fmt.Printf("You ordered the building of %d granaries during your rule.\n", state.granary)
	fmt.Printf("The city began with 10 acres per person and ended with %d.\n", state.acres/state.population)
	fmt.Printf("\tAvg Bushels at turn start: %d; Avg Bushels eaten by rats: %d\n", state.avgBushelsAvail/state.turns, state.avgPestEaten/state.turns)

	os.Exit(0)
}

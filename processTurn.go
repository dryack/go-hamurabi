package main

import (
	"fmt"
	"math/rand"
	"os"
)

func checkForPlague(state *cityState) {
	if state.year > 0 && rand.Intn(15) == 0 {
		fmt.Println("A horrible plague has struck! Many have died!")
		state.died = state.population / (rand.Intn(4) + 2)
		state.cows = state.cows / 4
		state.population -= state.died
		state.totalDead += state.died
	}
}

func printYearResults(state *cityState) {
	if state.year > 0 {
		doNumbers(state)
	}
	state.tradeVal = 17 + rand.Intn(10)
	fmt.Printf("\nMy lord, in the year %d, I beg to report to you that %d people starved, %d were born, and %d "+
		"came to the city.\n", state.year, state.starved, state.born, state.migrated)
	checkForPlague(state)
	fmt.Printf("Population is now %d.\n", state.population)
	fmt.Printf("The city owns %d acres of land.\n", state.acres)
	var cowsFed int
	if state.cows*15 > state.population {
		cowsFed = state.population
	} else {
		cowsFed = state.cows * 15
	}
	fmt.Printf("The city keeps %d cows whose produce fed %d people this year.\n", state.cows, cowsFed)
	fmt.Printf("We have harvested %d bushels per acre.\n", state.bYield)
	fmt.Printf("Rats ate %d bushels of grain.\n", state.pests)
	fmt.Printf("We now have %d bushels in store.\n", state.bushels)
	fmt.Printf("Land is trading at %d bushels per acre.\n", state.tradeVal)
	state.year += 1
}

func doNumbers(state *cityState) {
	state.bYield = rand.Intn(9) + 1
	state.starved = state.population - (state.popFed + state.cows*20)
	if state.starved < 0 {
		state.starved = 0
	}
	state.avgStarved = int(float64(state.starved) / float64(state.population) * 100)
	state.born = int(float64(state.population) / float64(rand.Intn(8)+2))
	state.population += state.born
	checkForOverthrow(state)
	state.avgStarved = int(float64(state.starved) / float64(state.population) * 100)
	state.population -= state.starved // children die too
	state.migrated = int(0.1 * float64(rand.Intn(state.population)+1))
	state.population += state.migrated
	unprotectedGrain := state.bushels - state.granary*3000
	if unprotectedGrain < 0 {
		unprotectedGrain = 0
	}
	state.pests = int(float64(unprotectedGrain) / float64(rand.Intn(4)+3))
	state.bushels += (state.planted - state.cows*3) * state.bYield
	state.bushels -= state.pests
	if state.bushels < 0 {
		state.bushels = 0
	}
	state.totalDead += state.starved
	state.avgPestEaten += state.pests
	state.avgBushelsAvail += state.bushels
}

func checkForOverthrow(state *cityState) {
	if state.starved > int(0.45*float64(state.population)) {
		fmt.Printf("\nYou starved %d out of your population of only %d, this has cause you to be deposed by force!\n",
			state.starved, state.population)
		state.totalDead += state.starved
		endOfReign(state)
	}
}

func endOfReign(state *cityState) {
	fmt.Printf("In your %d year reign %d percent of the population starved per year on average. A total of %d "+
		"people died during your reign.\n", state.year, state.avgStarved/state.year, state.totalDead)
	fmt.Printf("The city began with 10 acres per person and ended with %d.\n", state.acres/state.population)
	fmt.Printf("\tAvg Bushels at turn start: %d; Avg Bushels eaten by rats: %d\n", state.avgBushelsAvail/state.turns, state.avgPestEaten/state.turns)

	os.Exit(0)
}

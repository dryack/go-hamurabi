package main

import (
	"fmt"
	"os"
)

func (s *gameSession) endOfReign() {
	fmt.Printf("In your %d year reign %d percent of the population starved per year on average. A total of %d "+
		"people died during your reign.\n", s.state.year, s.avgStarved/s.state.year, s.totalDead)
	fmt.Printf("The city began with 100 citizens and ended with %d.\n", s.state.population)
	fmt.Printf("You ordered the building of %d granaries during your rule.\n", s.state.granary)
	fmt.Printf("The city began with 10 acres per person and ended with %d.\n", s.state.acres/s.state.population)
	fmt.Printf("We maintained a herd of %d cows.\n", s.state.cows)
	fmt.Printf("\tAvg Bushels at turn start: %d; Avg Bushels eaten by rats: %d\n", s.avgBushelsAvail/s.turns, s.avgPestEaten/s.turns) // DEBUG

	os.Exit(0)
}

package main

import (
	"fmt"
	"os"
)

func (s *gameSession) endOfReign() {
	fmt.Printf("\nIn your %d year reign %d percent of the population starved per year on average. A total of %d "+
		"people died during your reign.\n", s.state.year, s.avgStarved/s.state.year, s.totalDead)
	fmt.Printf("The city began with 100 citizens and ended with %d.\n", s.state.resources.population)
	fmt.Printf("You ordered the building of %d granaries during your rule.\n", s.state.resources.granary)
	// FIXME: can result in a divide by zero error in the case of drastic loss of population due to plague
	fmt.Printf("The city began with 10 acres per person and ended with %d.\n", s.state.resources.acres/s.state.resources.population)
	fmt.Printf("We maintained a herd of %d cows.\n", s.state.resources.cows)
	fmt.Printf("Your final score was %d\n", s.calcScore())
	fmt.Printf("\tAvg Bushels at turn start: %d; Avg Bushels eaten by rats: %d\n", s.avgBushelsAvail/s.turns, s.avgPestEaten/s.turns) // DEBUG

	os.Exit(0)
}

func (s *gameSession) calcScore() int {
	var (
		population = 1
		cows       = 5
		land       = 10
		plows      = 2
		palace1    = 200
		palace2    = 1000
		palace3    = 5000
		granary    = 15
	)
	score := (population * s.state.resources.population) +
		(cows * s.state.resources.cows) +
		(land * s.state.resources.acres) +
		(plows * s.state.resources.plows) +
		(granary * s.state.resources.granary) +
		(palace1 * func(p bool) int {
			if p {
				return 1
			}
			return 0
		}(s.state.structures.palace1)) +
		(palace2 * func(p bool) int {
			if p {
				return 1
			}
			return 0
		}(s.state.structures.palace2)) +
		(palace3 * func(p bool) int {
			if p {
				return 1
			}
			return 0
		}(s.state.structures.palace3))
	return score
}

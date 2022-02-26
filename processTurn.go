package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

func (s *gameSession) checkForPlague() bool {
	if s.state.year > 0 && rand.Intn(15) == 0 {
		fmt.Println("\nA horrible plague has struck! Many have died!")
		s.state.died = s.state.population / (rand.Intn(4) + 2)
		s.state.cows = s.state.cows / 4
		s.state.population -= s.state.died
		s.totalDead += s.state.died
		return true
	}
	return false
}

func (s *gameSession) printYearResults() {
	if s.state.year > 0 {
		s.doNumbers()
	}
	fmt.Printf("\nMy lord, in the year %d, I beg to report to you that %d people starved, %d were born, and %d "+
		"came to the city.\n", s.state.year, s.state.starved, s.state.born, s.state.migrated)
	fmt.Printf("Population is now %d.\n", s.state.population)
	fmt.Printf("The city owns %d acres of land, and has %d granaries.\n", s.state.acres, s.state.granary)

	// we can't support the cows - so they are killed
	if s.state.forceSlaughtered > 0 {
		fmt.Printf("As we lacked the land to support them, %d cows were slaughtered!\n", s.state.forceSlaughtered)
	}
	fmt.Printf("The city keeps %d cows whose product fed %d people this year.\n", s.state.cows, s.state.cowsFed)

	if s.state.acres < 1 || s.state.planted == 0 {
		fmt.Printf("Traders report that %s harvested %d bushels per acrs.\n", s.otherCityStates[rand.Intn(len(s.otherCityStates)-1)], s.state.bYield)
	} else {
		fmt.Printf("We have harvested %d bushels per acre.\n", s.state.bYield)
	}

	if s.state.nonFarmer > 0 && s.state.tradeGoods > 0 {
		fmt.Printf("Thanks to having %d citizens not required to farm, trade goods and vegatables brought in %d "+
			"bushels of grain.\n", s.state.nonFarmer, s.state.tradeGoods)
	}

	fmt.Printf("Rats ate %d bushels of grain.\n", s.state.pests)
	fmt.Printf("We now have %d bushels in store.\n", s.state.bushels)
	fmt.Printf("We have distributed a total of %d hand plows amongst the people.\n", s.state.plows)
	fmt.Printf("Land is trading at %d bushels per acre.\n", s.state.tradeVal)
	s.state.year += 1
}

func (s *gameSession) doNumbers() {
	plague := s.checkForPlague()

	s.state.tradeVal = 17 + rand.Intn(10)
	s.state.bYield = rand.Intn(9) + 1
	// cows
	s.doCows()
	// starvation & population
	s.doStarvation(plague)
	s.checkForOverthrow()

	s.state.population += s.state.born
	s.avgStarved = int(float64(s.state.starved) / float64(s.state.population) * 100)
	s.state.population -= s.state.starved // children die too
	// migration
	s.doMigration(plague)
	// pests
	s.doPests()
	// agricultural results
	s.doAgriculture()

	s.state.tradeGoods = s.state.nonFarmer * (rand.Intn(49) + 1)
	s.state.bushels += s.state.tradeGoods
	s.totalDead += s.state.starved
	s.avgPestEaten += s.state.pests
	s.avgBushelsAvail += s.state.bushels
}

func (s *gameSession) doAgriculture() {
	s.state.bushels += (s.state.planted - s.state.cows*3) * s.state.bYield
	s.state.bushels -= s.state.pests
	if s.state.bushels < 0 {
		s.state.bushels = 0
	}

	// although the peasants don't have to sow, land must be tended or it will become wasted and be reclaimed by nature
	// some lands are tended by the royal staff, and although they can be sold, they CAN'T go to waste
	royalLands := 500
	fieldMaintPerPop := 30
	maxAcresMaint := s.state.population * fieldMaintPerPop
	// we don't lose the royal-held lands to wastage from lack of peasants
	if s.state.acres > royalLands {
		// if there aren't enough peasants to maintain our acreage
		if maxAcresMaint < s.state.acres {
			s.state.acresWastage = int(math.Abs(float64(maxAcresMaint - (s.state.acres - royalLands))))
			fmt.Printf("Due to a lack of peasants to work the land, %d acres have wasted and are lost!\n", s.state.acresWastage)
		} else {
			s.state.acresWastage = 0
		}
	} else {
		s.state.acresWastage = 0
	}
	s.state.acres -= s.state.acresWastage
	if s.state.acres < royalLands {
		s.state.acres = royalLands
		fmt.Println("However your personal retainers protected your personal estate!")
	}
}

func (s *gameSession) doPests() {
	granaryProtectMultiplier := 3000
	unprotectedGrain := s.state.bushels - s.state.granary*granaryProtectMultiplier
	if unprotectedGrain < 0 {
		unprotectedGrain = 0
	}
	s.state.pests = int(float64(unprotectedGrain) / float64(rand.Intn(4)+3))
}

func (s *gameSession) doMigration(plague bool) {
	var cowMigrantAttraction int
	if s.state.cows > 3 {
		cowMigrantAttraction = s.state.cows * 5
	} else {
		cowMigrantAttraction = 0
	}
	if plague {
		// people don't come to a place with a plague
		s.state.migrated = (int(0.1*float64(rand.Intn(s.state.population)+1)) + cowMigrantAttraction) / 5
	} else {
		s.state.migrated = int(0.1*float64(rand.Intn(s.state.population)+1)) + cowMigrantAttraction
	}
	s.state.population += s.state.migrated
}

func (s *gameSession) doStarvation(plague bool) {
	s.state.starved = s.state.population - (s.state.popFed + s.state.cows*s.state.cowMultiplier)
	if s.state.starved < 0 {
		s.state.starved = 0
	}
	s.avgStarved = int(float64(s.state.starved) / float64(s.state.population) * 100)
	s.state.born = int(float64(s.state.population) / float64(rand.Intn(8)+2))
	if plague {
		s.state.born /= 2 // children die from the plague as well
	}
}

func (s *gameSession) doCows() {
	if s.state.acres < s.state.cows*3 {
		s.state.forceSlaughtered = 0
		if s.state.acres <= 2 {
			s.state.forceSlaughtered = s.state.cows
		} else {
			s.state.forceSlaughtered = (s.state.acres / 3) % s.state.cows
		}
		s.state.cows -= s.state.forceSlaughtered
	}
	if s.state.cows*s.state.cowMultiplier > s.state.population {
		s.state.cowsFed = s.state.population
	} else {
		s.state.cowsFed = s.state.cows * s.state.cowMultiplier
	}
}

func (s *gameSession) checkForOverthrow() {
	if s.state.starved > int(0.45*float64(s.state.population)) {
		fmt.Printf("\nYou starved %d out of your population of only %d, this has caused you to be deposed by force!\n",
			s.state.starved, s.state.population)
		s.totalDead += s.state.starved
		s.endOfReign()
	}

	if s.state.population < 10 {
		fmt.Printf("\nYour continued mismanagement caused your population to decline to the point that the " +
			"remaining peasants fled your land\nYou are left ruling an empty city, as your royal guards and staff escape.\n")
		s.state.population = 0
		s.endOfReign()
	}
}

func (s *gameSession) endOfReign() {
	fmt.Printf("In your %d year reign %d percent of the population starved per year on average. A total of %d "+
		"people died during your reign.\n", s.state.year, s.avgStarved/s.state.year, s.totalDead)
	fmt.Printf("The city began with 100 citizens and ended with %d.\n", s.state.population)
	fmt.Printf("You ordered the building of %d granaries during your rule.\n", s.state.granary)
	fmt.Printf("The city began with 10 acres per person and ended with %d.\n", s.state.acres/s.state.population)
	fmt.Printf("\tAvg Bushels at turn start: %d; Avg Bushels eaten by rats: %d\n", s.avgBushelsAvail/s.turns, s.avgPestEaten/s.turns) // DEBUG

	os.Exit(0)
}

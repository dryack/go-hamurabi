package main

import (
	"fmt"
	"github.com/muesli/termenv"
	"math"
	"math/rand"
)

func (s *gameSession) checkForPlague() bool {
	chance := 7 + (s.state.resources.population / 5000)
	// a plague has a 50% chance of continuing each year
	if s.state.plagueLY {
		chance = 50 + (s.state.resources.population / 3000)
	}
	if s.state.year > 0 && rand.Intn(99)+1 < chance {
		return true
	}
	return false

	// if s.state.year > 0 && rand.Intn(15) == 0 {
	// 	// TODO: let's do this under printYearResults()
	//
	// 	// TODO: let's do this under doNumbers()
	// 	s.state.died = s.state.resources.population / (rand.Intn(4) + 2)
	// 	s.state.resources.population -= s.state.died
	// 	s.totalDead += s.state.died
	// 	return true

}

func (s *gameSession) printYearResults(term *terminal) {
	const (
		heading = "\nMy lord, in the year %s, I beg to report to you that %s people starved, our native population %s %s, and %s " +
			"came to the city.\n"
		newPopulation        = "Population is now %s.\n"
		acresAndGranaries    = "The city owns %s acres of land, and has %s granaries.\n"
		palaceCompletion     = "My Lord, your workers have completed work on your palace!"
		palaceMsg1           = "You are residing in a large palace, together with your family and closest retainers."
		palaceMsg2           = "You are residing in a huge palace, together with your family and many retainers."
		palaceMsg3           = "You are residing in a massive bustling palace, together with you family, many retainers, royal merchants, and visiting diplomats."
		noPalaceWork         = "My Dread Lord, I regret that due to the plague, no work was done on your palace!"
		buildingPalace       = "Construction on your palace is underway, and will be completed in %s years.\n"
		expandingPalace      = "Expansion of your palace is underway, and will be completed in %s years.\n"
		forceSlaughteredCows = "%s cows were slaughtered, as we lacked the land to support them!\n"
		cowsMsg              = "The city keeps %s cows whose product fed %s people this year.\n"
		traderReportsYield   = "Traders report that %s harvested %s bushels per acre.\n"
		harvestYield         = "We have harvested %s bushels per acre.\n"
		citizenTraders       = "Thanks to having %s citizens not required to farm, trade goods and vegetables brought in %s " +
			"bushels of grain.\n"
		rats          = "Rats ate %s bushels of grain.\n"
		storedBushels = "We now have %s bushels in store.\n"
		plows         = "We have distributed a total of %s hand plows amongst the people.\n"
		landValue     = "Land is trading at %s bushels per acre.\n"
	)

	var plague bool
	var palaceComplete bool
	if s.state.year > 0 {
		plague, palaceComplete = s.doNumbers(term)
	}

	fmt.Printf(heading, term.pink(s.state.year), term.pink(s.state.starved), s.lostOrGained(), term.pink(Abs(s.state.born)), term.pink(s.state.migrated))
	fmt.Printf(newPopulation, term.pink(s.state.resources.population))
	fmt.Printf(acresAndGranaries, term.pink(s.state.resources.acres), term.pink(s.state.resources.granary))
	if palaceComplete {
		// TODO: colorCode() or new func to support full string coloration
		fmt.Println(termenv.String(palaceCompletion).Bold().Foreground(term.p.Color("226")))
	}

	switch {
	case s.state.structures.palace3:
		// TODO: colorCode() or new func to support full string coloration
		fmt.Println(termenv.String(palaceMsg3).Bold().Background(term.p.Color("214")).Foreground(term.p.Color("16")))
	case s.state.structures.palace2:
		fmt.Println(termenv.String(palaceMsg2).Bold().Background(term.p.Color("220")).Foreground(term.p.Color("16")))
	case s.state.structures.palace1:
		// TODO: color 226 needs to be double checked
		fmt.Println(termenv.String(palaceMsg1).Bold().Background(term.p.Color("226")).Foreground(term.p.Color("16")))
	}

	// if we're supposed to be building a palace, how many more years are left to build it?
	if s.state.buildingPalace > -1 {
		switch {
		case plague: // no work is done during a plague year
			// TODO: colorCode() or new func to support full string coloration
			fmt.Println(termenv.String(noPalaceWork).Bold().Foreground(term.p.Color("226")))
			fallthrough
		case !s.state.structures.palace1:
			fmt.Printf(buildingPalace, term.pink(5-s.state.buildingPalace))
		case s.state.structures.palace1:
			fmt.Printf(expandingPalace, term.pink(5-s.state.buildingPalace))
		case s.state.structures.palace2:
			fmt.Printf(expandingPalace, term.pink(5-s.state.buildingPalace))
		}
	}

	// we can't support the cows - so they are killed
	if s.state.forceSlaughtered > 0 {
		fmt.Printf(forceSlaughteredCows, term.red(s.state.forceSlaughtered))
	}
	if s.state.resources.cows > 0 {
		fmt.Printf(cowsMsg, term.pink(s.state.resources.cows), term.pink(s.state.cowsFed))
	}

	if s.state.resources.acres < 1 || s.state.planted == 0 {
		neighboringCity := s.otherCityStates[rand.Intn(len(s.otherCityStates)-1)]
		fmt.Printf(traderReportsYield, term.pink(s.state.bYield), neighboringCity)
	} else {
		fmt.Printf(harvestYield, term.pink(s.state.bYield))
	}

	if s.state.nonFarmer > 0 && s.state.tradeGoods > 0 {
		fmt.Printf(citizenTraders, term.pink(s.state.nonFarmer), term.pink(s.state.tradeGoods))
	}

	fmt.Printf(rats, term.pink(s.state.pests))
	fmt.Printf(storedBushels, term.pink(s.state.resources.bushels))
	fmt.Printf(plows, term.pink(s.state.resources.plows))
	fmt.Printf(landValue, term.pink(s.state.tradeVal))
	// TODO: state change shouldn't happen in messages
	s.state.year += 1
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s *gameSession) doNumbers(term *terminal) (bool, bool) {
	s.state.forceSlaughtered = 0 // reset this each turn
	s.state.tradeVal = 17 + rand.Intn(10)

	plague := s.checkForPlague()
	if !plague {
		s.state.plagueLY = false // there's no ongoing plague
	}

	s.doPopulation(plague, term) // starvation & population
	s.doCows(plague)
	s.checkForOverthrow(term)
	s.doMigration(plague)
	s.doPests()
	s.doAgriculture(term)
	palaceComplete := s.doConstruction(plague)

	// trade is reduced during plague
	if plague {
		s.state.tradeGoods = s.state.nonFarmer * (rand.Intn(15) + 1)
	} else {
		s.state.tradeGoods = s.state.nonFarmer * (rand.Intn(49) + 1)
	}
	// 10% of trade is taxed, and converted to silver shekels.  for future consideration, 3600 mina equal a talent
	// and 1 mina equals 60 shekels
	if s.state.technology.silver {
		taxes := s.state.tradeGoods / 10
		s.state.tradeGoods -= taxes
		s.state.resources.minaOfSilver += taxes
	}

	s.state.resources.bushels += s.state.tradeGoods
	s.totalDead += s.state.starved
	s.avgPestEaten += s.state.pests
	s.avgBushelsAvail += s.state.resources.bushels
	return plague, palaceComplete
}

func (s *gameSession) doConstruction(plague bool) bool {
	palaceComplete := false

	// if there's no plague, advance construction of the palace by a year
	if !plague && s.state.buildingPalace > -1 {
		s.state.buildingPalace++
	}
	// if construction is complete, set the correct palace type and its technologies
	if s.state.buildingPalace > 4 {
		switch s.palaceBuilding {
		case 1:
			s.state.structures.palace1 = true
			s.state.technology.silver = true
			s.state.technology.orchard = true
		case 2:
			s.state.structures.palace2 = true
			s.state.technology.stela = true
		case 3:
			s.state.structures.palace3 = true
			s.state.technology.ziggurat = true
		}
		palaceComplete = true // returned to processTurn()
		s.state.buildingPalace = -1
		s.palaceBuilding = -1 // reset palace construction once it's been completed
	}
	return palaceComplete
}

func (s *gameSession) doAgriculture(term *terminal) {
	s.state.bYield = rand.Intn(9) + 1
	s.state.resources.bushels += (s.state.planted - s.state.resources.cows*3) * s.state.bYield
	s.state.resources.bushels -= s.state.pests
	if s.state.resources.bushels < 0 {
		s.state.resources.bushels = 0
	}

	// although the peasants don't have to sow, land must be tended or it will become wasted and be reclaimed by nature
	// some lands are tended by the royal staff, and although they can be sold, they CAN'T go to waste
	royalLands := 500
	fieldMaintPerPop := 30
	maxAcresMaint := s.state.populationLY * fieldMaintPerPop
	// we don't lose the royal-held lands to wastage from lack of peasants
	if s.state.resources.acres > royalLands {
		// if there aren't enough peasants to maintain our acreage
		if maxAcresMaint < s.state.resources.acres {
			s.state.acresWastage = int(math.Abs(float64(maxAcresMaint - (s.state.resources.acres - royalLands))))
			// TODO: messages shouldn't be happening during state changes
			fmt.Printf("\nDue to a lack of peasants to work the land, %s acres have wasted and are lost!\n", term.colorCode("196", s.state.acresWastage))
			s.state.resources.acres -= s.state.acresWastage
			s.totAcresWasted += s.state.acresWastage
			if s.state.resources.acres < royalLands {
				s.state.resources.acres = royalLands
				// TODO: messages shouldn't be happening during state changes
				fmt.Println("However your personal retainers protected your personal estate!")
			}
		} else {
			s.state.acresWastage = 0
		}
	} else {
		s.state.acresWastage = 0
	}
}

func (s *gameSession) doPests() {
	granaryProtectMultiplier := 3000
	unprotectedGrain := s.state.resources.bushels - s.state.resources.granary*granaryProtectMultiplier
	if unprotectedGrain < 0 {
		unprotectedGrain = 0
	}
	s.state.pests = int(float64(unprotectedGrain) / float64(rand.Intn(4)+3))
}

func (s *gameSession) doMigration(plague bool) {
	var cowMigrantAttraction int
	switch {
	case s.state.resources.cows <= 3:
		cowMigrantAttraction = 0
	case s.state.resources.cows > 3 && s.state.populationLY <= 500:
		cowMigrantAttraction = s.state.resources.cows * 5
	case s.state.resources.cows > 3 && s.state.populationLY <= 10000:
		cowMigrantAttraction = s.state.resources.cows * 3
	case s.state.resources.cows > 3 && s.state.populationLY > 10000:
	default:
		cowMigrantAttraction = 0
	}
	if plague {
		// people don't come to a place with a plague
		s.state.migrated = (int(0.1*float64(rand.Intn(s.state.populationLY)+1)) + cowMigrantAttraction) / 10
	} else {
		s.state.migrated = int(0.1*float64(rand.Intn(s.state.populationLY)+1)) + cowMigrantAttraction
	}
	s.state.resources.population += s.state.migrated
}

func (s *gameSession) doPopulation(plague bool, term *terminal) {
	dead := 0
	childDeaths := 1
	if plague {
		dead, childDeaths = s.doPlague(term)
	}
	femalePopulation := s.state.resources.population / 2
	s.state.populationLY = s.state.resources.population

	s.state.starved = s.state.populationLY - (s.state.popFed + s.state.resources.cows*s.state.technology.cowMultiplier)
	if s.state.starved < 0 {
		s.state.starved = 0
	}
	s.state.resources.population -= s.state.starved // children die too
	s.avgStarved = int(float64(s.state.starved) / float64(s.state.resources.population) * 100)
	// half the population (the females), minus half the dead if there's a plague, have a birth rate between 9% and 95%
	// which is divided by a randomly determined number (between 2 and 3 if there's a plague).  from this result we
	// subtract half of the deaths if there's a plague
	s.state.born = (int(float64(femalePopulation-(dead/2))*math.Max(rand.Float64()-0.05, 0.15)) / childDeaths) - dead/2
	s.state.resources.population += s.state.born
}

func (s *gameSession) doCows(plague bool) {
	// TODO: implement cowsLY mirroring populationLY, to prevent time travel
	// if the plague is killing cows, kill em before anything else is considered
	if plague {
		s.state.resources.cows = s.state.resources.cows/rand.Intn(5) + 1
	}

	if s.state.resources.acres < s.state.resources.cows*3 {
		s.state.forceSlaughtered = 0
		if s.state.resources.acres <= 2 {
			s.state.forceSlaughtered = s.state.resources.cows
		} else {
			s.state.forceSlaughtered = (s.state.resources.acres / 3) % s.state.resources.cows
		}
		s.state.resources.cows -= s.state.forceSlaughtered
	}
	if s.state.resources.cows*s.state.technology.cowMultiplier > s.state.populationLY {
		s.state.cowsFed = s.state.populationLY
	} else {
		s.state.cowsFed = s.state.resources.cows * s.state.technology.cowMultiplier
	}
}

func (s *gameSession) checkForOverthrow(term *terminal) {
	const (
		deposedMsg         = "\nYou starved %s out of your population of only %d, this has caused you to be deposed by force!\n"
		populationDeclined = "\nYour continued mismanagement caused your population to decline to the point that the " +
			"remaining peasants fled your land\nYou are left ruling an empty city, as your royal guards and staff escape.\n"
	)
	var stability float64 = 0.45
	switch {
	case s.state.structures.palace1:
		stability += .01
		fallthrough
	case s.state.structures.palace2:
		stability += .02
		fallthrough
	case s.state.structures.palace3:
		stability += .03
	}

	if s.state.starved > int(stability*float64(s.state.resources.population)) {
		// TODO: colorCode() or new func to support full string coloration
		fmt.Printf(deposedMsg, term.red(s.state.starved), s.state.populationLY)
		s.totalDead += s.state.starved
		s.endOfReign()
	}

	if s.state.resources.population < 10 {
		// TODO: colorCode() or new func to support full string coloration
		fmt.Print(termenv.String(populationDeclined).Bold().Background(term.p.Color("196")))
		s.state.resources.population = 1 // prevents divide by zero
		s.endOfReign()
	}
}

func (s *gameSession) doPlague(term *terminal) (int, int) {
	plagueMsg := "A horrible plague has struck! Many have died!\n"
	plagueContinuesMsg := "The terrible plague continues! Our people die in droves!"
	dead := 0
	childDeaths := 1

	// TODO: messages moved to printYearResults()
	if s.state.plagueLY {
		fmt.Print(termenv.String(plagueContinuesMsg).Bold().Background(term.p.Color("196")))
	} else {
		fmt.Print(termenv.String(plagueMsg).Bold().Background(term.p.Color("196")))
	}
	dead = s.state.populationLY / (rand.Intn(1) + 2)
	s.totalDead += dead
	childDeaths = rand.Intn(1) + 2
	return dead, childDeaths
}

func (s *gameSession) lostOrGained() string {
	switch {
	case s.state.born < 0:
		return "decreased by"
	case s.state.born >= 0:
		return "increased by"
	}
	return ""
}

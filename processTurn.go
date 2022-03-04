package main

import (
	"fmt"
	"github.com/muesli/termenv"
	"math"
	"math/rand"
)

func (s *gameSession) checkForPlague() bool {
	const plagueMsg = "\nA horrible plague has struck! Many have died!"

	if s.state.year > 0 && rand.Intn(15) == 0 {
		// TODO: let's do this under printYearResults()
		fmt.Println(plagueMsg)
		// TODO: let's do this under doNumbers()
		s.state.died = s.state.population / (rand.Intn(4) + 2)
		s.state.cows = s.state.cows / 4
		s.state.population -= s.state.died
		s.totalDead += s.state.died
		return true
	}
	return false
}

func (s *gameSession) printYearResults() {
	const (
		heading = "\nMy lord, in the year %d, I beg to report to you that %d people starved, %d were born, and %d " +
			"came to the city.\n"
		newPopulation        = "Population is now %d.\n"
		acresAndGranaries    = "The city owns %d acres of land, and has %d granaries.\n"
		palaceCompletion     = "My Lord, your workers have completed work on your palace!"
		palaceMsg1           = "You are residing in a large palace, together with your family and closest retainers."
		palaceMsg2           = "You are residing in a huge palace, together with your family and many retainers."
		palaceMsg3           = "You are residing in a massive bustling palace, together with you family, many retainers, royal merchants, and visiting diplomats."
		noPalaceWork         = "My Dread Lord, I regret that due to the plague, no work was done on your palace!"
		buildingPalace       = "Construction on your palace is underway, and will be completed in %d years.\n"
		expandingPalace      = "Expansion of your palace is underway, and will be completed in %d years.\n"
		forceSlaughteredCows = "%d cows were slaughtered, as we lacked the land to support them!\n"
		cowsMsg              = "The city keeps %d cows whose product fed %d people this year.\n"
		traderReportsYield   = "Traders report that %s harvested %d bushels per acre.\n"
		harvestYield         = "We have harvested %d bushels per acre.\n"
		citizenTraders       = "Thanks to having %d citizens not required to farm, trade goods and vegetables brought in %d " +
			"bushels of grain.\n"
		rats          = "Rats ate %d bushels of grain.\n"
		storedBushels = "We now have %d bushels in store.\n"
		plows         = "We have distributed a total of %d hand plows amongst the people.\n"
		landValue     = "Land is trading at %d bushels per acre.\n"
	)

	var plague bool
	var palaceComplete bool
	if s.state.year > 0 {
		plague, palaceComplete = s.doNumbers()
	}
	msg := s.fOut(heading, "199", s.state.year, s.state.starved, s.state.born, s.state.migrated)
	fmt.Print(msg)

	msg = s.fOut(newPopulation, "199", s.state.population)
	fmt.Print(msg)
	msg = s.fOut(acresAndGranaries, "199", s.state.acres, s.state.granary)
	fmt.Print(msg)
	if palaceComplete {
		fmt.Println(termenv.String(palaceCompletion).Bold().Foreground(s.p.Color("226")))
	}

	switch {
	case s.state.palace3:
		fmt.Println(termenv.String(palaceMsg3).Bold().Background(s.p.Color("214")).Foreground(s.p.Color("16")))
	case s.state.palace2:
		fmt.Println(termenv.String(palaceMsg2).Bold().Background(s.p.Color("220")).Foreground(s.p.Color("16")))
	case s.state.palace1:
		fmt.Println(termenv.String(palaceMsg1).Bold().Background(s.p.Color("226")).Foreground(s.p.Color("16")))
	}

	if s.state.buildingPalace > -1 {
		switch {
		case plague:
			fmt.Println(termenv.String(noPalaceWork).Bold().Foreground(s.p.Color("226")))
			fallthrough
		case !s.state.palace1:
			fmt.Printf(buildingPalace, 5-s.state.buildingPalace)
		case s.state.palace1:
			fmt.Printf(expandingPalace, 5-s.state.buildingPalace)
		case s.state.palace2:
			fmt.Printf(expandingPalace, 5-s.state.buildingPalace)
		}
	}

	// we can't support the cows - so they are killed
	if s.state.forceSlaughtered > 0 {
		msg = s.fOut(forceSlaughteredCows, "196", s.state.forceSlaughtered)
		fmt.Print(msg)
	}
	if s.state.cows > 0 {
		msg = s.fOut(cowsMsg, "199", s.state.cows, s.state.cowsFed)
		fmt.Print(msg)
	}

	if s.state.acres < 1 || s.state.planted == 0 {
		msg = fmt.Sprintf(s.fOut(traderReportsYield, "199", s.state.bYield), s.otherCityStates[rand.Intn(len(s.otherCityStates)-1)])
		fmt.Print(msg)
	} else {
		msg = fmt.Sprintf(s.fOut(harvestYield, "199", s.state.bYield))
		fmt.Printf(msg)
	}

	if s.state.nonFarmer > 0 && s.state.tradeGoods > 0 {
		msg = fmt.Sprintf(s.fOut(citizenTraders, "199", s.state.nonFarmer, s.state.tradeGoods))
		fmt.Print(msg)
	}

	msg = fmt.Sprintf(s.fOut(rats, "199", s.state.pests))
	fmt.Print(msg)
	msg = fmt.Sprintf(s.fOut(storedBushels, "199", s.state.bushels))
	fmt.Print(msg)
	msg = fmt.Sprintf(s.fOut(plows, "199", s.state.plows))
	fmt.Print(msg)
	msg = fmt.Sprintf(s.fOut(landValue, "199", s.state.tradeVal))
	fmt.Print(msg)
	s.state.year += 1
}

func (s *gameSession) doNumbers() (bool, bool) {
	plague := s.checkForPlague()
	var palaceComplete bool
	s.state.forceSlaughtered = 0 // reset this each turn

	if !plague && s.state.buildingPalace > -1 {
		s.state.buildingPalace++
	}
	if s.state.buildingPalace > 4 {
		switch s.palaceBuilding {
		case 1:
			s.state.palace1 = true
		case 2:
			s.state.palace2 = true
		case 3:
			s.state.palace3 = true
		}
		palaceComplete = true
		s.state.buildingPalace = -1
		s.palaceBuilding = -1
	}

	s.state.tradeVal = 17 + rand.Intn(10)
	s.state.bYield = rand.Intn(9) + 1
	// cows
	s.doCows()
	// starvation & population
	s.doPopulation(plague)
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

	// trade is reduced during plague
	if plague {
		s.state.tradeGoods = s.state.nonFarmer * (rand.Intn(15) + 1)
	} else {
		s.state.tradeGoods = s.state.nonFarmer * (rand.Intn(49) + 1)
	}

	s.state.bushels += s.state.tradeGoods
	s.totalDead += s.state.starved
	s.avgPestEaten += s.state.pests
	s.avgBushelsAvail += s.state.bushels
	return plague, palaceComplete
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
			// TODO: messages shouldn't be happening during state changes
			msg := fmt.Sprintf(s.fOut("\nDue to a lack of peasants to work the land, %d acres have wasted and are lost!\n", "196", s.state.acresWastage))
			fmt.Print(msg)
			s.state.acres -= s.state.acresWastage
			s.totAcresWasted += s.state.acresWastage
			if s.state.acres < royalLands {
				s.state.acres = royalLands
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
	unprotectedGrain := s.state.bushels - s.state.granary*granaryProtectMultiplier
	if unprotectedGrain < 0 {
		unprotectedGrain = 0
	}
	s.state.pests = int(float64(unprotectedGrain) / float64(rand.Intn(4)+3))
}

func (s *gameSession) doMigration(plague bool) {
	var cowMigrantAttraction int
	switch {
	case s.state.cows <= 3:
		cowMigrantAttraction = 0
	case s.state.cows > 3 && s.state.population <= 500:
		cowMigrantAttraction = s.state.cows * 5
	case s.state.cows > 3 && s.state.population <= 10000:
		cowMigrantAttraction = s.state.cows * 3
	case s.state.cows > 3 && s.state.population > 10000:
	default:
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

func (s *gameSession) doPopulation(plague bool) {
	carryingCapacity := 100000
	s.state.starved = s.state.population - (s.state.popFed + s.state.cows*s.state.cowMultiplier)
	if s.state.starved < 0 {
		s.state.starved = 0
	}
	s.avgStarved = int(float64(s.state.starved) / float64(s.state.population) * 100)
	switch {
	case s.state.population < 2000:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(8)+2))
	case s.state.population < carryingCapacity/20:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(7)+4))
	case s.state.population < carryingCapacity/10:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(6)+5))
	case s.state.population < carryingCapacity/5:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(6)+6))
	case s.state.population < carryingCapacity/2:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(5)+7))
	case s.state.population < carryingCapacity-carryingCapacity/3:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(5)+8))
	case s.state.population < carryingCapacity-carryingCapacity/5:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(4)+9))
	case s.state.population < carryingCapacity-carryingCapacity/10:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(4)+10))
	case s.state.population < carryingCapacity-carryingCapacity/20:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(3)+11))
	case s.state.population < carryingCapacity:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(2)+12))
	case s.state.population > carryingCapacity:
		s.state.born = int(float64(s.state.population) / float64(rand.Intn(2)+14))
	}
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
	const (
		deposedMsg         = "\nYou starved %d out of your population of only %d, this has caused you to be deposed by force!\n"
		populationDeclined = "\nYour continued mismanagement caused your population to decline to the point that the " +
			"remaining peasants fled your land\nYou are left ruling an empty city, as your royal guards and staff escape.\n"
	)
	if s.state.starved > int(0.45*float64(s.state.population)) {
		msg := fmt.Sprintf(s.fOut(deposedMsg, "196",
			s.state.starved, s.state.population))
		fmt.Print(msg)
		s.totalDead += s.state.starved
		s.endOfReign()
	}

	if s.state.population < 10 {
		fmt.Print(termenv.String(populationDeclined).Bold().Background(s.p.Color("196")))
		s.state.population = 0
		s.endOfReign()
	}
}

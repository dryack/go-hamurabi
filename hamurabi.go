package main

import (
	"fmt"
	"github.com/muesli/termenv"
	"math/rand"
	"time"
)

func main() {
	// mode, err := termenv.()
	// fmt.Println(mode)
	// if err != nil {
	//	panic(err)
	// }
	// defer termenv.RestoreWindowsConsole(mode)
	profile := termenv.ColorProfile()
	fmt.Println(profile)
	test := termenv.String("hello world")
	test.Background(profile.Color("75")).Foreground(profile.Color("52")).Blink()

	fmt.Println(test)

	gameLoop()
}

func gameLoop() {
	rand.Seed(time.Now().UnixNano())
	sumer, test := newGameSession()

	if !test {
		orientation()
		for t := 0; t <= sumer.turns; t++ {
			sumer.printYearResults()
			sumer.getAcres()
			sumer.construction()
			sumer.technology()
			sumer.feedPeople()
			sumer.agriculture()
		}
		sumer.endOfReign()
	}
	for t := 0; t <= sumer.turns; t++ {
		sumer.printYearResults()
		sumer.getAcres()
		sumer.construction()
		sumer.technology()
		sumer.feedPeople()
		sumer.agriculture()
	}
	sumer.endOfReign()
}

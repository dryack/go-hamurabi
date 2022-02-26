package main

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
)

func playerInput(prompt string, defChoice int, maxVal int, failMsg string) int {
	var res string
	os := runtime.GOOS

	// pc, _, _, ok := runtime.Caller(1) // DEBUG
	// details := runtime.FuncForPC(pc)
	// if ok && details != nil {
	//	fmt.Printf("called from %s\n", details.Name())
	// }

	finalPrompt := prompt + " [" + strconv.Itoa(defChoice) + "]" + " => "
	fmt.Print(finalPrompt)
	if os == "windows" {
		_, _ = fmt.Scanf("%s\n", &res)
	} else {
		_, _ = fmt.Scanf("%s", &res)
	}
	if res == "" {
		res = strconv.Itoa(defChoice)
	}
	choice, err := checkInput(res, maxVal)
	if errors.Is(err, strconv.ErrSyntax) {
		fmt.Println("Have you gone mad Hamurabi?! Try again.")
		return playerInput(prompt, defChoice, maxVal, failMsg)
	} else if err != nil {
		fmt.Println(failMsg)
		return playerInput(prompt, defChoice, maxVal, failMsg)
	}
	return choice
}

func checkInput(input string, maxVal int) (int, error) {
	choice, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}
	if choice < 0 || choice > maxVal {
		return 0, fmt.Errorf("invalid choice")
	}
	return choice, nil
}

func (s *gameSession) grainRemaining(res int) {
	if res == 0 {
		return
	}
	fmt.Printf("You have %d bushels of grain remaining.\n", s.state.bushels)
}

func orientation() {
	fmt.Print("\nOh great Hamurabi, congratulations on your accession to the rule of our beautiful Sumer! I am here" +
		" to provide you with the knowledge you need to help us thrive.\n")

	enterToCont()
	fmt.Print("Each year I will present you with an accounting of our prior year's incomes and expenditures. " +
		"\nI will explain the manner in which our population changed, the number of acres of land we hold, the number of " +
		"cows we have available to help feed our population, and how many people they fed. \nIn addition I will tell you " +
		"the number of bushels of barley we reaped per acre of planted land. Finally I will inform you of the toll taken" +
		" on our grain reserves by the rats.\n")
	enterToCont()
	fmt.Print("Oh dread lord Hamurabi, I beg you to know our people each require 20 bushels of grain a year or " +
		"else they will starve. \nThe milk, butter, and cheese provided by our cows will feed 15 people a year for each" +
		" cow, but no income is gained in this way. Cows require 3 acres of land each on which to graze.\nDuring hard" +
		" times, you may order the slaughter of cows, and in this way feed many people.\n")
	enterToCont()
	fmt.Print("Acres may be purchased from land-owning peasants at a price which varies each year. \nOur people may " +
		"only plant 10 acres of land each, although we have great hopes that our strongest men may come to plant more " +
		"as we begin to make use of hand plows.\nToo much land untended by peasants is liable to be reclaimed by nature.\n")
	fmt.Print("Our priests tell us that Enki has delivered unto them the plans for " +
		"a building that will protect our barley from the rats, but warn construction will be expensive.\nPlows multiply" +
		"the labor of our adult male peasants, freeing others to trade and craft. This may have significant impact on" +
		"our income.")
	enterToCont()
	fmt.Print("You must know that starvation; whether caused by cruelty, plague, or pest, will eventually" +
		" lead to the people rising up and deposing you. \nSometimes it is best to let some starve today, in the hope" +
		" that the gods and your steady hand will feed us all in the future. And now my lord, you are ready to rule!\n")
	enterToCont()
}

func enterToCont() {
	os := runtime.GOOS
	fmt.Print("<ENTER> to continue\n")
	if os == "windows" {
		_, _ = fmt.Scanf("%s\n", nil)
	} else {
		_, _ = fmt.Scanf("%s", nil)
	}
}

func yn(prompt string) bool {
	var res string
	os := runtime.GOOS
	fmt.Print(prompt, "[y,n]? => ")
	if os == "windows" {
		_, _ = fmt.Scanf("%s\n", &res)
	} else {
		_, _ = fmt.Scanf("%s", &res)
	}
	r, _ := regexp.Compile("([yYnN])")
	if !r.MatchString(res) {
		fmt.Println("My lord you are incoherent, I need a yes or no to proceed!")
		return yn(prompt)
	}
	if res == "y" || res == "Y" {
		return true
	} else {
		return false
	}
}

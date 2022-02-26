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

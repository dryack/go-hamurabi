package main

import (
	"errors"
	"fmt"
	"github.com/mattn/go-tty"
	"github.com/muesli/termenv"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func playerInput(prompt string, defChoice int, maxVal int, failMsg string, verb string) int {
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
		time.Sleep(2 * time.Second)
		termenv.ClearLines(2)
		return playerInput(prompt, defChoice, maxVal, failMsg, verb)
	} else if err != nil {
		fmt.Println(failMsg)
		time.Sleep(2 * time.Second)
		termenv.ClearLines(2)
		return playerInput(prompt, defChoice, maxVal, failMsg, verb)
	}

	defer func(choice int, verb string) {
		if choice != 0 {
			p := termenv.EnvColorProfile()
			fmt.Println(termenv.String("You " + verb + " " + strconv.Itoa(choice)).Bold().Foreground(p.Color("40")))
		}
	}(choice, verb)
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

func enterToCont() bool {
	fmt.Print("<ENTER> to continue, <Q> to quit")
	t, err := tty.Open()
	if err != nil {
		fmt.Println(err)
	}
	defer t.Close()

	for {
		r, err := t.ReadRune()
		if err != nil {
			fmt.Println(err)
		}
		switch {
		case r == '\r' || r == '\n':
			termenv.ClearScreen()
			return false
		case r == 'Q' || r == 'q':
			termenv.ClearScreen()
			return true
		}
	}
}

func yn(prompt string) bool {
	var res string
	os := runtime.GOOS
	fmt.Print(prompt, " (y,n) ? [n] => ")
	if os == "windows" {
		_, _ = fmt.Scanf("%s\n", &res)
	} else {
		_, _ = fmt.Scanf("%s", &res)
	}
	if res == "" {
		return false
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

//
func (s *gameSession) fOut(str string, clr string, nums ...int) string {
	ansiFmtStr := make([]string, 0)
	for _, n := range nums {
		ansiFmtStr = append(ansiFmtStr, termenv.String(strconv.Itoa(n)).Bold().Foreground(s.p.Color(clr)).String())
	}
	for _, f := range ansiFmtStr {
		str = strings.Replace(str, "%d", f, 1)
	}
	return str
}
package common

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"math"
	"os"
	"strings"
)

const indentChars = "  "

type (
	ProgressUpdate struct {
		Key     string `json:"key"`
		Message string `json:"msg"`
		Count   int    `json:"c"`
	}
)

var (
	NoProgress         = false
	lastProgressUpdate = ProgressUpdate{}
	completion         = 0
	termWidth          = 80
)

func init() {
	w, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		// Assume 80 columns
		termWidth = 80
	} else {
		termWidth = w
	}
}

func ShowProgress(u ProgressUpdate) {
	if NoProgress {
		return
	}

	if completion > 0 && u.Key != "" && u.Key != lastProgressUpdate.Key {
		completion = 0
		switchMessage("✓", color.FgHiGreen)
		showMessage(u)
	} else {
		clearLine()
		if u.Key == "" {
			// Detail message
			showMessage(lastProgressUpdate)
			color.New(color.FgWhite).Printf(" %s", u.Message)
		} else {
			completion++
			// Normal message
			showMessage(u)
		}
	}

	if u.Key != "" {
		lastProgressUpdate = u
	}
}

func ShowPrompt(msg string, v ...interface{}) string {
	if NoProgress {
		return ""
	}

	switchMessage("✓", color.FgHiGreen)
	color.New(color.FgHiWhite).Printf("%s...\n", fmt.Sprintf(msg, v...))
	resetProgress()

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		ShowError("Unable to read user input")
		os.Exit(1)
	}
	return text
}

func ShowWorking(msg string, v ...interface{}) {
	if NoProgress {
		return
	}

	switchMessage("✓", color.FgHiGreen)
	color.New(color.FgHiWhite).Printf("%s...\n", fmt.Sprintf(msg, v...))
	resetProgress()
}

func ShowError(msg string, v ...interface{}) {
	if NoProgress {
		return
	}

	switchMessage("✗", color.FgHiRed)
	color.New(color.FgRed).Printf("%s\n", fmt.Sprintf(msg, v...))
	resetProgress()
}

func ShowDone(msg string, v ...interface{}) {
	if NoProgress {
		return
	}

	switchMessage("✓", color.FgHiGreen)
	color.New(color.FgGreen).Printf("%s\n", fmt.Sprintf(msg, v...))
	resetProgress()
}

func clearLine() {
	fmt.Printf("\r%s\r", strings.Repeat(" ", termWidth))
}

func resetProgress() {
	lastProgressUpdate = ProgressUpdate{}
	completion = 0
}

func switchMessage(symbol string, attr color.Attribute) {
	if lastProgressUpdate.Message != "" {
		clearLine()
		color.New(color.FgHiWhite).Printf("%s%s [", indentChars, lastProgressUpdate.Message)
		color.New(attr).Print(symbol)
		color.New(color.FgHiWhite).Printf("]\n")
	}
}

func showMessage(u ProgressUpdate) {
	color.New(color.FgHiWhite).Printf("%s%s... [", indentChars, u.Message)
	if u.Count > 1 {
		// We have a goal to calculate percentage to
		var percentage int
		if completion == 0 {
			percentage = 0
		} else {
			percentage = int(math.Min(math.Abs(float64(completion)/float64(u.Count))*100, 100))
		}
		color.New(color.FgHiYellow).Printf("%d%%", percentage)
	} else {
		// No goal, just display an hourglass
		color.New(color.FgHiYellow).Print("⧖")
	}
	color.New(color.FgHiWhite).Print("]")
}

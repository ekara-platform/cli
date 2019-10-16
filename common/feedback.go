package common

import (
	"bufio"
	"fmt"
	"github.com/ekara-platform/engine/util"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"math"
	"os"
	"strings"
)

type (
	consoleFeedbackNotifier struct {
		util.FeedbackNotifier
		last struct {
			key     string
			message string
			count   int
		}
		completion int
	}
)

var (
	CliFeedbackNotifier = &consoleFeedbackNotifier{}
	NoFeedback          = false
	termWidth           = 80
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

func (l *consoleFeedbackNotifier) Error(msg string, v ...interface{}) {
	if NoFeedback {
		return
	}
	l.switchMessage("✗", color.FgHiRed)
	color.New(color.FgRed).Printf("%s\n", fmt.Sprintf(msg, v...))
}

func (l *consoleFeedbackNotifier) Info(msg string, v ...interface{}) {
	if NoFeedback {
		return
	}
	l.switchMessage("✓", color.FgHiGreen)
	color.New(color.FgGreen).Printf("%s\n", fmt.Sprintf(msg, v...))
}

func (l *consoleFeedbackNotifier) Progress(key string, message string, v ...interface{}) {
	l.ProgressG(key, 1, message, v...)
}

func (l *consoleFeedbackNotifier) ProgressG(key string, goal int, message string, v ...interface{}) {
	if NoFeedback {
		return
	}

	var formattedMessage = fmt.Sprintf(message, v...)
	if key != "" && key != l.last.key {
		l.completion = 0
		l.switchMessage("✓", color.FgHiGreen)
		l.displayMessage(formattedMessage, goal)
	} else {
		l.completion++
		l.clearLine()
		l.displayMessage(formattedMessage, goal)
	}

	l.last.key = key
	l.last.count = goal
	l.last.message = formattedMessage
}

func (l *consoleFeedbackNotifier) Detail(message string, v ...interface{}) {
	l.clearLine()
	l.displayMessage(l.last.message, l.last.count)
	color.New(color.FgWhite).Printf(" %s", fmt.Sprintf(message, v...))
}

func (l *consoleFeedbackNotifier) Prompt(msg string, v ...interface{}) string {
	if NoFeedback {
		return ""
	}

	l.switchMessage("✓", color.FgHiGreen)
	color.New(color.FgHiWhite).Printf("%s...\n", fmt.Sprintf(msg, v...))

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		l.Error("Unable to read user input")
		os.Exit(1)
	}
	return text
}

func (l *consoleFeedbackNotifier) reset() {
	l.last.key = ""
	l.last.message = ""
	l.last.count = 0
	l.completion = 0
}

func (l *consoleFeedbackNotifier) switchMessage(symbol string, attr color.Attribute) {
	if l.last.message != "" {
		l.clearLine()
		color.New(color.FgHiWhite).Printf("%s [", l.last.message)
		color.New(attr).Print(symbol)
		color.New(color.FgHiWhite).Printf("]\n")
	}
	l.reset()
}

func (l *consoleFeedbackNotifier) displayMessage(message string, goal int) {
	color.New(color.FgHiWhite).Printf("%s... [", message)
	if goal == 0 {
		// no goal: info
		color.New(color.FgHiBlue).Print("🛈")
	} else if goal == 1 {
		// no progression: 0% or 100% only
		color.New(color.FgHiYellow).Print("⧖")
	} else {
		// progression: calculate percentage
		var percentage int
		if l.completion == 0 {
			percentage = 0
		} else {
			percentage = int(math.Min(math.Abs(float64(l.completion)/float64(goal))*100, 100))
		}
		color.New(color.FgHiYellow).Printf("%d%%", percentage)
	}
	color.New(color.FgHiWhite).Print("]")
}

func (l *consoleFeedbackNotifier) clearLine() {
	fmt.Printf("\r%s\r", strings.Repeat(" ", termWidth))
}

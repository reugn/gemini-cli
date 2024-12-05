package terminal

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"
	"github.com/reugn/gemini-cli/internal/terminal/color"
)

const (
	geminiUser = "gemini"
	cliUser    = "cli"
)

type Prompt struct {
	User     string
	UserNext string
	Gemini   string
	Cli      string
}

type promptColor struct {
	user   func(string) string
	gemini func(string) string
	cli    func(string) string
}

func newPromptColor() *promptColor {
	if termenv.HasDarkBackground() {
		return &promptColor{
			user:   color.Cyan,
			gemini: color.Green,
			cli:    color.Yellow,
		}
	}
	return &promptColor{
		user:   color.Blue,
		gemini: color.Green,
		cli:    color.Magenta,
	}
}

func NewPrompt(currentUser string) *Prompt {
	maxLength := maxLength(currentUser, geminiUser, cliUser)
	pc := newPromptColor()
	return &Prompt{
		User:     pc.user(buildPrompt(currentUser, maxLength)),
		UserNext: pc.user(buildPrompt(strings.Repeat(" ", len(currentUser)), maxLength)),
		Gemini:   pc.gemini(buildPrompt(geminiUser, maxLength)),
		Cli:      pc.cli(buildPrompt(cliUser, maxLength)),
	}
}

func maxLength(str ...string) int {
	var maxLength int
	for _, s := range str {
		length := len(s)
		if maxLength < length {
			maxLength = length
		}
	}
	return maxLength
}

func buildPrompt(user string, length int) string {
	return fmt.Sprintf("%s>%s", user, strings.Repeat(" ", length-len(user)+1))
}

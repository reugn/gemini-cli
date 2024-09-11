package cli

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"
	"github.com/reugn/gemini-cli/internal/cli/color"
)

const (
	geminiUser = "gemini"
	cliUser    = "cli"
)

type prompt struct {
	user     string
	userNext string
	gemini   string
	cli      string
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

func newPrompt(currentUser string) *prompt {
	maxLength := maxLength(currentUser, geminiUser, cliUser)
	pc := newPromptColor()
	return &prompt{
		user:     pc.user(buildPrompt(currentUser, maxLength)),
		userNext: pc.user(buildPrompt(strings.Repeat(" ", len(currentUser)), maxLength)),
		gemini:   pc.gemini(buildPrompt(geminiUser, maxLength)),
		cli:      pc.cli(buildPrompt(cliUser, maxLength)),
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

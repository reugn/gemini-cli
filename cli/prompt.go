package cli

import (
	"fmt"
	"strings"

	"github.com/reugn/gemini-cli/cli/color"
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

func newPrompt(currentUser string) *prompt {
	maxLength := maxLength(currentUser, geminiUser, cliUser)
	return &prompt{
		user:     color.Blue(buildPrompt(currentUser, maxLength)),
		userNext: color.Blue(buildPrompt(strings.Repeat(" ", len(currentUser)), maxLength)),
		gemini:   color.Green(buildPrompt(geminiUser, maxLength)),
		cli:      color.Yellow(buildPrompt(cliUser, maxLength)),
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

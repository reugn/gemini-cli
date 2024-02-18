package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/reugn/gemini-cli/gemini"
)

// ChatOpts represents Chat configuration options.
type ChatOpts struct {
	Format bool
	Style  string
}

// Chat controls the chat flow.
type Chat struct {
	model  *gemini.ChatSession
	prompt *prompt
	reader *bufio.Reader
	opts   *ChatOpts
}

// NewChat returns a new Chat.
func NewChat(user string, model *gemini.ChatSession, opts *ChatOpts) *Chat {
	return &Chat{
		model:  model,
		prompt: newPrompt(user),
		reader: bufio.NewReader(os.Stdin),
		opts:   opts,
	}
}

// StartChat starts the chat loop.
func (c *Chat) StartChat() {
	for {
		fmt.Print(c.prompt.user)
		message, ok := c.readLine()
		if !ok {
			continue
		}
		command := c.parseCommand(message)
		if quit := command.run(message); quit {
			break
		}
	}
}

func (c *Chat) readLine() (string, bool) {
	input, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("%s%s\n", c.prompt.cli, err)
		return "", false
	}
	input = strings.ReplaceAll(input, "\n", "")
	if strings.TrimSpace(input) == "" {
		return "", false
	}
	return input, true
}

func (c *Chat) parseCommand(message string) command {
	if strings.HasPrefix(message, systemCmdPrefix) {
		return newSystemCommand(c.model, c.prompt)
	}
	return newGeminiCommand(c.model, c.prompt, c.opts)
}

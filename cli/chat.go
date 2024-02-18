package cli

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
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
	reader *readline.Instance
	opts   *ChatOpts
}

// NewChat returns a new Chat.
func NewChat(user string, model *gemini.ChatSession, opts *ChatOpts) (*Chat, error) {
	reader, err := readline.NewEx(&readline.Config{})
	if err != nil {
		return nil, err
	}
	prompt := newPrompt(user)
	reader.SetPrompt(prompt.user)
	return &Chat{
		model:  model,
		prompt: prompt,
		reader: reader,
		opts:   opts,
	}, nil
}

// StartChat starts the chat loop.
func (c *Chat) StartChat() {
	for {
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
	input, err := c.reader.Readline()
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

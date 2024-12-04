package chat

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/cli"
	"github.com/reugn/gemini-cli/internal/config"
	"github.com/reugn/gemini-cli/internal/handler"
)

// Chat handles the interactive exchange of messages between user and model.
type Chat struct {
	reader         *readline.Instance
	writer         io.Writer
	terminalPrompt *cli.Prompt
	opts           *Opts

	geminiHandler handler.MessageHandler
	systemHandler handler.MessageHandler
}

// New returns a new Chat.
func New(
	user string, session *gemini.ChatSession,
	configuration *config.Configuration, writer io.Writer, opts *Opts,
) (*Chat, error) {
	reader, err := readline.NewEx(&readline.Config{})
	if err != nil {
		return nil, err
	}

	terminalPrompt := cli.NewPrompt(user)
	reader.SetPrompt(terminalPrompt.User)
	if opts.Multiline {
		// disable history for multiline input mode
		reader.HistoryDisable()
	}

	spinner := cli.NewSpinner(writer, time.Second, 5)
	geminiHandler, err := handler.NewGeminiQuery(session, spinner, opts.Style)
	if err != nil {
		return nil, err
	}

	systemHandler := handler.NewSystemCommand(session, configuration, reader,
		&opts.Multiline, opts.GenerativeModel)

	return &Chat{
		terminalPrompt: terminalPrompt,
		reader:         reader,
		writer:         writer,
		opts:           opts,
		geminiHandler:  geminiHandler,
		systemHandler:  systemHandler,
	}, nil
}

// Start starts the main chat loop between user and model.
func (c *Chat) Start() {
	for {
		message, ok := c.read()
		if !ok {
			continue
		}

		// process the message
		messageHandler, terminalPrompt := c.getHandlerPrompt(message)
		response, quit := messageHandler.Handle(message)
		_ = response.Print(c.writer, terminalPrompt)

		if quit {
			break
		}
	}
}

func (c *Chat) read() (string, bool) {
	if c.opts.Multiline {
		return c.readMultiLine()
	}
	return c.readLine()
}

func (c *Chat) readLine() (string, bool) {
	input, err := c.reader.Readline()
	if err != nil {
		return c.handleReadError(len(input), err)
	}
	return validateInput(input)
}

func (c *Chat) readMultiLine() (string, bool) {
	var builder strings.Builder
	term := c.opts.LineTerminator
	for {
		input, err := c.reader.Readline()
		if err != nil {
			c.reader.SetPrompt(c.terminalPrompt.User)
			return c.handleReadError(builder.Len()+len(input), err)
		}

		if strings.HasSuffix(input, term) ||
			strings.HasPrefix(input, handler.SystemCmdPrefix) {
			builder.WriteString(strings.TrimSuffix(input, term))
			break
		}

		if builder.Len() == 0 {
			c.reader.SetPrompt(c.terminalPrompt.UserNext)
		}

		builder.WriteString(input + "\n")
	}
	c.reader.SetPrompt(c.terminalPrompt.User)
	return validateInput(builder.String())
}

func (c *Chat) handleReadError(inputLen int, err error) (string, bool) {
	if errors.Is(err, readline.ErrInterrupt) {
		if inputLen == 0 {
			return handler.SystemCmdPrefix + handler.SystemCmdQuit, true
		}
	} else {
		handler.PrintError(c.writer, c.terminalPrompt.Cli, err)
	}
	return "", false
}

func (c *Chat) getHandlerPrompt(message string) (handler.MessageHandler, string) {
	if strings.HasPrefix(message, handler.SystemCmdPrefix) {
		return c.systemHandler, c.terminalPrompt.Cli
	}
	_, _ = fmt.Fprint(c.writer, c.terminalPrompt.Gemini)
	return c.geminiHandler, ""
}

func validateInput(input string) (string, bool) {
	input = strings.TrimSpace(input)
	return input, input != ""
}

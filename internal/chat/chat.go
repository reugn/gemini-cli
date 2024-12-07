package chat

import (
	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/cli"
	"github.com/reugn/gemini-cli/internal/config"
	"github.com/reugn/gemini-cli/internal/handler"
	"github.com/reugn/gemini-cli/internal/terminal"
)

// Chat handles the interactive exchange of messages between user and model.
type Chat struct {
	io *terminal.IO

	geminiHandler handler.MessageHandler
	systemHandler handler.MessageHandler
}

// New returns a new Chat.
func New(
	user string, session *gemini.ChatSession,
	configuration *config.Configuration, opts *Opts,
) (*Chat, error) {
	terminalIOConfig := &terminal.IOConfig{
		User:           user,
		Multiline:      opts.Multiline,
		LineTerminator: opts.LineTerminator,
	}

	terminalIO, err := terminal.NewIO(terminalIOConfig)
	if err != nil {
		return nil, err
	}

	geminiIO := handler.NewIO(terminalIO, terminalIO.Prompt.Gemini)
	geminiHandler, err := handler.NewGeminiQuery(geminiIO, session, opts.rendererOptions())
	if err != nil {
		return nil, err
	}

	systemIO := handler.NewIO(terminalIO, terminalIO.Prompt.Cli)
	systemHandler, err := handler.NewSystemCommand(systemIO, session, configuration,
		opts.GenerativeModel, opts.rendererOptions())
	if err != nil {
		return nil, err
	}

	return &Chat{
		io:            terminalIO,
		geminiHandler: geminiHandler,
		systemHandler: systemHandler,
	}, nil
}

// Start starts the main chat loop between user and model.
func (c *Chat) Start() {
	for {
		// read query from the user
		message := c.io.Read()
		if message == "" {
			continue
		}

		// get handler for the read message
		// the message is not empty here
		messageHandler := c.getHandler(message[:1])

		// write the agent terminal prompt
		c.io.Write(messageHandler.TerminalPrompt())

		// process the message
		response, quit := messageHandler.Handle(message)

		// write the response
		c.io.Write(response.String())

		if quit {
			break
		}
	}
}

// getHandler returns the handler for the message.
func (c *Chat) getHandler(prefix string) handler.MessageHandler {
	if prefix == cli.SystemCmdPrefix {
		return c.systemHandler
	}
	return c.geminiHandler
}

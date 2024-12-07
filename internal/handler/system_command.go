package handler

import (
	"fmt"
	"strings"

	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/cli"
	"github.com/reugn/gemini-cli/internal/config"
)

const (
	empty            = "Empty"
	unchangedMessage = "The selection is unchanged."
)

// SystemCommand processes chat system commands; implements the MessageHandler interface.
// It aggregates the processing by delegating it to one of the underlying handlers.
type SystemCommand struct {
	*IO
	handlers map[string]MessageHandler
}

var _ MessageHandler = (*SystemCommand)(nil)

// NewSystemCommand returns a new SystemCommand.
func NewSystemCommand(io *IO, session *gemini.ChatSession, configuration *config.Configuration,
	modelName string, rendererOptions RendererOptions) (*SystemCommand, error) {
	helpCommandHandler, err := NewHelpCommand(io, rendererOptions)
	if err != nil {
		return nil, err
	}

	handlers := map[string]MessageHandler{
		cli.SystemCmdHelp:            helpCommandHandler,
		cli.SystemCmdQuit:            NewQuitCommand(io),
		cli.SystemCmdSelectPrompt:    NewSystemPromptCommand(io, session, configuration.Data),
		cli.SystemCmdSelectInputMode: NewInputModeCommand(io),
		cli.SystemCmdModel:           NewModelCommand(io, session, modelName),
		cli.SystemCmdHistory:         NewHistoryCommand(io, session, configuration),
	}

	return &SystemCommand{
		IO:       io,
		handlers: handlers,
	}, nil
}

// Handle processes the chat system command.
func (s *SystemCommand) Handle(message string) (Response, bool) {
	if !strings.HasPrefix(message, cli.SystemCmdPrefix) {
		return newErrorResponse(fmt.Errorf("system command mismatch")), false
	}

	var args string
	t := strings.SplitN(message, " ", 2)
	if len(t) == 2 {
		args = t[1]
	}

	systemHandler, ok := s.handlers[message[1:]]
	if !ok {
		return newErrorResponse(fmt.Errorf("unknown system command")), false
	}

	return systemHandler.Handle(args)
}

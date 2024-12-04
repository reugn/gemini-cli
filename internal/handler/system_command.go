package handler

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/config"
)

const (
	SystemCmdPrefix          = "!"
	SystemCmdQuit            = "q"
	systemCmdSelectPrompt    = "p"
	systemCmdSelectInputMode = "i"
	systemCmdModel           = "m"
	systemCmdHistory         = "h"
)

// SystemCommand processes chat system commands; implements the MessageHandler interface.
// It aggregates the processing by delegating it to one of the underlying handlers.
type SystemCommand struct {
	handlers map[string]MessageHandler
}

var _ MessageHandler = (*SystemCommand)(nil)

// NewSystemCommand returns a new SystemCommand.
func NewSystemCommand(session *gemini.ChatSession, configuration *config.Configuration,
	reader *readline.Instance, multiline *bool, modelName string) *SystemCommand {
	handlers := map[string]MessageHandler{
		SystemCmdQuit:            NewQuitCommand(),
		systemCmdSelectPrompt:    NewSystemPromptCommand(session, configuration.Data),
		systemCmdSelectInputMode: NewInputModeCommand(reader, multiline),
		systemCmdModel:           NewModelCommand(session, modelName),
		systemCmdHistory:         NewHistoryCommand(session, configuration),
	}

	return &SystemCommand{
		handlers: handlers,
	}
}

// Handle processes the chat system command.
func (s *SystemCommand) Handle(message string) (Response, bool) {
	if !strings.HasPrefix(message, SystemCmdPrefix) {
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

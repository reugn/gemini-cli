package handler

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

var inputModeOptions = []string{
	"Single-line",
	"Multi-line",
}

// InputModeCommand processes the chat input mode system command.
// It implements the MessageHandler interface.
type InputModeCommand struct {
	*IO
}

var _ MessageHandler = (*InputModeCommand)(nil)

// NewInputModeCommand returns a new InputModeCommand.
func NewInputModeCommand(io *IO) *InputModeCommand {
	return &InputModeCommand{
		IO: io,
	}
}

// Handle processes the chat input mode system command.
func (h *InputModeCommand) Handle(_ string) (Response, bool) {
	defer h.terminal.Write(h.terminalPrompt)
	multiline, err := h.selectInputMode()
	if err != nil {
		return newErrorResponse(err), false
	}

	if h.terminal.Config.Multiline == multiline {
		// the same input mode is selected
		return dataResponse(unchangedMessage), false
	}

	h.terminal.Config.Multiline = multiline
	h.terminal.SetUserPrompt()
	if h.terminal.Config.Multiline {
		// disable history for multi-line messages since it is
		// unusable for future requests
		h.terminal.Reader.HistoryDisable()
	} else {
		h.terminal.Reader.HistoryEnable()
	}

	mode := inputModeOptions[modeIndex(h.terminal.Config.Multiline)]
	return dataResponse(fmt.Sprintf("Switched to %q input mode.", mode)), false
}

// selectInputMode returns true if multiline input is selected;
// otherwise, it returns false.
func (h *InputModeCommand) selectInputMode() (bool, error) {
	prompt := promptui.Select{
		Label:        "Select input mode",
		HideSelected: true,
		Items:        inputModeOptions,
		CursorPos:    modeIndex(h.terminal.Config.Multiline),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == inputModeOptions[1], nil
}

func modeIndex(b bool) int {
	if b {
		return 1
	}
	return 0
}

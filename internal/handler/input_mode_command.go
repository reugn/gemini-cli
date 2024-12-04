package handler

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
)

var inputModeOptions = []string{
	"Single-line",
	"Multi-line",
}

// InputModeCommand processes the chat input mode system command.
// It implements the MessageHandler interface.
type InputModeCommand struct {
	reader    *readline.Instance
	multiline *bool
}

var _ MessageHandler = (*InputModeCommand)(nil)

// NewInputModeCommand returns a new InputModeCommand.
func NewInputModeCommand(reader *readline.Instance, multiline *bool) *InputModeCommand {
	return &InputModeCommand{
		reader:    reader,
		multiline: multiline,
	}
}

// Handle processes the chat input mode system command.
func (h *InputModeCommand) Handle(_ string) (Response, bool) {
	multiline, err := h.selectInputMode()
	if err != nil {
		return newErrorResponse(err), false
	}

	if *h.multiline == multiline {
		// the same input mode is selected
		return dataResponse(unchangedMessage), false
	}

	*h.multiline = multiline
	if *h.multiline {
		// disable history for multi-line messages since it is
		// unusable for future requests
		h.reader.HistoryDisable()
	} else {
		h.reader.HistoryEnable()
	}

	mode := inputModeOptions[modeIndex(*h.multiline)]
	return dataResponse(fmt.Sprintf("Switched to %q input mode.", mode)), false
}

// selectInputMode returns true if multiline input is selected;
// otherwise, it returns false.
func (h *InputModeCommand) selectInputMode() (bool, error) {
	prompt := promptui.Select{
		Label:        "Select input mode",
		HideSelected: true,
		Items:        inputModeOptions,
		CursorPos:    modeIndex(*h.multiline),
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

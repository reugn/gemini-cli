package handler

import (
	"fmt"
	"slices"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/config"
	"google.golang.org/genai"
)

var historyOptions = []string{
	"Clear chat history",
	"Store chat history",
	"Load chat history",
	"Delete stored history records",
}

// HistoryCommand processes the chat history system commands.
// It implements the MessageHandler interface.
type HistoryCommand struct {
	*IO
	session       *gemini.ChatSession
	configuration *config.Configuration
}

var _ MessageHandler = (*HistoryCommand)(nil)

// NewHistoryCommand returns a new HistoryCommand.
func NewHistoryCommand(io *IO, session *gemini.ChatSession,
	configuration *config.Configuration) *HistoryCommand {
	return &HistoryCommand{
		IO:            io,
		session:       session,
		configuration: configuration,
	}
}

// Handle processes the history system command.
func (h *HistoryCommand) Handle(_ string) (Response, bool) {
	option, err := h.selectHistoryOption()
	if err != nil {
		return newErrorResponse(err), false
	}
	var response Response
	switch option {
	case historyOptions[0]:
		response = h.handleClear()
	case historyOptions[1]:
		response = h.handleStore()
	case historyOptions[2]:
		response = h.handleLoad()
	case historyOptions[3]:
		response = h.handleDelete()
	default:
		response = newErrorResponse(fmt.Errorf("unsupported option: %s", option))
	}
	return response, false
}

// handleClear handles the chat history clear request.
func (h *HistoryCommand) handleClear() Response {
	h.terminal.Write(h.terminalPrompt)
	if err := h.session.ClearHistory(); err != nil {
		return newErrorResponse(err)
	}

	return dataResponse("Cleared the chat history.")
}

// handleStore handles the chat history store request.
func (h *HistoryCommand) handleStore() Response {
	defer h.terminal.Write(h.terminalPrompt)
	historyLabel, err := h.promptHistoryLabel()
	if err != nil {
		return newErrorResponse(err)
	}

	timeLabel := time.Now().In(time.Local).Format(time.DateTime)
	recordLabel := fmt.Sprintf("%s - %s", timeLabel, historyLabel)

	h.configuration.Data.AddHistoryRecord(
		recordLabel,
		h.session.GetHistory(),
	)

	if err := h.configuration.Flush(); err != nil {
		return newErrorResponse(err)
	}

	return dataResponse(fmt.Sprintf("%q has been saved to the file.", recordLabel))
}

// handleLoad handles the chat history load request.
func (h *HistoryCommand) handleLoad() Response {
	defer h.terminal.Write(h.terminalPrompt)
	label, history, err := h.loadHistory()
	if err != nil {
		return newErrorResponse(err)
	}

	if err := h.session.SetHistory(history); err != nil {
		return newErrorResponse(err)
	}

	return dataResponse(fmt.Sprintf("%q has been loaded to the chat history.", label))
}

// handleDelete handles deletion of the stored history records.
func (h *HistoryCommand) handleDelete() Response {
	h.terminal.Write(h.terminalPrompt)
	h.configuration.Data.History = make(map[string][]*gemini.SerializableContent)
	if err := h.configuration.Flush(); err != nil {
		return newErrorResponse(err)
	}
	return dataResponse("History records have been removed from the file.")
}

// loadHistory returns history data to be set.
func (h *HistoryCommand) loadHistory() (string, []*genai.Content, error) {
	promptNames := make([]string, len(h.configuration.Data.History)+1)
	promptNames[0] = empty
	i := 1
	for p := range h.configuration.Data.History {
		promptNames[i] = p
		i++
	}
	prompt := promptui.Select{
		Label:        "Select conversation history to load",
		HideSelected: true,
		Items:        promptNames,
		CursorPos:    slices.Index(promptNames, empty),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return result, nil, err
	}

	if result == empty {
		return result, nil, nil
	}

	serializedContent := h.configuration.Data.History[result]
	content := make([]*genai.Content, len(serializedContent))
	for i, c := range serializedContent {
		content[i] = c.ToContent()
	}

	return result, content, nil
}

// promptHistoryLabel returns a label for the history record.
func (h *HistoryCommand) promptHistoryLabel() (string, error) {
	prompt := promptui.Prompt{
		Label:       "Enter a label for the history record",
		HideEntered: true,
	}

	label, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return label, nil
}

// selectHistoryOption returns the selected history action name.
func (h *HistoryCommand) selectHistoryOption() (string, error) {
	prompt := promptui.Select{
		Label:        "Select history option",
		HideSelected: true,
		Items:        historyOptions,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

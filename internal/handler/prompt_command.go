package handler

import (
	"fmt"
	"slices"

	"github.com/manifoldco/promptui"
	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/config"
	"google.golang.org/genai"
)

// SystemPromptCommand processes the chat prompt system command.
// It implements the MessageHandler interface.
type SystemPromptCommand struct {
	*IO
	session         *gemini.ChatSession
	applicationData *config.ApplicationData

	systemPrompt string
}

var _ MessageHandler = (*SystemPromptCommand)(nil)

// NewSystemPromptCommand returns a new SystemPromptCommand.
func NewSystemPromptCommand(io *IO, session *gemini.ChatSession,
	applicationData *config.ApplicationData) *SystemPromptCommand {
	return &SystemPromptCommand{
		IO:              io,
		session:         session,
		applicationData: applicationData,
	}
}

// Handle processes the chat prompt system command.
func (h *SystemPromptCommand) Handle(_ string) (Response, bool) {
	defer h.terminal.Write(h.terminalPrompt)
	label, systemPrompt, err := h.selectSystemPrompt()
	if err != nil {
		return newErrorResponse(err), false
	}

	if err := h.session.SetSystemInstruction(systemPrompt); err != nil {
		return newErrorResponse(err), false
	}

	return dataResponse(fmt.Sprintf("Selected %q system instruction.", label)), false
}

// selectSystemPrompt returns a system instruction to be set.
func (h *SystemPromptCommand) selectSystemPrompt() (string, *genai.Content, error) {
	promptNames := make([]string, len(h.applicationData.SystemPrompts)+1)
	promptNames[0] = empty
	i := 1
	for p := range h.applicationData.SystemPrompts {
		promptNames[i] = p
		i++
	}
	prompt := promptui.Select{
		Label:        "Select system instruction",
		HideSelected: true,
		Items:        promptNames,
		CursorPos:    slices.Index(promptNames, h.systemPrompt),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return result, nil, err
	}

	h.systemPrompt = result
	if result == empty {
		return result, nil, nil
	}

	systemInstruction := h.applicationData.SystemPrompts[result]
	return result, systemInstruction.ToContent(), nil
}

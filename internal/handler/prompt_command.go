package handler

import (
	"fmt"
	"slices"

	"github.com/google/generative-ai-go/genai"
	"github.com/manifoldco/promptui"
	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/config"
)

// SystemPromptCommand processes the chat prompt system command.
// It implements the MessageHandler interface.
type SystemPromptCommand struct {
	session         *gemini.ChatSession
	applicationData *config.ApplicationData

	currentPrompt string
}

var _ MessageHandler = (*SystemPromptCommand)(nil)

// NewSystemPromptCommand returns a new SystemPromptCommand.
func NewSystemPromptCommand(session *gemini.ChatSession,
	applicationData *config.ApplicationData) *SystemPromptCommand {
	return &SystemPromptCommand{
		session:         session,
		applicationData: applicationData,
	}
}

// Handle processes the chat prompt system command.
func (h *SystemPromptCommand) Handle(_ string) (Response, bool) {
	label, systemPrompt, err := h.selectSystemPrompt()
	if err != nil {
		return newErrorResponse(err), false
	}

	modelBuilder := h.session.CopyModelBuilder().
		WithSystemInstruction(systemPrompt)
	h.session.SetModel(modelBuilder)

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
		CursorPos:    slices.Index(promptNames, h.currentPrompt),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return result, nil, err
	}

	h.currentPrompt = result
	if result == empty {
		return result, nil, nil
	}

	systemInstruction := h.applicationData.SystemPrompts[result]
	return result, systemInstruction.ToContent(), nil
}

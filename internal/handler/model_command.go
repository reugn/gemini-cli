package handler

import (
	"fmt"
	"slices"

	"github.com/manifoldco/promptui"
	"github.com/reugn/gemini-cli/gemini"
)

var modelOptions = []string{
	"Select generative model",
	"Chat model info",
}

// ModelCommand processes the chat model system commands.
// It implements the MessageHandler interface.
type ModelCommand struct {
	session      *gemini.ChatSession
	currentModel string
}

var _ MessageHandler = (*ModelCommand)(nil)

// NewModelCommand returns a new ModelCommand.
func NewModelCommand(session *gemini.ChatSession, modelName string) *ModelCommand {
	return &ModelCommand{
		session:      session,
		currentModel: modelName,
	}
}

// Handle processes the chat model system command.
func (h *ModelCommand) Handle(_ string) (Response, bool) {
	option, err := h.selectModelOption()
	if err != nil {
		return newErrorResponse(err), false
	}

	var response Response
	switch option {
	case modelOptions[0]:
		response = h.handleSelectModel()
	case modelOptions[1]:
		response = h.handleModelInfo()
	default:
		response = newErrorResponse(fmt.Errorf("unsupported option: %s", option))
	}
	return response, false
}

// handleSelectModel handles the generative model selection.
func (h *ModelCommand) handleSelectModel() Response {
	model, err := h.selectModel(h.session.ListModels())
	if err != nil {
		return newErrorResponse(err)
	}

	if h.currentModel == model {
		return dataResponse(unchangedMessage)
	}

	modelBuilder := h.session.CopyModelBuilder().WithName(model)
	h.session.SetModel(modelBuilder)
	h.currentModel = model

	return dataResponse(fmt.Sprintf("Selected %q generative model.", model))
}

// handleSelectModel handles the current generative model info request.
func (h *ModelCommand) handleModelInfo() Response {
	modelInfo, err := h.session.ModelInfo()
	if err != nil {
		return newErrorResponse(err)
	}
	return dataResponse(modelInfo)
}

// selectModelOption returns the selected action name.
func (h *ModelCommand) selectModelOption() (string, error) {
	prompt := promptui.Select{
		Label:        "Select model option",
		HideSelected: true,
		Items:        modelOptions,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// selectModel returns the selected generative model name.
func (h *ModelCommand) selectModel(models []string) (string, error) {
	prompt := promptui.Select{
		Label:        "Select generative session",
		HideSelected: true,
		Items:        models,
		CursorPos:    slices.Index(models, h.currentModel),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

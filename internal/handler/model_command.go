package handler

import (
	"fmt"
	"slices"
	"strings"

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
	*IO
	session             *gemini.ChatSession
	generativeModelName string
}

var _ MessageHandler = (*ModelCommand)(nil)

// NewModelCommand returns a new ModelCommand.
func NewModelCommand(io *IO, session *gemini.ChatSession, modelName string) *ModelCommand {
	return &ModelCommand{
		IO:                  io,
		session:             session,
		generativeModelName: modelName,
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
	defer h.terminal.Write(h.terminalPrompt)
	modelName, err := h.selectModel(h.session.ListModels())
	if err != nil {
		return newErrorResponse(err)
	}

	if h.generativeModelName == modelName {
		return dataResponse(unchangedMessage)
	}

	modelBuilder := h.session.CopyModelBuilder().WithName(modelName)
	h.session.SetModel(modelBuilder)
	h.generativeModelName = modelName

	return dataResponse(fmt.Sprintf("Selected %q generative model.", modelName))
}

// handleSelectModel handles the current generative model info request.
func (h *ModelCommand) handleModelInfo() Response {
	h.terminal.Write(h.terminalPrompt)
	h.terminal.Spinner.Start()
	defer h.terminal.Spinner.Stop()

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
		Label:        modelOptions[0],
		HideSelected: true,
		Items:        models,
		CursorPos:    slices.Index(models, h.generativeModelName),
		Searcher: func(input string, index int) bool {
			return strings.Contains(models[index], input)
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

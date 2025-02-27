package handler

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/reugn/gemini-cli/internal/cli"
)

// HelpCommand handles the help system command request.
type HelpCommand struct {
	*IO
	renderer *glamour.TermRenderer
}

var _ MessageHandler = (*HelpCommand)(nil)

// NewHelpCommand returns a new HelpCommand.
func NewHelpCommand(io *IO, opts RendererOptions) (*HelpCommand, error) {
	renderer, err := opts.newTermRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate terminal renderer: %w", err)
	}

	return &HelpCommand{
		IO:       io,
		renderer: renderer,
	}, nil
}

// Handle processes the help system command.
func (h *HelpCommand) Handle(_ string) (Response, bool) {
	var b strings.Builder
	b.WriteString("# System commands\n")
	b.WriteString("Use a command prefixed with an exclamation mark (e.g., `!h`).\n")
	fmt.Fprintf(&b, "* `%s` - Select the generative model system prompt.\n", cli.SystemCmdSelectPrompt)
	fmt.Fprintf(&b, "* `%s` - Select from a list of generative model operations.\n", cli.SystemCmdModel)
	fmt.Fprintf(&b, "* `%s` - Select from a list of chat history operations.\n", cli.SystemCmdHistory)
	fmt.Fprintf(&b, "* `%s` - Toggle the input mode.\n", cli.SystemCmdSelectInputMode)
	fmt.Fprintf(&b, "* `%s` - Exit the application.\n", cli.SystemCmdQuit)

	rendered, err := h.renderer.Render(b.String())
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to format instructions: %w", err)), false
	}

	return dataResponse(rendered), false
}

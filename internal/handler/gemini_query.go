package handler

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/cli"
)

// GeminiQuery processes queries to gemini models.
// It implements the MessageHandler interface.
type GeminiQuery struct {
	session  *gemini.ChatSession
	spinner  *cli.Spinner
	renderer *glamour.TermRenderer
}

var _ MessageHandler = (*GeminiQuery)(nil)

// NewGeminiQuery returns a new GeminiQuery message handler.
func NewGeminiQuery(session *gemini.ChatSession, spinner *cli.Spinner,
	style string) (*GeminiQuery, error) {
	renderer, err := glamour.NewTermRenderer(glamour.WithStylePath(style))
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate terminal renderer: %w", err)
	}

	return &GeminiQuery{
		session:  session,
		spinner:  spinner,
		renderer: renderer,
	}, nil
}

// Handle processes the chat message.
func (h *GeminiQuery) Handle(message string) (Response, bool) {
	h.spinner.Start()
	response, err := h.session.SendMessage(message)
	h.spinner.Stop()
	if err != nil {
		return newErrorResponse(err), false
	}

	var b strings.Builder
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			_, _ = fmt.Fprintf(&b, "%s", part)
		}
	}

	rendered, err := h.renderer.Render(b.String())
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to format response: %w", err)), false
	}

	return dataResponse(rendered), false
}

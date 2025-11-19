package handler

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/reugn/gemini-cli/gemini"
)

// GeminiQuery processes queries to gemini models.
// It implements the MessageHandler interface.
type GeminiQuery struct {
	*IO
	session  *gemini.ChatSession
	renderer *glamour.TermRenderer
}

var _ MessageHandler = (*GeminiQuery)(nil)

// NewGeminiQuery returns a new GeminiQuery message handler.
func NewGeminiQuery(io *IO, session *gemini.ChatSession, opts RendererOptions) (*GeminiQuery, error) {
	renderer, err := opts.newTermRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate terminal renderer: %w", err)
	}

	return &GeminiQuery{
		IO:       io,
		session:  session,
		renderer: renderer,
	}, nil
}

// Handle processes the chat message.
func (h *GeminiQuery) Handle(message string) (Response, bool) {
	h.terminal.Spinner.Start()
	defer h.terminal.Spinner.Stop()

	response, err := h.session.SendMessage(message)
	if err != nil {
		return newErrorResponse(err), false
	}

	var b strings.Builder
	for _, candidate := range response.Candidates {
		for _, part := range candidate.Content.Parts {
			_, _ = fmt.Fprint(&b, part.Text)
		}
	}

	rendered, err := h.renderer.Render(b.String())
	if err != nil {
		return newErrorResponse(fmt.Errorf("failed to format response: %w", err)), false
	}

	return dataResponse(rendered), false
}

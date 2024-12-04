package handler

import (
	"fmt"
	"io"

	"github.com/reugn/gemini-cli/internal/cli/color"
)

const (
	empty            = "Empty"
	unchangedMessage = "The selection is unchanged."
)

// Response represents a response from a chat message handler.
type Response interface {
	Print(w io.Writer, prompt string) error
}

type dataResponse string

var _ Response = (*dataResponse)(nil)

func (r dataResponse) Print(w io.Writer, prompt string) error {
	_, err := fmt.Fprintf(w, "%s%s\n", prompt, r)
	return err
}

type errorResponse struct {
	error
}

func newErrorResponse(err error) errorResponse {
	return errorResponse{error: err}
}

var _ Response = (*errorResponse)(nil)

func (r errorResponse) Print(w io.Writer, prompt string) error {
	_, err := fmt.Fprintf(w, "%s%s\n", prompt, color.Red(r.Error()))
	return err
}

func PrintError(w io.Writer, prompt string, err error) {
	_ = newErrorResponse(err).Print(w, prompt)
}

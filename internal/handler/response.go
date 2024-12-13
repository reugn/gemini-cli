package handler

import (
	"fmt"

	"github.com/reugn/gemini-cli/internal/terminal"
)

// Response represents a response from a chat message handler.
type Response interface {
	fmt.Stringer
}

type dataResponse string

var _ Response = (*dataResponse)(nil)

func (r dataResponse) String() string {
	return fmt.Sprintf("%s\n", string(r))
}

//nolint:errname
type errorResponse struct {
	error
}

func newErrorResponse(err error) errorResponse {
	return errorResponse{error: err}
}

var _ Response = (*errorResponse)(nil)

func (r errorResponse) String() string {
	return fmt.Sprintf("%s\n", terminal.Error(r.Error()))
}

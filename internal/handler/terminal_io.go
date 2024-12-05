package handler

import "github.com/reugn/gemini-cli/internal/terminal"

// IO encapsulates terminal details for handlers.
type IO struct {
	terminal       *terminal.IO
	terminalPrompt string
}

// NewIO returns a new IO.
func NewIO(terminal *terminal.IO, terminalPrompt string) *IO {
	return &IO{
		terminal:       terminal,
		terminalPrompt: terminalPrompt,
	}
}

// TerminalPrompt returns the terminal prompt string.
func (io *IO) TerminalPrompt() string {
	return io.terminalPrompt
}

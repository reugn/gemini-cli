package terminal

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/reugn/gemini-cli/internal/cli"
	"github.com/reugn/gemini-cli/internal/terminal/color"
)

var Error = color.Red

// IOConfig represents the configuration settings for IO.
type IOConfig struct {
	User           string
	Multiline      bool
	LineTerminator string
}

// IO encapsulates input/output operations.
type IO struct {
	Reader  *readline.Instance
	Prompt  *Prompt
	Spinner *Spinner
	writer  io.Writer

	Config *IOConfig
}

// NewIO returns a new IO based on the provided configuration.
func NewIO(config *IOConfig) (*IO, error) {
	reader, err := readline.NewEx(&readline.Config{})
	if err != nil {
		return nil, err
	}

	terminalPrompt := NewPrompt(config.User)
	reader.SetPrompt(terminalPrompt.User)
	if config.Multiline {
		// disable history for multiline input mode
		reader.HistoryDisable()
	}

	return &IO{
		Reader:  reader,
		Prompt:  terminalPrompt,
		Spinner: NewSpinner(reader.Stdout(), time.Second, 5),
		writer:  reader.Stdout(),
		Config:  config,
	}, nil
}

// Close releases underlying terminal resources.
func (io *IO) Close() error {
	if io.Reader != nil {
		return io.Reader.Close()
	}
	return nil
}

// Read reads input from the underlying source and returns it as a string.
// If multiline is true, it reads all available lines; otherwise, it reads a single line.
func (io *IO) Read() string {
	if io.Config.Multiline {
		return io.readMultiLine()
	}
	return io.readLine()
}

// Write writes the given string data to the underlying data stream.
func (io *IO) Write(data string) {
	_, _ = fmt.Fprint(io.writer, data)
}

func (io *IO) readLine() string {
	input, err := io.Reader.Readline()
	if err != nil {
		return io.handleReadError(err, len(input))
	}
	return strings.TrimSpace(input)
}

func (io *IO) readMultiLine() string {
	defer io.SetUserPrompt()
	var builder strings.Builder
	for {
		input, err := io.Reader.Readline()
		if err != nil {
			return io.handleReadError(err, builder.Len()+len(input))
		}

		if strings.HasSuffix(input, io.Config.LineTerminator) ||
			strings.HasPrefix(input, cli.SystemCmdPrefix) {
			builder.WriteString(strings.TrimSuffix(input, io.Config.LineTerminator))
			break
		}

		if builder.Len() == 0 {
			io.Reader.SetPrompt(io.Prompt.UserMultilineNext)
		}

		builder.WriteString(input)
		builder.WriteRune('\n')
	}
	return strings.TrimSpace(builder.String())
}

func (io *IO) handleReadError(err error, inputLen int) string {
	if errors.Is(err, readline.ErrInterrupt) {
		if inputLen == 0 {
			// handle as the quit command
			return cli.SystemCmdPrefix + cli.SystemCmdQuit
		}
	} else {
		io.Write(fmt.Sprintf("%s%s\n", io.Prompt.Cli, Error(err.Error())))
	}
	return ""
}

// SetUserPrompt sets the terminal prompt according to the current input mode.
func (io *IO) SetUserPrompt() {
	if io.Config.Multiline {
		io.Reader.SetPrompt(io.Prompt.UserMultiline)
	} else {
		io.Reader.SetPrompt(io.Prompt.User)
	}
}

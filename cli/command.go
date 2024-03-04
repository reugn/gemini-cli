package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/reugn/gemini-cli/cli/color"
	"github.com/reugn/gemini-cli/gemini"
	"google.golang.org/api/iterator"
)

const (
	systemCmdPrefix = "!"
	systemCmdQuit   = "!q"
)

type command interface {
	run(message string) bool
}

type systemCommand struct {
	model  *gemini.ChatSession
	prompt *prompt
}

var _ command = (*systemCommand)(nil)

func newSystemCommand(model *gemini.ChatSession, prompt *prompt) command {
	return &systemCommand{
		model:  model,
		prompt: prompt,
	}
}

func (c *systemCommand) run(message string) bool {
	message = strings.TrimPrefix(message, systemCmdPrefix)
	switch message {
	case "q":
		c.print("Exiting gemini-cli...")
		return true
	case "p":
		c.model.ClearHistory()
		c.print("Cleared the chat history.")
	default:
		c.print("Unknown system command.")
	}
	return false
}

func (c *systemCommand) print(message string) {
	fmt.Printf("%s%s\n", c.prompt.cli, message)
}

type geminiCommand struct {
	model   *gemini.ChatSession
	prompt  *prompt
	spinner *spinner
	writer  *bufio.Writer
	opts    *ChatOpts
}

var _ command = (*geminiCommand)(nil)

func newGeminiCommand(model *gemini.ChatSession, prompt *prompt, opts *ChatOpts) command {
	writer := bufio.NewWriter(os.Stdout)
	return &geminiCommand{
		model:   model,
		prompt:  prompt,
		spinner: newSpinner(5, time.Second, writer),
		writer:  writer,
		opts:    opts,
	}
}

func (c *geminiCommand) run(message string) bool {
	c.printFlush(c.prompt.gemini)
	c.spinner.start()
	if c.opts.Format {
		// requires full markdown for formatting
		c.runBlocking(message)
	} else {
		c.runStreaming(message)
	}
	return false
}

func (c *geminiCommand) runBlocking(message string) {
	response, err := c.model.SendMessage(message)
	c.spinner.stop()
	if err != nil {
		fmt.Println(color.Red(err.Error()))
	} else {
		var buf strings.Builder
		for _, candidate := range response.Candidates {
			for _, part := range candidate.Content.Parts {
				buf.WriteString(fmt.Sprintf("%s", part))
			}
		}
		output, err := glamour.Render(buf.String(), c.opts.Style)
		if err != nil {
			fmt.Printf(color.Red("Failed to format: %s\n"), err)
			fmt.Println(buf.String())
			return
		}
		fmt.Print(output)
	}
}

func (c *geminiCommand) runStreaming(message string) {
	responseIterator := c.model.SendMessageStream(message)
	c.spinner.stop()
	for {
		response, err := responseIterator.Next()
		if err != nil {
			if !errors.Is(err, iterator.Done) {
				fmt.Print(color.Red(err.Error()))
			}
			break
		}
		for _, candidate := range response.Candidates {
			for _, part := range candidate.Content.Parts {
				c.printFlush(fmt.Sprintf("%s", part))
			}
		}
	}
	fmt.Print("\n")
}

func (c *geminiCommand) printFlush(message string) {
	fmt.Fprintf(c.writer, "%s", message)
	c.writer.Flush()
}

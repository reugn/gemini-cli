package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/reugn/gemini-cli/cli/color"
	"github.com/reugn/gemini-cli/gemini"
	"google.golang.org/api/iterator"
)

const systemCmdPrefix = "!"

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
	stream  bool
}

var _ command = (*geminiCommand)(nil)

func newGeminiCommand(model *gemini.ChatSession, prompt *prompt, opts *ChatOpts) command {
	writer := bufio.NewWriter(os.Stdout)
	return &geminiCommand{
		model:   model,
		prompt:  prompt,
		spinner: newSpinner(5, time.Second, writer),
		writer:  writer,
		stream:  opts.Stream,
	}
}

func (c *geminiCommand) run(message string) bool {
	c.printFlush(c.prompt.gemini)
	c.spinner.start()
	if c.stream {
		c.runStreaming(message)
	} else {
		c.runBlocking(message)
	}
	return false
}

func (c *geminiCommand) runBlocking(message string) {
	response, err := c.model.SendMessage(message)
	c.spinner.stop()
	if err != nil {
		fmt.Print(color.Red(err.Error()))
	} else {
		for _, candidate := range response.Candidates {
			for _, part := range candidate.Content.Parts {
				c.printFlush(fmt.Sprintf("%s", part))
			}
		}
	}
	fmt.Print("\n")
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

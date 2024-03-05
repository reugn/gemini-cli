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
	"google.golang.org/api/iterator"
)

const (
	systemCmdPrefix          = "!"
	systemCmdQuit            = "!q"
	systemCmdPurgeHistory    = "!p"
	systemCmdToggleInputMode = "!m"
)

type command interface {
	run(message string) bool
}

type systemCommand struct {
	chat *Chat
}

var _ command = (*systemCommand)(nil)

func newSystemCommand(chat *Chat) command {
	return &systemCommand{
		chat: chat,
	}
}

func (c *systemCommand) run(message string) bool {
	switch message {
	case systemCmdQuit:
		c.print("Exiting gemini-cli...")
		return true
	case systemCmdPurgeHistory:
		c.chat.model.ClearHistory()
		c.print("Cleared the chat history.")
	case systemCmdToggleInputMode:
		if c.chat.opts.Multiline {
			c.print("Switched to single-line input mode.")
			c.chat.reader.HistoryEnable()
			c.chat.opts.Multiline = false
		} else {
			c.print("Switched to multi-line input mode.")
			// disable history for multi-line messages since it is
			// unusable for future requests
			c.chat.reader.HistoryDisable()
			c.chat.opts.Multiline = true
		}
	default:
		c.print("Unknown system command.")
	}
	return false
}

func (c *systemCommand) print(message string) {
	fmt.Printf("%s%s\n", c.chat.prompt.cli, message)
}

type geminiCommand struct {
	chat    *Chat
	spinner *spinner
	writer  *bufio.Writer
}

var _ command = (*geminiCommand)(nil)

func newGeminiCommand(chat *Chat) command {
	writer := bufio.NewWriter(os.Stdout)
	return &geminiCommand{
		chat:    chat,
		spinner: newSpinner(5, time.Second, writer),
		writer:  writer,
	}
}

func (c *geminiCommand) run(message string) bool {
	c.printFlush(c.chat.prompt.gemini)
	c.spinner.start()
	if c.chat.opts.Format {
		// requires the entire response to be formatted
		c.runBlocking(message)
	} else {
		c.runStreaming(message)
	}
	return false
}

func (c *geminiCommand) runBlocking(message string) {
	response, err := c.chat.model.SendMessage(message)
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
		output, err := glamour.Render(buf.String(), c.chat.opts.Style)
		if err != nil {
			fmt.Printf(color.Red("Failed to format: %s\n"), err)
			fmt.Println(buf.String())
			return
		}
		fmt.Print(output)
	}
}

func (c *geminiCommand) runStreaming(message string) {
	responseIterator := c.chat.model.SendMessageStream(message)
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

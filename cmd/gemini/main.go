package main

import (
	"context"
	"errors"
	"os"
	"os/user"

	"github.com/reugn/gemini-cli/cli"
	"github.com/reugn/gemini-cli/gemini"
	"github.com/spf13/cobra"
)

const (
	version   = "0.1.0"
	apiKeyEnv = "GEMINI_API_KEY"
)

var (
	chat   = true
	stream = true
)

// run parses the CLI parameters and executes backup.
func run() int {
	rootCmd := &cobra.Command{
		Short:   "Gemini CLI Tool",
		Version: version,
	}

	rootCmd.Flags().BoolVar(&chat, "chat", true, "start chat session")
	rootCmd.Flags().BoolVar(&stream, "stream", true, "use streaming")

	rootCmd.RunE = func(_ *cobra.Command, _ []string) error {
		if !chat {
			return errors.New("only chat session is supported")
		}
		apiKey := os.Getenv(apiKeyEnv)
		chatSession, err := gemini.NewChatSession(context.Background(), apiKey)
		if err != nil {
			return err
		}
		opts := &cli.ChatOpts{
			Stream: stream,
		}
		chat := cli.NewChat(getCurrentUser(), chatSession, opts)
		chat.StartChat()

		chatSession.Close()
		return nil
	}

	err := rootCmd.Execute()
	if err != nil {
		return 1
	}
	return 0
}

func getCurrentUser() string {
	currentUser, err := user.Current()
	if err != nil {
		return "user"
	}
	return currentUser.Username
}

func main() {
	// start the application
	os.Exit(run())
}

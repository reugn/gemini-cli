package main

import (
	"context"
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

var opts = cli.ChatOpts{}

// run parses the CLI parameters and executes backup.
func run() int {
	rootCmd := &cobra.Command{
		Short:   "Gemini CLI Tool",
		Version: version,
	}

	rootCmd.Flags().BoolVarP(&opts.Format, "format", "f", true, "render markdown-formatted response")
	rootCmd.Flags().StringVarP(&opts.Style, "style", "s", "auto",
		"markdown format style (ascii, dark, light, pink, notty, dracula)")

	rootCmd.RunE = func(_ *cobra.Command, _ []string) error {
		apiKey := os.Getenv(apiKeyEnv)
		chatSession, err := gemini.NewChatSession(context.Background(), apiKey)
		if err != nil {
			return err
		}
		chat := cli.NewChat(getCurrentUser(), chatSession, &opts)
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

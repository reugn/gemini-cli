package main

import (
	"context"
	"os"
	"os/user"

	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/cli"
	"github.com/spf13/cobra"
)

const (
	version   = "0.3.1"
	apiKeyEnv = "GEMINI_API_KEY" //nolint:gosec
)

func run() int {
	rootCmd := &cobra.Command{
		Short:   "Gemini CLI Tool",
		Version: version,
	}

	var opts cli.ChatOpts
	rootCmd.Flags().StringVarP(&opts.Model, "model", "m", gemini.DefaultModel, "generative model name")
	rootCmd.Flags().BoolVarP(&opts.Format, "format", "f", true, "render markdown-formatted response")
	rootCmd.Flags().StringVarP(&opts.Style, "style", "s", "auto",
		"markdown format style (ascii, dark, light, pink, notty, dracula)")
	rootCmd.Flags().BoolVar(&opts.Multiline, "multiline", false, "read input as a multi-line string")
	rootCmd.Flags().StringVarP(&opts.Terminator, "term", "t", "$", "multi-line input terminator")

	rootCmd.RunE = func(_ *cobra.Command, _ []string) error {
		apiKey := os.Getenv(apiKeyEnv)
		chatSession, err := gemini.NewChatSession(context.Background(), opts.Model, apiKey)
		if err != nil {
			return err
		}

		chat, err := cli.NewChat(getCurrentUser(), chatSession, &opts)
		if err != nil {
			return err
		}
		chat.StartChat()

		return chatSession.Close()
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

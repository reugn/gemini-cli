package main

import (
	"context"
	"os"
	"os/user"

	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/chat"
	"github.com/reugn/gemini-cli/internal/config"
	"github.com/spf13/cobra"
)

const (
	version           = "0.3.1"
	apiKeyEnv         = "GEMINI_API_KEY" //nolint:gosec
	defaultConfigPath = "gemini_cli_config.json"
)

func run() int {
	rootCmd := &cobra.Command{
		Short:   "Gemini CLI Tool",
		Version: version,
	}

	var opts chat.Opts
	var configPath string
	rootCmd.Flags().StringVarP(&opts.GenerativeModel, "model", "m", gemini.DefaultModel,
		"generative model name")
	rootCmd.Flags().StringVarP(&opts.Style, "style", "s", "auto",
		"markdown format style (ascii, dark, light, pink, notty, dracula)")
	rootCmd.Flags().BoolVar(&opts.Multiline, "multiline", false,
		"read input as a multi-line string")
	rootCmd.Flags().StringVarP(&opts.LineTerminator, "term", "t", "$",
		"multi-line input terminator")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", defaultConfigPath,
		"path to configuration file in JSON format")

	rootCmd.RunE = func(_ *cobra.Command, _ []string) error {
		configuration, err := config.NewConfiguration(configPath)
		if err != nil {
			return err
		}

		modelBuilder := gemini.NewGenerativeModelBuilder().
			WithName(opts.GenerativeModel).
			WithSafetySettings(configuration.Data.SafetySettings)
		apiKey := os.Getenv(apiKeyEnv)
		chatSession, err := gemini.NewChatSession(context.Background(), modelBuilder, apiKey)
		if err != nil {
			return err
		}

		chatHandler, err := chat.New(getCurrentUser(), chatSession, configuration, &opts)
		if err != nil {
			return err
		}
		chatHandler.Start()

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

package main

import (
	"context"
	"errors"
	"os"
	"os/user"

	"github.com/reugn/gemini-cli/gemini"
	"github.com/reugn/gemini-cli/internal/chat"
	"github.com/reugn/gemini-cli/internal/config"
	"github.com/spf13/cobra"
)

const (
	version           = "0.4.0"
	defaultConfigPath = "gemini_cli_config.json"
)

func run() int {
	rootCmd := &cobra.Command{
		Short:   "Gemini CLI Tool",
		Version: version,
	}

	var (
		opts       chat.Opts
		configPath string
	)

	rootCmd.Flags().StringVarP(&opts.GenerativeModel, "model", "m", gemini.DefaultModel,
		"generative model name")
	rootCmd.Flags().BoolVar(&opts.Multiline, "multiline", false,
		"read input as a multi-line string")
	rootCmd.Flags().StringVarP(&opts.LineTerminator, "term", "t", "$",
		"multi-line input terminator")
	rootCmd.Flags().StringVarP(&opts.StylePath, "style", "s", "auto",
		"markdown format style (ascii, dark, light, pink, notty, dracula, tokyo-night)")
	rootCmd.Flags().IntVarP(&opts.WordWrap, "wrap", "w", 80,
		"line length for response word wrapping")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", defaultConfigPath,
		"path to configuration file in JSON format")

	rootCmd.RunE = func(_ *cobra.Command, _ []string) (err error) {
		configuration, err := config.NewConfiguration(configPath)
		if err != nil {
			return err
		}

		chatSession, err := gemini.NewChatSession(context.Background(), opts.GenerativeModel,
			configuration.Data.GenaiContentConfig())
		if err != nil {
			return err
		}

		chatHandler, err := chat.New(getCurrentUser(), chatSession, configuration, &opts)
		if err != nil {
			return err
		}
		defer func() { err = errors.Join(err, chatHandler.Close()) }()

		chatHandler.Start()
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
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

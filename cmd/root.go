// Package cmd provides the command-line interface for the gic application.
package cmd

import (
	"fmt"

	"gic/internal/config"
	"gic/internal/git"
	"gic/internal/llm"
	"gic/internal/logger"

	"github.com/spf13/cobra"
)

var (
	hash               string
	verbose            bool
	createSampleConfig bool
	pullRequest        bool
	rootCmd            = &cobra.Command{
		Use:   "gic",
		Short: "gic",
		Long:  "gic generates git commit messages based on staged changes.",
		PersistentPreRunE: func(_ *cobra.Command, args []string) error {
			// Set logger level based on the verbose flag
			if verbose {
				logger.SetLogLevel("debug")
			} else {
				logger.SetLogLevel("info")
			}

			// Check for non-flag arguments
			if len(args) > 0 {
				return fmt.Errorf("unexpected arguments: %v", args)
			}

			return nil
		},
	}
)

// Execute runs the root command of the application.
func Execute(version, commit string) error {
	rootCmd.Version = version
	hash = commit
	setVersion()
	rootCmd.RunE = executeCmd
	return rootCmd.Execute()
}

func executeCmd(_ *cobra.Command, _ []string) error {
	l := logger.GetLogger()
	if createSampleConfig {
		l.Debug("Started creating sample configuration")
		err := config.CreateSampleConfig()
		if err != nil {
			return err
		}
		l.Debug("Finish creating sample configuration")
		return nil
	}

	l.Debug("Started executing command")
	l.Debug("Start loading configuration")
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	l.Debug("Finish loading configuration")

	// Include the pullRequest flag in the configuration
	cfg.PR = pullRequest

	gitDiff, err := git.GetGitDiff(cfg)
	if err != nil {
		return err
	}

	l.Debug("Start generating commit message")
	commitMessage, err := llm.GenerateCommitMessage(cfg, gitDiff)
	if err != nil {
		return err
	}
	if commitMessage == "### NO STAGED CHAGES ###" {
		return nil
	}
	l.Info("commit message: " + commitMessage)
	return git.Commit(commitMessage, cfg)
}

func setVersion() {
	template := fmt.Sprintf("gic version: %s commit: %s \n", rootCmd.Version, hash)
	rootCmd.SetVersionTemplate(template)
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "set logging level to verbose")
	rootCmd.PersistentFlags().BoolVarP(
		&createSampleConfig,
		"create-sample-config",
		"s",
		false,
		"create a sample configuration file in the running directory",
	)
	// create a flag for --pull-request or -p default false and is used to generate a message against source main
	rootCmd.PersistentFlags().BoolVarP(
		&pullRequest,
		"pull-request",
		"p",
		false,
		"generate a commit message comparing against main branch",
	)
}

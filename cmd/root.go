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

// const exitCodeFailure = 1

var (
	hash    string
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "gic",
		Short: "gic",
		Long:  "gic generates git commit messages based on staged changes.",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			// Set logger level based on the verbose flag
			if verbose {
				logger.SetLogLevel("debug")
			} else {
				logger.SetLogLevel("info")
			}
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
	// _ = cmd
	// _ = args
	l := logger.GetLogger()
	l.Debug("Started executing command")
	l.Debug("Start loading configuration")
	cfg, err := config.LoadConfig()
	l.Debug("Finish loading configuration")
	if err != nil {
		return err
	}

	l.Debug("Start getting staged changes")
	gitDiff, err := git.GetStagedChanges()
	if err != nil {
		return err
	}
	l.Debug("Finish getting staged changes")

	l.Debug("Start generating commit message")
	commitMessage, err := llm.GenerateCommitMessage(cfg, gitDiff)
	if err != nil {
		return err
	}
	l.Debug("Finish generating commit message")

	l.Debug("Start validating commit message includes changes")
	if commitMessage == "### NO STAGED CHAGES ###" {
		return nil
	}
	l.Debug("Finish validating commit message includes changes")

	return git.Commit(commitMessage, cfg)
}

func setVersion() {
	template := fmt.Sprintf("gic version: %s commit: %s \n", rootCmd.Version, hash)
	rootCmd.SetVersionTemplate(template)
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "set logging level to verbose")
}

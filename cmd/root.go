// Package cmd provides the command-line interface for the gic application.
package cmd

import (
	"fmt"

	"gic/internal/config"
	"gic/internal/git"
	"gic/internal/llm"

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
	}
)

// Execute runs the root command of the application.
func Execute(version, commit string) error {
	rootCmd.Version = version
	hash = commit

	setVersion()

	rootCmd = &cobra.Command{
		Use:   "gic",
		Short: "gic generates git commit messages based on staged changes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd
			_ = args
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			gitDiff, err := git.GetStagedChanges()
			if err != nil {
				return err
			}

			// retrieve the commit message
			commitMessage, err := llm.GenerateCommitMessage(cfg, gitDiff)
			if err != nil {
				return err
			}

			return git.Commit(commitMessage, cfg.Commit)
		},
	}

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		// log.Fatal(err)
		return err
		// os.Exit(exitCodeFailure)
	}
	return nil
}

func setVersion() {
	template := fmt.Sprintf("gic version: %s commit: %s \n", rootCmd.Version, hash)
	rootCmd.SetVersionTemplate(template)
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "set logging level to verbose")
}

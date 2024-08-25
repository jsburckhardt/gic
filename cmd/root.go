package cmd

import (
	"fmt"
	"log"
	"os"

	"gic/internal/config"
	"gic/internal/git"
	"gic/internal/llm"

	"github.com/spf13/cobra"
)

var (
	hash    string
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "gic",
		Short: "gic",
		Long:  "gic generates git commit messages based on staged changes.",
	}
)

func Execute(version, commit string) {
	rootCmd.Version = version
	hash = commit

	setVersion()

	var rootCmd = &cobra.Command{
		Use:   "gic",
		Short: "gic generates git commit messages based on staged changes.",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.LoadConfig()
			if err != nil {
				log.Fatal(err)
			}

			err = config.ValidateConfig(cfg)
			if err != nil {
				log.Fatal(err)
			}

			gitDiff, err := git.GetStagedChanges()
			if err != nil {
				log.Fatal(err)
			}

			// if diff is empty finish
			if gitDiff == "" {
				fmt.Println("No staged changes found.")
				return
			}

			// retrieve the commit message
			commitMessage, err := llm.GenerateCommitMessage(cfg, gitDiff)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Suggested Commit Message:", commitMessage)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setVersion() {
	template := fmt.Sprintf("gic version: %s commit: %s \n", rootCmd.Version, hash)
	rootCmd.SetVersionTemplate(template)
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "set logging level to verbose")
}

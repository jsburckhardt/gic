package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jsburckhardt/gic/internal/config"
	"github.com/jsburckhardt/gic/internal/git"
	"github.com/jsburckhardt/gic/internal/llm"
	"github.com/spf13/cobra"
)

func main() {
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

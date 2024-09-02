// Package git provides functions for interacting with Git repositories.
package git

import (
	"fmt"
	"os/exec"
)

// GetStagedChanges returns the staged changes in the git repository.
// It executes the "git diff --cached" command and returns the output as a string.
// If an error occurs during the execution of the command, it returns an empty string and the error.
func GetStagedChanges() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// Commit commits the staged changes with the generated message.
// it will only print the message unless commit is set to true.
func Commit(message string, commit bool) error {
	cmd := exec.Command("git", "commit", "-m", message)
	if commit {
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		fmt.Println("Suggested commit message:" + message)
	}
	return nil
}

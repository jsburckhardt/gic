// Package git provides functions for interacting with Git repositories.
package git

import (
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

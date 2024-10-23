// Package git provides functions for interacting with Git repositories.
package git

import (
	"bytes"
	"gic/internal/config"
	"gic/internal/logger"
	"os/exec"
	"strings"
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

// GetDiffWithMain returns the diff between the current branch and the main branch.
func GetDiffWithMain() (string, error) {
	// check if it is behind and if it is, return error saying it is behind origin/main
	_, err := isLocalMainBehind()
	if err != nil {
		return "", err
	}
	cmd := exec.Command("git", "diff", "origin/main")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// IsLocalMainBehind checks if the local main branch is behind the origin main branch.
func isLocalMainBehind() (bool, error) {
	// Fetch the latest changes from the origin
	cmd := exec.Command("git", "fetch", "origin")
	if err := cmd.Run(); err != nil {
		return false, err
	}

	// Compare the local main branch with the origin main branch
	cmd = exec.Command("git", "rev-list", "--left-right", "--count", "main...origin/main")
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// Parse the output
	parts := strings.Fields(string(out))
	if len(parts) != 2 {
		return false, err
	}

	// Check if the local main branch is behind
	behind := parts[1] != "0"
	return behind, nil
}

// Commit commits the staged changes with the generated message.
// it will only print the message unless commit is set to true.
func Commit(message string, cfg config.Config, pr bool) error {
	l := logger.GetLogger()
	var err error
	cmd := exec.Command("git", "commit", "-m", message)
	if cfg.ShouldCommit && !pr {
		l.Debug("ShouldCommit True. Committing changes...")
		l.Debug("Commit message: " + message)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err = cmd.Run(); err != nil {
			l.Error("Failed to commit changes", "error", err, "stderr", stderr.String())
			return err
		}
	}
	return nil
}

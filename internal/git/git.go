// Package git provides functions for interacting with Git repositories.
package git

import (
	"bytes"
	"gic/internal/config"
	"gic/internal/logger"
	"os/exec"
	"strings"
)

// consts
const (
	emptyString     = ""
	gitString       = "git"
	diffOutputLimit = 2
	diffsResults    = 1
)

// GetStagedChanges returns the staged changes in the git repository.
// It executes the "git diff --cached" command and returns the output as a string.
// If an error occurs during the execution of the command, it returns an empty string and the error.
func getStagedChanges() (string, error) {
	cmd := exec.Command(gitString, "diff", "--cached")
	out, err := cmd.Output()
	if err != nil {
		return emptyString, err
	}
	return string(out), nil
}

// getDiffWithMain returns the diff between the current branch and the main branch.
func getDiffWithMain() (string, error) {
	// check if it is behind and if it is, return error saying it is behind origin/main
	_, err := isLocalMainBehind()
	if err != nil {
		return emptyString, err
	}
	cmd := exec.Command(gitString, "diff", "origin/main")
	output, err := cmd.Output()
	if err != nil {
		return emptyString, err
	}
	return string(output), nil
}

// IsLocalMainBehind checks if the local main branch is behind the origin main branch.
func isLocalMainBehind() (bool, error) {
	// Fetch the latest changes from the origin
	cmd := exec.Command(gitString, "fetch", "origin")
	if err := cmd.Run(); err != nil {
		return false, err
	}

	// Compare the local main branch with the origin main branch
	cmd = exec.Command(gitString, "rev-list", "--left-right", "--count", "main...origin/main")
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	// Parse the output
	parts := strings.Fields(string(out))
	if len(parts) != diffOutputLimit {
		return false, err
	}

	// Check if the local main branch is behind
	behind := parts[diffsResults] != "0"
	return behind, nil
}

// Commit commits the staged changes with the generated message.
// it will only print the message unless commit is set to true.
func Commit(message string, cfg config.Config) error {
	l := logger.GetLogger()
	var err error
	cmd := exec.Command(gitString, "commit", "-m", message)
	if cfg.ShouldCommit && !cfg.PR {
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

// GetGitDiff returns the diff of the git repository based on the configuration.
func GetGitDiff(cfg config.Config) (string, error) {
	l := logger.GetLogger()
	if cfg.PR {
		l.Debug("Start getting diff with main branch")
		return getDiffWithMain()
	}
	l.Debug("Start getting staged changes")
	return getStagedChanges()
}

package git_test

import (
	"testing"

	"gic/internal/config"
	"gic/internal/git"
)

// create a empty Config struct
var testConfig config.Config

func init() {
	testConfig.PR = false
}

func TestGetStagedChanges(t *testing.T) {
	diff, err := git.GetGitDiff(testConfig)
	if err != nil {
		t.Fatal(err)
	}

	if diff == "" {
		t.Fatal("Expected diff output, got empty string")
	}
}

package git_test

import (
	"testing"

	"github.com/jsburckhardt/gic/internal/git"
)

func TestGetStagedChanges(t *testing.T) {
	diff, err := git.GetStagedChanges()
	if err != nil {
		t.Fatal(err)
	}

	if diff == "" {
		t.Fatal("Expected diff output, got empty string")
	}
}

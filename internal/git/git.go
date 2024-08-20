package git

import (
    "os/exec"
)

func GetStagedChanges() (string, error) {
    cmd := exec.Command("git", "diff", "--cached")
    out, err := cmd.Output()
    if err != nil {
        return "", err
    }
    return string(out), nil
}

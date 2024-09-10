// main package for the gic
package main

import (
	"gic/cmd"
	"gic/internal/logger"
)

var (
	version = "local"
	commit  = "n/a"
)

func main() {
	logger.InitLogger()
	l := logger.GetLogger()
	err := cmd.Execute(version, commit)
	if err != nil {
		l.Error("failed to execute command", "error", err)
		panic(err)
	}
}

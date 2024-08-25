// main package for the gic
package main

import (
	"gic/cmd"
)

var (
	version = "edge"
	commit  = "n/a"
)

func main() {
	cmd.Execute(version, commit)
}

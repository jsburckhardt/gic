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
	err := cmd.Execute(version, commit)
	if err != nil {
		panic(err)
	}
}

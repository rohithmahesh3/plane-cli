package main

import (
	"os"

	"github.com/rohithmahesh3/plane-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

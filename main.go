package main

import (
	"fmt"
	"os"

	"github.com/kolo/ledger-go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		exitWithError(err)
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

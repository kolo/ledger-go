package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

const (
	ledgerDirEnvKey = "LEDGER_DIR"
)

func main() {
	ledgerDir := os.Getenv(ledgerDirEnvKey)
	if ledgerDir == "" {
		exitWithErr(errors.New(fmt.Sprintf("%s is not defined", ledgerDirEnvKey)))
	}

	fmt.Println(len(read(ledgerDir)))
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}

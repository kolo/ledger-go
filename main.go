package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	configFilename  = "config.json"
	ledgerDirEnvKey = "LEDGER_DIR"
)

func main() {
	ledgerDir := os.Getenv(ledgerDirEnvKey)
	if ledgerDir == "" {
		exitWithErr(errors.New(fmt.Sprintf("%s is not defined", ledgerDirEnvKey)))
	}

	config, err := loadUserConfig(filepath.Join(ledgerDir, configFilename))
	if err != nil {
		exitWithErr(err)
	}

	balanceReport(newSimpleReader(config.Accounts, read(ledgerDir)))
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}

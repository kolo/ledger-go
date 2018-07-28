package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	configFilename  = "config.json"
	ledgerDirEnvKey = "LEDGER_DIR"
)

func main() {
	cmd := &cobra.Command{
		Use: "ledger",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(newCommand("balance", func(ledgerDir string, config *userConfig) {
		balanceReport(newSimpleReader(config.Accounts, read(ledgerDir)))
	}))
	cmd.AddCommand(newCommand("expense", func(ledgerDir string, config *userConfig) {
		expensesReport(newSimpleReader(config.Accounts, read(ledgerDir)), config.Accounts)
	}))

	if err := cmd.Execute(); err != nil {
		exitWithErr(err)
	}
}

func newCommand(name string, run func(string, *userConfig)) *cobra.Command {
	return &cobra.Command{
		Use: name,
		Run: func(*cobra.Command, []string) {
			ledgerDir := os.Getenv(ledgerDirEnvKey)
			if ledgerDir == "" {
				exitWithErr(errors.New(fmt.Sprintf("%s is not defined", ledgerDirEnvKey)))
			}

			config, err := loadUserConfig(filepath.Join(ledgerDir, configFilename))
			if err != nil {
				exitWithErr(err)
			}

			run(ledgerDir, config)
		},
	}
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}

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
	cmd.AddCommand(balanceCommand())
	cmd.AddCommand(expensesCommand())

	if err := cmd.Execute(); err != nil {
		exitWithErr(err)
	}
}

func balanceCommand() *cobra.Command {
	return &cobra.Command{
		Use: "balance",
		Run: func(*cobra.Command, []string) {
			ledgerDir := os.Getenv(ledgerDirEnvKey)
			if ledgerDir == "" {
				exitWithErr(errors.New(fmt.Sprintf("%s is not defined", ledgerDirEnvKey)))
			}

			config, err := loadUserConfig(filepath.Join(ledgerDir, configFilename))
			if err != nil {
				exitWithErr(err)
			}

			balanceReport(newSimpleReader(config.Accounts, read(ledgerDir)))
		},
	}
}

func expensesCommand() *cobra.Command {
	return &cobra.Command{
		Use: "expenses",
		Run: func(*cobra.Command, []string) {
			ledgerDir := os.Getenv(ledgerDirEnvKey)
			if ledgerDir == "" {
				exitWithErr(errors.New(fmt.Sprintf("%s is not defined", ledgerDirEnvKey)))
			}

			config, err := loadUserConfig(filepath.Join(ledgerDir, configFilename))
			if err != nil {
				exitWithErr(err)
			}

			expensesReport(newSimpleReader(config.Accounts, read(ledgerDir)), config.Accounts)
		},
	}
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	configFilename  = "config.json"
	ledgerDirEnvKey = "LEDGER_DIR"
)

func main() {
	env, err := loadEnvironment()
	if err != nil {
		exitWithErr(err)
	}

	cmd := &cobra.Command{
		Use: "ledger",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(newExpensesCommand(env))
	cmd.AddCommand(newBalanceCommand(env))
	cmd.AddCommand(newExportCommand(env))
	cmd.AddCommand(newLogCommand(env))

	if err := cmd.Execute(); err != nil {
		exitWithErr(err)
	}
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}

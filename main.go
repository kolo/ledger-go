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
	cmd := &cobra.Command{
		Use: "ledger",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(newExpensesCommand())
	cmd.AddCommand(newBalanceCommand())

	if err := cmd.Execute(); err != nil {
		exitWithErr(err)
	}
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %v", err)
	os.Exit(1)
}

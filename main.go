package main

import (
	"fmt"
	"os"

	"github.com/kolo/ledger-go/cmd"
)

const (
	configFilename  = "config.json"
	ledgerDirEnvKey = "LEDGER_DIR"
)

func main() {
	if err := cmd.Execute(); err != nil {
		exitWithError(err)
	}

	// env, err := loadEnvironment()
	// if err != nil {
	// 	exitWithErr(err)
	// }

	// cmd := &cobra.Command{
	// 	Use: "ledger",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		cmd.Help()
	// 	},
	// }

	// cmd.AddCommand(newExpensesCommand(env))
	// cmd.AddCommand(newBalanceCommand(env))
	// cmd.AddCommand(newExportCommand(env))
	// cmd.AddCommand(newLogCommand(env))

	// if err := cmd.Execute(); err != nil {
	// 	exitWithErr(err)
	// }
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

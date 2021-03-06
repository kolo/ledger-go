package cmd

import "github.com/spf13/cobra"

var ledgerCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

func Execute() error {
	ledgerCmd.AddCommand(newBalanceCmd())
	ledgerCmd.AddCommand(newExpensesCmd())
	ledgerCmd.AddCommand(newExportCmd())
	ledgerCmd.AddCommand(newLogCmd())

	return ledgerCmd.Execute()
}

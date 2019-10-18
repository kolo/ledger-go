package cmd

import "github.com/spf13/cobra"

var ledgerCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

func Execute() error {
	ledgerCmd.AddCommand(newBalanceCmd())
	return ledgerCmd.Execute()
}

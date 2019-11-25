package cmd

import (
	"io"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
)

func newBalanceCmd() *cobra.Command {
	balance := newBaseCmd()

	cmd := &cobra.Command{
		Use: "balance",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return balance.Run(cmd.OutOrStdout())
		},
	}

	balance.addFlags(cmd.Flags())

	balance.run = func(cfg *ledger.Config, iter ledger.RecordIterator, stdout io.Writer) error {
		ledger.BalanceReport(iter, cfg.Assets, stdout)
		return nil
	}

	return cmd
}

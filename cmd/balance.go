package cmd

import (
	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
)

func newBalanceCmd() *cobra.Command {
	balance := newBaseCmd()

	cmd := &cobra.Command{
		Use: "balance",
		RunE: func(*cobra.Command, []string) error {
			return balance.Run()
		},
	}

	balance.addFlags(cmd.Flags())

	balance.run = func(cfg *ledger.Config, iter ledger.RecordIterator) error {
		ledger.BalanceReport(iter, cfg.Assets)
		return nil
	}

	return cmd
}

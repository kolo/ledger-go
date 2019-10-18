package cmd

import (
	"time"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
)

var unixEpoch = time.Unix(0, 0)

type balanceCmd struct {
	*baseCmd
}

func newBalanceCmd() *cobra.Command {
	balance := &balanceCmd{
		baseCmd: newBaseCmd(),
	}

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

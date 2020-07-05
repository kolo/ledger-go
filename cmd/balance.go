package cmd

import (
	"io"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type balanceCmd struct {
	*baseCmd

	monthlyReport bool
}

func newBalanceCmd() *cobra.Command {
	balance := &balanceCmd{
		baseCmd: newBaseCmd(),
	}

	cmd := &cobra.Command{
		Use: "balance",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return balance.Run(cmd.OutOrStdout())
		},
	}

	balance.addFlags(cmd.Flags())

	balance.run = func(cfg *ledger.Config, iter ledger.RecordIterator, stdout io.Writer) error {
		if balance.monthlyReport {
			ledger.MonthlyBalanceReport(iter, cfg.Assets, stdout)
			return nil
		}
		ledger.BalanceReport(iter, cfg.Assets, stdout)
		return nil
	}

	return cmd
}

func (c *balanceCmd) addFlags(flags *pflag.FlagSet) {
	c.baseCmd.addFlags(flags)

	flags.BoolVarP(&c.monthlyReport, "monthly", "", false, "display balance by month")
}

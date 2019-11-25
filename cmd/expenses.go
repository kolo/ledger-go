package cmd

import (
	"io"

	"github.com/kolo/ledger-go/ledger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type expensesCmd struct {
	*baseCmd
	weeklyReport bool
}

func newExpensesCmd() *cobra.Command {
	expenses := &expensesCmd{
		baseCmd: newBaseCmd(),
	}

	cmd := &cobra.Command{
		Use: "expenses",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return expenses.Run(cmd.OutOrStdout())
		},
	}

	expenses.addFlags(cmd.Flags())

	expenses.run = func(cfg *ledger.Config, iter ledger.RecordIterator, stdout io.Writer) error {
		if expenses.weeklyReport {
			ledger.WeeklyExpensesReport(iter, cfg.Assets, stdout)
			return nil
		}

		ledger.ExpensesReport(iter, stdout)
		return nil
	}

	return cmd
}

func (c *expensesCmd) addFlags(flags *pflag.FlagSet) {
	c.baseCmd.addFlags(flags)

	flags.BoolVarP(&c.weeklyReport, "weekly", "", false, "group expenses by week")
}

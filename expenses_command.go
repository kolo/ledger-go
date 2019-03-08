package main

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

// expensesCommand implements "expenses" command.
type expensesCommand struct {
	cmd *cobra.Command
	env *environment

	// filters
	accounts *accountFilter
	period   *dateRangeFilter

	// grouping
	weekly *bool
}

type filterFunc func(*record) *record

func newExpensesCommand(env *environment) *cobra.Command {
	c := &expensesCommand{
		env:      env,
		accounts: newAccountFilter(),
		period:   newDateRangeFilter(),
	}

	c.cmd = &cobra.Command{
		Use: "expenses [OPTIONS]",
		Run: func(*cobra.Command, []string) {
			c.expenses()
		},
	}
	c.addFlags()

	return c.Cmd()
}

// Cmd initializes an instance of cobra.Command.
func (c *expensesCommand) Cmd() *cobra.Command {
	return c.cmd
}

func (c *expensesCommand) addFlags() {
	flags := c.cmd.Flags()

	c.accounts.addFlags(flags)
	c.period.addFlags(flags)

	c.weekly = flags.BoolP("weekly", "", false, "group expenses by week")
}

func (c *expensesCommand) expenses() {
	reader := newFilteredReader(c.env.reader(), c.filter())

	if *c.weekly {
		weeklyExpensesReport(reader, c.assets())
		return
	}

	expensesReport(reader, c.assets())
}

func (c *expensesCommand) assets() []string {
	return c.env.Accounts
}

func (c *expensesCommand) filter() filterFunc {
	return func(r *record) *record {
		if r == nil {
			return nil
		}

		if r.recordType() != recordTypeExpense {
			return nil
		}

		r = c.accounts.filter(r)
		if r == nil {
			return nil
		}

		r = c.period.filter(r)
		if r == nil {
			return nil
		}

		return r
	}
}

func expensesReport(rd recordReader, assets []string) {
	expenses := report{}

	for _, asset := range assets {
		expenses[asset] = &reportItem{
			account: &account{
				name:  asset,
				asset: true,
			},
			total: decimal.Zero,
		}
	}

	for {
		r := rd.Next()
		if r == nil {
			break
		}

		if ri, found := expenses[r.credit.name]; found {
			ri.increase(r.amount)
		}
	}

	printReport(expenses)
}

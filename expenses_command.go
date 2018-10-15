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
	credit *string
	debit  *string
	period *dateRangeFilter

	// grouping
	weekly *bool
}

type filterFunc func(*record) *record

func newExpensesCommand(env *environment) *cobra.Command {
	c := &expensesCommand{
		env:    env,
		period: newDateRangeFilter(),
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

	c.period.addFlags(flags)

	c.credit = flags.StringP("credit", "", "", "filter by the credit account")
	c.debit = flags.StringP("debit", "", "", "filter by the debit account")
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
	credit := *c.credit
	if credit == "" {
		return c.env.Accounts
	}

	return []string{credit}
}

func (c *expensesCommand) filter() filterFunc {
	debit := *c.debit

	return func(r *record) *record {
		if r == nil {
			return nil
		}

		if r.recordType() != recordTypeExpense {
			return nil
		}

		r = c.period.filter(r)
		if r == nil {
			return nil
		}

		if debit != "" && debit != r.debit.name {
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

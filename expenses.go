package main

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
)

// expensesCommand implements "expenses" command.
type expensesCommand struct {
	cmd *cobra.Command
	env *environment

	from *dateFlag
}

type filterFunc func(*record) *record

func newExpensesCommand(env *environment) *cobra.Command {
	c := &expensesCommand{
		env:  env,
		from: &dateFlag{},
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
	flags.VarP(c.from, "from", "", "set a starting date")
}

func (c *expensesCommand) expenses() {
	expensesReport(c.env.reader(), c.env.Accounts, c.filter())
}

func (c *expensesCommand) filter() filterFunc {
	from := c.from.value

	return func(r *record) *record {
		if r == nil {
			return nil
		}

		if r.recordType() != recordTypeExpense {
			return nil
		}

		if r.recordedAt.Before(from) {
			return nil
		}

		return r
	}
}

func expensesReport(rd recordReader, assets []string, filter filterFunc) {
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

		r = filter(r)
		if r == nil {
			continue
		}

		if ri, found := expenses[r.credit.name]; found {
			ri.increase(r.amount)
		}
	}

	for _, ri := range expenses {
		fmt.Printf("%5s: %6s\n", ri.account.name, ri.total.StringFixed(2))
	}
	fmt.Printf("Total: %v\n", expenses.total())
}

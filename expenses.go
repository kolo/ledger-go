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

	// filters
	credit *string
	debit  *string
	from   *dateFlag
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

	c.credit = flags.StringP("credit", "", "", "filter by the credit account")
	c.debit = flags.StringP("debit", "", "", "filter by the debit account")
}

func (c *expensesCommand) expenses() {
	expensesReport(c.env.reader(), c.assets(), c.filter())
}

func (c *expensesCommand) assets() []string {
	credit := *c.credit
	if credit == "" {
		return c.env.Accounts
	}

	return []string{credit}
}

func (c *expensesCommand) filter() filterFunc {
	from := c.from.value
	debit := *c.debit

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

		if debit != "" && debit != r.debit.name {
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
		if ri.total.Equal(decimal.Zero) {
			continue
		}
		fmt.Printf("%5s: %6s\n", ri.account.name, ri.total.StringFixed(2))
	}
	fmt.Printf("Total: %6s\n", expenses.total().StringFixed(2))
}
